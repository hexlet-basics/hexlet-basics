package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"strings"
	"testing"

	"github.com/go-pkgz/auth/v2/token"
	"github.com/golang-jwt/jwt/v5"
	ht "github.com/ogen-go/ogen/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob/memblob"

	"hexletbasics/internal/api"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/config"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/ids"
	"hexletbasics/internal/progress"
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

type authenticatedHandler struct {
	next    http.Handler
	cookies []*http.Cookie
	xsrf    string
}

func (h authenticatedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range h.cookies {
		r.AddCookie(cookie)
	}
	r.Header.Set("X-XSRF-TOKEN", h.xsrf)
	h.next.ServeHTTP(w, r)
}

// newAttachmentRouterStack builds the real router over an in-memory blob
// bucket and a transaction-bound ent client. Uploads run through the generated
// ogen server — its multipart decoder and SecurityHandler — exactly as in
// production; downloads use the blob read route mounted beside it.
func newAttachmentRouterStack(t *testing.T, admin bool) (http.Handler, []*http.Cookie, string) {
	t.Helper()
	db, transactor := testsupport.NewClientWithTransactor(t)
	bucket := memblob.OpenBucket(nil)
	t.Cleanup(func() { _ = bucket.Close() })
	assets := assetstore.New(db, bucket, "http://assets.example.test")
	translator := testsupport.NewTranslator(t)
	errorHandler := testsupport.NewAPIErrorHandler(t, translator)
	cfg := &config.Config{JWTSecret: "test-secret", EmailTokenSecret: "test-email-secret"}
	enqueuer := &testsupport.RecordingEnqueuer{DB: db}
	handler := handlers.NewServer(
		db,
		cfg,
		enqueuer,
		enqueuer,
		enqueuer,
		progress.New(db, transactor, &testsupport.RecordingEventPublisher{}, testsupport.NewStubExerciseRunner()),
		assets,
		testsupport.NewRecordingRegistrar(db),
		&testsupport.RecordingEventPublisher{},
		translator,
		errorHandler,
	)
	apiServer, err := api.NewServer(
		handler,
		handler.AuthHandler(),
		api.WithErrorHandler(errorHandler.Write),
		api.WithNotFound(handlers.NewNotFoundHandler(translator)),
		api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
	)
	require.NoError(t, err)
	att := handlers.NewAttachmentHandler(assets, errorHandler)
	gh := handlers.NewGitHubWebhookHandler(db, enqueuer, "", translator)
	router := translator.Middleware(handlers.NewRouter(apiServer, att, gh, handler.AuthHandler()))

	jti := ids.New()
	rec := httptest.NewRecorder()
	u, err := db.User.Create().
		SetEmail("attachment-user-" + ids.New() + "@example.com").
		SetAdmin(admin).
		Save(t.Context())
	require.NoError(t, err)
	authUser := &token.User{ID: strconv.Itoa(u.ID), Name: "attachment-admin@example.com"}
	_, err = token.NewService(token.Opts{
		SecretReader: token.SecretFunc(func(string) (string, error) {
			return cfg.JWTSecret, nil
		}),
	}).Set(rec, token.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       jti,
			Audience: jwt.ClaimStrings{"hexlet-basics"},
		},
		User:         authUser,
		AuthProvider: &token.AuthProvider{Name: "password"},
	})
	require.NoError(t, err)

	return router, rec.Result().Cookies(), jti
}

func newAttachmentRouter(t *testing.T) http.Handler {
	t.Helper()
	router, cookies, xsrf := newAttachmentRouterStack(t, true)
	return authenticatedHandler{next: router, cookies: cookies, xsrf: xsrf}
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

func uploadAttachmentURL(t *testing.T, router http.Handler) string {
	t.Helper()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, uploadRequest(t, "cover.png", "image/png", tinyPNG))
	require.Equal(t, http.StatusCreated, rec.Code)

	var response struct {
		URL string `json:"url"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	require.NotEmpty(t, response.URL)
	return response.URL
}

func TestUploadAttachmentAndDownload(t *testing.T) {
	router := newAttachmentRouter(t)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, uploadRequest(t, "cover.png", "application/octet-stream", tinyPNG))

	require.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

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
	assert.Equal(t, strconv.Itoa(len(tinyPNG)), dl.Header().Get("Content-Length"))
	assert.NotEmpty(t, dl.Header().Get("Last-Modified"))
	assert.Equal(t, "nosniff", dl.Header().Get("X-Content-Type-Options"))
	got, err := io.ReadAll(dl.Body)
	require.NoError(t, err)
	assert.Equal(t, tinyPNG, got)
}

// TestAdminUploadAttachmentThroughGeneratedClient drives the upload with the
// generated client, so the generated multipart encoder and decoder meet over
// the contract and the stored row is what the response names.
func TestAdminUploadAttachmentThroughGeneratedClient(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()

	att, err := h.Client.AdminUploadAttachment(ctx, &api.AttachmentUploadFormMultipart{
		File: ht.MultipartFile{Name: "cover.png", File: bytes.NewReader(tinyPNG), Size: int64(len(tinyPNG))},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	assert.Equal(t, "cover.png", att.Filename)
	assert.Equal(t, "image/png", att.ContentType)
	assert.Equal(t, int64(len(tinyPNG)), att.ByteSize)
	stored, err := h.DB.Attachment.Get(ctx, int(att.ID))
	require.NoError(t, err)
	assert.Equal(t, "cover.png", stored.Filename)
}

func TestUploadAttachmentUsesExactContractAuthentication(t *testing.T) {
	t.Run("missing session", func(t *testing.T) {
		router, _, _ := newAttachmentRouterStack(t, true)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, uploadRequest(t, "cover.png", "image/png", tinyPNG))

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))
	})

	t.Run("current database user is not admin", func(t *testing.T) {
		router, cookies, xsrf := newAttachmentRouterStack(t, false)
		rec := httptest.NewRecorder()

		authenticatedHandler{next: router, cookies: cookies, xsrf: xsrf}.
			ServeHTTP(rec, uploadRequest(t, "cover.png", "image/png", tinyPNG))

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("missing xsrf header", func(t *testing.T) {
		router, cookies, _ := newAttachmentRouterStack(t, true)
		rec := httptest.NewRecorder()

		authenticatedHandler{next: router, cookies: cookies}.
			ServeHTTP(rec, uploadRequest(t, "cover.png", "image/png", tinyPNG))

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestDownloadAttachmentHTTPFeatures(t *testing.T) {
	router := newAttachmentRouter(t)
	url := uploadAttachmentURL(t, router)

	full := httptest.NewRecorder()
	router.ServeHTTP(full, httptest.NewRequest(http.MethodGet, url, nil))
	require.Equal(t, http.StatusOK, full.Code)
	lastModified := full.Header().Get("Last-Modified")
	require.NotEmpty(t, lastModified)

	t.Run("range", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Range", "bytes=1-8")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusPartialContent, rec.Code)
		assert.Equal(t, "bytes 1-8/"+strconv.Itoa(len(tinyPNG)), rec.Header().Get("Content-Range"))
		assert.Equal(t, "8", rec.Header().Get("Content-Length"))
		assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, tinyPNG[1:9], rec.Body.Bytes())
	})

	t.Run("invalid range", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Range", "bytes=1000-2000")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusRequestedRangeNotSatisfiable, rec.Code)
		assert.Equal(t, "bytes */"+strconv.Itoa(len(tinyPNG)), rec.Header().Get("Content-Range"))
		assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	})

	t.Run("head", func(t *testing.T) {
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, httptest.NewRequest(http.MethodHead, url, nil))

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, strconv.Itoa(len(tinyPNG)), rec.Header().Get("Content-Length"))
		assert.Equal(t, lastModified, rec.Header().Get("Last-Modified"))
		assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
		assert.Empty(t, rec.Body.Bytes())
	})

	t.Run("not modified", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("If-Modified-Since", lastModified)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotModified, rec.Code)
		assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
		assert.Empty(t, rec.Body.Bytes())
	})

	t.Run("if range matches", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Range", "bytes=1-8")
		req.Header.Set("If-Range", lastModified)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusPartialContent, rec.Code)
		assert.Equal(t, tinyPNG[1:9], rec.Body.Bytes())
	})

	t.Run("if range is stale", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Range", "bytes=1-8")
		req.Header.Set("If-Range", "Mon, 02 Jan 2006 15:04:05 GMT")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, tinyPNG, rec.Body.Bytes())
	})
}

func TestUploadValidationErrorUsesRequestLocale(t *testing.T) {
	router := newAttachmentRouter(t)
	req := uploadRequest(t, "notes.png", "image/png", []byte("hello"))
	req.Header.Set("Accept-Language", "es")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var body api.ValidationError
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, []string{"Tipo de archivo no compatible"}, body.Errors["file"])
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

	// A form without the required part is malformed per the contract, so the
	// generated decoder rejects it like any other undecodable request.
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))
}

func TestDownloadUnknownKeyIsNotFound(t *testing.T) {
	router := newAttachmentRouter(t)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/storage/does-not-exist.png", nil))

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
		"type":"about:blank",
		"title":"Not Found",
		"status":404
	}`, rec.Body.String())
}

func TestUploadAttachmentRejectsOversizedBody(t *testing.T) {
	router := newAttachmentRouter(t)
	// Valid PNG bytes padded past the whole-body cap, so only the router can stop
	// it before the bytes reach the bucket.
	data := append(append([]byte{}, tinyPNG...), make([]byte, handlers.UploadBodyLimit)...)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, uploadRequest(t, "huge.png", "image/png", data))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))
}

func TestUploadAttachmentRejectsFileOverAssetLimit(t *testing.T) {
	router := newAttachmentRouter(t)
	// Over the file limit but inside the body cap's envelope allowance, so the
	// request reaches assetstore and fails as a field error, not a transport one.
	data := append(append([]byte{}, tinyPNG...), make([]byte, assetstore.MaxUploadBytes)...)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, uploadRequest(t, "large.png", "image/png", data))

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var body api.ValidationError
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.NotEmpty(t, body.Errors["file"])
}
