package handlers

import (
	"context"
	"time"

	"hexletbasics/internal/api"
	"hexletbasics/internal/assetstore"
)

// bookPath is the book page. Legacy routed it under the ru suffix only
// (RuSuffixConstraint), so there is no unprefixed or es variant to link to.
const bookPath = "/ru/book"

// bookDelivery is how legacy's send_file handed the PDF over: shown in the
// browser, saved under a readable name. The URL only has to outlive the click
// that follows the redirect.
var bookDelivery = assetstore.Delivery{
	ContentType:        "application/pdf",
	ContentDisposition: `inline; filename="profession-developer-hexlet-book.pdf"`,
	Expiry:             time.Hour,
}

// GetBook answers the book page (legacy `books#show`): whether the visitor has
// requested the book, in any state, which picks the request or the download
// button. A visitor has no request, so the page renders for them too.
func (s *Server) GetBook(ctx context.Context) (*api.BookView, error) {
	view := &api.BookView{}

	u, ok := AuthenticatedUser(ctx)
	if !ok {
		return view, nil
	}

	requested, err := s.books.Requested(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	view.Requested = requested

	return view, nil
}

// CreateBookRequest records the signed-in user's request (legacy
// `books#create_request`). A repeat succeeds and changes nothing, as legacy
// flashed the same success for it; only a first request raises BookRequested.
func (s *Server) CreateBookRequest(ctx context.Context) error {
	u := MustAuthenticatedUser(ctx)
	if err := s.books.Request(ctx, u.ID, s.i18n.Locale(ctx)); err != nil {
		return err
	}

	return nil
}

// DownloadBook sends the signed-in user to the PDF (legacy `books#download`).
// Legacy streamed the file; here the browser is redirected to the bucket so
// 28 MB never pass through the API (ADR-0005). The request is marked
// downloaded on every download, as legacy did, and a user who has not
// requested the book is sent back to the book page.
func (s *Server) DownloadBook(ctx context.Context) (*api.DownloadBookFound, error) {
	u := MustAuthenticatedUser(ctx)

	marked, err := s.books.MarkDownloaded(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	if !marked {
		return &api.DownloadBookFound{Location: s.settings.SiteURL + bookPath}, nil
	}

	location, err := s.assets.DownloadURL(ctx, s.settings.BookBlobKey, bookDelivery)
	if err != nil {
		return nil, err
	}

	return &api.DownloadBookFound{Location: location}, nil
}
