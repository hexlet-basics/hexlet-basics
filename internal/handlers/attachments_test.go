package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob/memblob"

	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/testsupport"
)

// tinyPNG is a minimal (1x1) PNG so uploads carry real image bytes.
var tinyPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
	0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

// newAttachmentRouter builds the real router over an in-memory blob bucket and a
// txdb-backed ent client. The api side is a stub — these tests exercise only the
// multipart/blob routes, which live outside the generated server.
func newAttachmentRouter(t *testing.T) http.Handler {
	t.Helper()
	db := testsupport.NewClient(t)
	bucket := memblob.OpenBucket(nil)
	t.Cleanup(func() { _ = bucket.Close() })
	assets := assetstore.New(db, bucket, "http://assets.example.test")
	att := handlers.NewAttachmentHandler(assets)
	gh := handlers.NewGitHubWebhookHandler(db, &testsupport.RecordingEnqueuer{}, "")
	apiStub := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	return handlers.NewRouter(apiStub, att, gh)
}

// uploadRequest builds a multipart POST with one file part whose declared
// Content-Type can differ from the bytes; assetstore derives the stored type
// from content rather than trusting this header.
func uploadRequest(t *testing.T, filename, contentType string, data []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	if contentType != "" {
		h.Set("Content-Type", contentType)
	}
	part, err := mw.CreatePart(h)
	require.NoError(t, err)
	_, err = part.Write(data)
	require.NoError(t, err)
	require.NoError(t, mw.Close())

	req := httptest.NewRequest(http.MethodPost, "/admin/attachments", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

func TestUploadAttachmentAndDownload(t *testing.T) {
	router := newAttachmentRouter(t)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, uploadRequest(t, "cover.png", "application/octet-stream", tinyPNG))

	require.Equal(t, http.StatusCreated, rec.Code)

	var att struct {
		ID          int    `json:"id"`
		URL         string `json:"url"`
		Filename    string `json:"filename"`
		ContentType string `json:"contentType"`
		ByteSize    int64  `json:"byteSize"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &att))

	assert.NotZero(t, att.ID)
	assert.Equal(t, "cover.png", att.Filename)
	assert.Equal(t, "image/png", att.ContentType)
	assert.Equal(t, int64(len(tinyPNG)), att.ByteSize)
	// The configured public origin is canonical even when the request arrived
	// through a different host or proxy.
	assert.True(t, strings.HasPrefix(att.URL, "http://assets.example.test/storage/"),
		"url must use the configured public origin, got %q", att.URL)
	assert.Contains(t, att.URL, ".png", "the read URL uses the detected image extension")

	// The returned url must actually serve the stored bytes.
	dl := httptest.NewRecorder()
	router.ServeHTTP(dl, httptest.NewRequest(http.MethodGet, att.URL, nil))

	require.Equal(t, http.StatusOK, dl.Code)
	assert.Equal(t, "image/png", dl.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", dl.Header().Get("X-Content-Type-Options"))
	got, err := io.ReadAll(dl.Body)
	require.NoError(t, err)
	assert.Equal(t, tinyPNG, got)
}

func TestUploadAttachmentRejectsUnsupportedType(t *testing.T) {
	router := newAttachmentRouter(t)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, uploadRequest(t, "notes.png", "image/png", []byte("hello")))

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var body struct {
		Errors map[string][]string `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.NotEmpty(t, body.Errors["file"], "validation error keyed by the file field")
}

func TestUploadAttachmentRequiresFilePart(t *testing.T) {
	router := newAttachmentRouter(t)

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	require.NoError(t, mw.WriteField("other", "x"))
	require.NoError(t, mw.Close())
	req := httptest.NewRequest(http.MethodPost, "/admin/attachments", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestDownloadUnknownKeyIsNotFound(t *testing.T) {
	router := newAttachmentRouter(t)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/storage/does-not-exist.png", nil))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
