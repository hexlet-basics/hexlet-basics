package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"hexletbasics/internal/api"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/localization"
)

// AdminUploadAttachment implements `POST /admin/attachments`: store the single
// multipart `file` part in the bucket under a fresh key, record its metadata,
// and return the Attachment the admin form references by id. ogen decodes the
// form and enforces admin+XSRF before this runs; the router caps the body size.
// The part's declared Content-Type is ignored on purpose — assetstore derives
// the stored type from the bytes.
func (s *Server) AdminUploadAttachment(
	ctx context.Context,
	req *api.AttachmentUploadFormMultipart,
) (api.AdminUploadAttachmentRes, error) {
	att, err := s.assets.Put(ctx, assetstore.Upload{
		Filename: req.File.Name,
		Body:     req.File.File,
	})
	switch {
	case errors.Is(err, assetstore.ErrUnsupportedMediaType):
		return validationError("file", s.i18n.Text(ctx, localization.UnsupportedFileType)), nil
	case errors.Is(err, assetstore.ErrTooLarge):
		return validationError("file", s.i18n.Text(ctx, localization.FileTooLarge)), nil
	case err != nil:
		return nil, fmt.Errorf("store uploaded attachment: %w", err)
	}
	attachment := s.conv.ToAttachment(att)
	return &attachment, nil
}

// AttachmentHandler serves the blob read path, `GET /api/storage/{key}`, mounted
// alongside the generated api.Server by NewRouter. It is not a contract
// operation: it streams stored bytes with standard conditional and range
// semantics, which http.ServeContent provides and a JSON contract cannot
// describe. Uploads are the generated adminUploadAttachment operation, and
// asset lifecycle invariants live in assetstore.Store rather than here.
type AttachmentHandler struct {
	assets *assetstore.Store
	errors *APIErrorHandler
}

// NewAttachmentHandler wires the blob read route to the shared asset store.
func NewAttachmentHandler(assets *assetstore.Store, errorHandler *APIErrorHandler) *AttachmentHandler {
	return &AttachmentHandler{assets: assets, errors: errorHandler}
}

// Download handles `GET /api/storage/{key}`: serve the stored bytes with standard
// conditional and range request semantics. This is the read side of the `url`
// the uploader returns — without it that url would point at nothing.
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
	http.ServeContent(w, r, key, reader.ModTime, reader)
}
