package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
	"strconv"

	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"

	"hexletbasics/ent"
	"hexletbasics/internal/ids"
)

// maxUploadBytes caps a single upload. Covers/outcomes images are small; the
// limit stops an oversized body from being buffered to disk. It doubles as the
// http.MaxBytesReader ceiling so the read aborts rather than filling storage.
const maxUploadBytes = 15 << 20 // 15 MiB

// allowedUploadTypes is the content-type allowlist for uploads. The admin
// uploader only ever supplies cover/outcomes images (legacy covers are raster,
// e.g. 800x400), so restrict to raster images — this is the sole validation
// guarding a public-ish write endpoint until auth lands, and it keeps arbitrary
// blobs out of the bucket. SVG is deliberately excluded: Download echoes the
// stored content type same-origin as the API (which sets the JWT cookie), so a
// script-bearing `image/svg+xml` would be a stored-XSS vector.
var allowedUploadTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

// AttachmentHandler serves the multipart upload + blob read path that lives
// OUTSIDE the ogen-generated router. ogen cannot generate multipart/form-data
// (see internal/apigen/ogen.yml), so `POST /admin/attachments` and the
// `GET /storage/{key}` read-back are plain net/http handlers mounted alongside
// the generated api.Server by NewRouter. It owns the ent client (for the
// attachment record) and the gocloud bucket (for the bytes).
type AttachmentHandler struct {
	db     *ent.Client
	bucket *blob.Bucket
}

// NewAttachmentHandler wires the uploader to its ent client and blob bucket.
func NewAttachmentHandler(db *ent.Client, bucket *blob.Bucket) *AttachmentHandler {
	return &AttachmentHandler{db: db, bucket: bucket}
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
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

	file, header, err := r.FormFile("file")
	if err != nil {
		// A body over the limit surfaces here as *http.MaxBytesError.
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			writeValidationError(w, "file", "file is too large")
			return
		}
		writeValidationError(w, "file", "a file part is required")
		return
	}
	defer func() { _ = file.Close() }()

	contentType := header.Header.Get("Content-Type")
	if !allowedUploadTypes[contentType] {
		writeValidationError(w, "file", "unsupported file type")
		return
	}

	// Preserve the original extension so the read URL is recognizable / servable;
	// the ULID keeps keys unique and unguessable-enough without a DB round-trip.
	key := ids.New() + path.Ext(header.Filename)

	written, err := h.writeBlob(ctx, key, contentType, file)
	if err != nil {
		// Storage failure is a server error, not a validation error.
		http.Error(w, "failed to store file", http.StatusInternalServerError)
		return
	}

	att, err := h.db.Attachment.Create().
		SetStorageKey(key).
		SetFilename(header.Filename).
		SetContentType(contentType).
		SetByteSize(written).
		Save(ctx)
	if err != nil {
		// Best-effort cleanup: the bytes are already in the bucket but the record
		// failed, so drop the orphan rather than leak it.
		_ = h.bucket.Delete(ctx, key)
		http.Error(w, "failed to record attachment", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, attachmentResponse{
		ID:          att.ID,
		URL:         attachmentURL(r, key),
		Filename:    att.Filename,
		ContentType: att.ContentType,
		ByteSize:    att.ByteSize,
	})
}

// writeBlob streams the uploaded reader into the bucket and returns the byte
// count actually written (so byte_size reflects the stored bytes, not a
// client-declared size). The writer's Close is where fileblob/s3blob surface
// most errors, so it is checked.
func (h *AttachmentHandler) writeBlob(ctx context.Context, key, contentType string, src io.Reader) (int64, error) {
	bw, err := h.bucket.NewWriter(ctx, key, &blob.WriterOptions{ContentType: contentType})
	if err != nil {
		return 0, err
	}
	written, copyErr := io.Copy(bw, src)
	closeErr := bw.Close()
	if copyErr != nil {
		return 0, copyErr
	}
	if closeErr != nil {
		return 0, closeErr
	}
	return written, nil
}

// Download handles `GET /storage/{key}`: stream the stored bytes back with the
// content type they were written with. This is the read side of the `url` the
// uploader returns — without it that url would point at nothing.
func (h *AttachmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := r.PathValue("key")

	reader, err := h.bucket.NewReader(ctx, key, nil)
	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}
	defer func() { _ = reader.Close() }()

	if ct := reader.ContentType(); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	// Never let a browser sniff the bytes into a different, possibly executable
	// type than the one we stored (defense in depth alongside the raster-only
	// upload allowlist).
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.FormatInt(reader.Size(), 10))
	_, _ = io.Copy(w, reader)
}

// attachmentURL is the ABSOLUTE read URL for a stored key, built from the request
// origin. It must be absolute: the frontend runs on a different origin than the
// API (dev Vite :3000 vs API :3001; VITE_API_URL in prod), so a root-relative
// `/storage/...` in an <img src> would resolve against the frontend and 404. The
// upload is always reached through the API's own public origin, so the request's
// scheme+host is the correct base for the read-back URL.
func attachmentURL(r *http.Request, key string) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host + "/storage/" + key
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeValidationError emits the OpenAPI ValidationError shape (422, errors keyed
// by field) so a failed upload is indistinguishable on the wire from the ogen
// writes' 422s the frontend already handles.
func writeValidationError(w http.ResponseWriter, field, message string) {
	writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
		"errors": map[string][]string{field: {message}},
	})
}
