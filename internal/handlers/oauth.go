package handlers

import (
	"context"
	"errors"
	"strings"

	"hexletbasics/internal/api"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/oauthserver"
)

// OAuthConsents is the consent seam of the OAuth server (ADR-0022): what the
// consent screen shows, and approving or denying it. *oauthserver.Server is the
// implementation; the dependency points this way only, handlers know the
// package and the package never knows handlers. nil means the OAuth server is
// switched off, and the operations answer as if the request did not exist.
type OAuthConsents interface {
	Pending(ctx context.Context, requestID string) (oauthserver.Consent, error)
	Approve(ctx context.Context, requestID string) (redirectURL string, err error)
	Deny(ctx context.Context, requestID string) (redirectURL string, err error)
	// ConnectedApps and Disconnect are the Connected Apps page: what the signed-in
	// user has approved, and withdrawing one.
	ConnectedApps(ctx context.Context) ([]oauthserver.ConnectedApp, error)
	Disconnect(ctx context.Context, appID int) error
}

// OAuthSession is the application's half of the seam: who is signed in. The
// consent operations are secured by the contract, so the user is the one the
// session security handler already loaded into the context.
type OAuthSession struct{}

var _ oauthserver.Application = OAuthSession{}

// SignedIn names the authenticated user for the consent screen: the name when
// the profile has one, the email otherwise.
func (OAuthSession) SignedIn(ctx context.Context) (oauthserver.Account, bool) {
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		return oauthserver.Account{}, false
	}
	name := strings.TrimSpace(ptrValue(u.FirstName) + " " + ptrValue(u.LastName))
	if name == "" {
		name = ptrValue(u.Email)
	}
	return oauthserver.Account{ID: u.ID, Name: name}, true
}

func ptrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// GetOAuthConsent returns what the consent screen shows for the authorization
// request the browser arrived with.
func (s *Server) GetOAuthConsent(ctx context.Context, params api.GetOAuthConsentParams) (api.GetOAuthConsentRes, error) {
	if s.oauth == nil {
		return s.oauthNotFound(ctx), nil
	}
	consent, err := s.oauth.Pending(ctx, params.RequestId)
	if errors.Is(err, oauthserver.ErrRequestNotFound) {
		return s.oauthNotFound(ctx), nil
	}
	if err != nil {
		return nil, err
	}
	return &api.OAuthConsent{
		ClientName:   consent.ClientName,
		RedirectHost: consent.RedirectHost,
		Scopes:       consent.Scopes,
		AccountName:  consent.Account.Name,
	}, nil
}

// ApproveOAuthConsent grants the request as the signed-in user.
func (s *Server) ApproveOAuthConsent(ctx context.Context, params api.ApproveOAuthConsentParams) (api.ApproveOAuthConsentRes, error) {
	return s.decideOAuthConsent(ctx, params.RequestId, s.oauthApprove)
}

// DenyOAuthConsent refuses the request.
func (s *Server) DenyOAuthConsent(ctx context.Context, params api.DenyOAuthConsentParams) (api.DenyOAuthConsentRes, error) {
	return s.decideOAuthConsent(ctx, params.RequestId, s.oauthDeny)
}

func (s *Server) oauthApprove(ctx context.Context, id string) (string, error) {
	return s.oauth.Approve(ctx, id)
}

func (s *Server) oauthDeny(ctx context.Context, id string) (string, error) {
	return s.oauth.Deny(ctx, id)
}

// decideOAuthConsent runs one of the two decisions. Both answer the same shape,
// so one function maps both, and the typed unions of the two operations accept
// the same three result types.
func (s *Server) decideOAuthConsent(
	ctx context.Context, requestID string, decide func(context.Context, string) (string, error),
) (decision, error) {
	if s.oauth == nil {
		return s.oauthNotFound(ctx), nil
	}
	redirect, err := decide(ctx, requestID)
	if errors.Is(err, oauthserver.ErrRequestNotFound) {
		return s.oauthNotFound(ctx), nil
	}
	if err != nil {
		return nil, err
	}
	return &api.OAuthDecision{RedirectUrl: redirect}, nil
}

// decision is the result both consent decisions share: the generated
// per-operation response unions have the same members.
type decision interface {
	api.ApproveOAuthConsentRes
	api.DenyOAuthConsentRes
}

func (s *Server) oauthNotFound(ctx context.Context) *api.NotFoundError {
	return &api.NotFoundError{Message: s.i18n.Text(ctx, localization.OAuthRequestNotFound)}
}

// ListConnectedApps returns the signed-in user's Connected Apps. With the OAuth
// server off no app can exist, so the list is empty rather than a 404: the
// account page then shows its empty state instead of an error.
func (s *Server) ListConnectedApps(ctx context.Context) (api.ListConnectedAppsRes, error) {
	if s.oauth == nil {
		return &api.ConnectedApps{Items: []api.ConnectedApp{}}, nil
	}
	apps, err := s.oauth.ConnectedApps(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]api.ConnectedApp, 0, len(apps))
	for _, a := range apps {
		item := api.ConnectedApp{
			ID:          int64(a.ID),
			Name:        a.Name,
			Scopes:      a.Scopes,
			ConnectedAt: a.ConnectedAt,
		}
		if a.ClientHost != "" {
			item.ClientHost = api.NewOptString(a.ClientHost)
		}
		items = append(items, item)
	}
	return &api.ConnectedApps{Items: items}, nil
}

// RevokeConnectedApp ends every token the signed-in user granted one app. An
// app that is not theirs, has lapsed or was revoked already is a 404, the same
// answer for each so the id reveals nothing about other users.
func (s *Server) RevokeConnectedApp(ctx context.Context, params api.RevokeConnectedAppParams) (api.RevokeConnectedAppRes, error) {
	notFound := &api.NotFoundError{Message: s.i18n.Text(ctx, localization.ConnectedAppNotFound)}
	if s.oauth == nil {
		return notFound, nil
	}
	err := s.oauth.Disconnect(ctx, int(params.ID))
	if errors.Is(err, oauthserver.ErrAppNotConnected) {
		return notFound, nil
	}
	if err != nil {
		return nil, err
	}
	return &api.RevokeConnectedAppNoContent{}, nil
}
