package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/localization"
)

// multipartOverheadBytes leaves room for the multipart envelope while the
// asset store independently enforces the exact file-byte limit.
const multipartOverheadBytes = 1 << 20

// AttachmentHandler serves the multipart upload + blob read path that lives
// OUTSIDE the ogen-generated router. ogen cannot generate multipart/form-data
// (see internal/apigen/ogen.yml), so `POST /admin/attachments` and the
// `GET /storage/{key}` read-back are plain net/http adapters mounted alongside
// the generated api.Server by NewRouter. Asset lifecycle invariants live in
// assetstore.Store rather than in this transport adapter.
type AttachmentHandler struct {
	assets *assetstore.Store
	i18n   *localization.Translator
	errors *APIErrorHandler
}

// NewAttachmentHandler wires the HTTP adapter to the shared asset store.
func NewAttachmentHandler(
	assets *assetstore.Store,
	translator *localization.Translator,
	errorHandler *APIErrorHandler,
) *AttachmentHandler {
	return &AttachmentHandler{assets: assets, i18n: translator, errors: errorHandler}
}

// attachmentResponse mirrors the OpenAPI `Attachment` schema exactly (camelCase,
// int64 byteSize). It is hand-serialized because the endpoint is outside the
// generated layer — a field-name drift here would silently break the hey-api
// client, so the tags are the contract.
type attachmentResponse struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	ByteSize    int64  `json:"byteSize"`
}

// Upload handles `POST /admin/attachments`: read the single multipart `file`
// part, store its bytes in the bucket under a fresh key, record the metadata in
// ent, and return the 201 Attachment the admin form references by id.
func (h *AttachmentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Bound the body before parsing so a huge upload can't be buffered whole.
	r.Body = http.MaxBytesReader(w, r.Body, assetstore.MaxUploadBytes+multipartOverheadBytes)

	file, header, err := r.FormFile("file")
	if err != nil {
		// A body over the limit surfaces here as *http.MaxBytesError.
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			h.errors.Write(ctx, w, r, newValidationError("file", h.i18n.Text(ctx, localization.FileTooLarge)))
			return
		}
		h.errors.Write(ctx, w, r, newValidationError("file", h.i18n.Text(ctx, localization.FileRequired)))
		return
	}
	defer func() { _ = file.Close() }()

	att, err := h.assets.Put(ctx, assetstore.Upload{
		Filename: header.Filename,
		Body:     file,
	})
	if errors.Is(err, assetstore.ErrUnsupportedMediaType) {
		h.errors.Write(ctx, w, r, newValidationError("file", h.i18n.Text(ctx, localization.UnsupportedFileType)))
		return
	}
	if errors.Is(err, assetstore.ErrTooLarge) {
		h.errors.Write(ctx, w, r, newValidationError("file", h.i18n.Text(ctx, localization.FileTooLarge)))
		return
	}
	if err != nil {
		h.errors.Write(ctx, w, r, fmt.Errorf("store uploaded attachment: %w", err))
		return
	}

	writeJSON(w, http.StatusCreated, attachmentResponse{
		ID:          att.ID,
		URL:         att.URL,
		Filename:    att.Filename,
		ContentType: att.ContentType,
		ByteSize:    att.ByteSize,
	})
}

// Download handles `GET /storage/{key}`: stream the stored bytes back with the
// content type they were written with. This is the read side of the `url` the
// uploader returns — without it that url would point at nothing.
func (h *AttachmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := r.PathValue("key")

	reader, err := h.assets.Open(ctx, key)
	if err != nil {
		if errors.Is(err, assetstore.ErrNotFound) {
			h.errors.Write(ctx, w, r, withHTTPStatus(http.StatusNotFound, err))
			return
		}
		h.errors.Write(ctx, w, r, fmt.Errorf("open stored attachment %q: %w", key, err))
		return
	}
	defer func() { _ = reader.Close() }()

	if reader.ContentType != "" {
		ct := reader.ContentType
		w.Header().Set("Content-Type", ct)
	}
	// Never let a browser sniff the bytes into a different, possibly executable
	// type than the one we stored (defense in depth alongside the raster-only
	// upload allowlist).
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.FormatInt(reader.Size, 10))
	_, _ = io.Copy(w, reader)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
