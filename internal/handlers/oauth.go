package handlers

import (
	"context"
	"errors"
	"strings"

	"hexletbasics/internal/api"
	"hexletbasics/internal/oauthserver"
)

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
func (s *Server) GetOAuthConsent(ctx context.Context, params api.GetOAuthConsentParams) (*api.OAuthConsent, error) {
	if s.oauth == nil {
		return nil, errOAuthRequestNotFound
	}
	consent, err := s.oauth.Pending(ctx, params.RequestId)
	if errors.Is(err, oauthserver.ErrRequestNotFound) {
		return nil, errOAuthRequestNotFound
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
func (s *Server) ApproveOAuthConsent(ctx context.Context, params api.ApproveOAuthConsentParams) (*api.OAuthDecision, error) {
	return s.decideOAuthConsent(ctx, params.RequestId, s.oauthApprove)
}

// DenyOAuthConsent refuses the request.
func (s *Server) DenyOAuthConsent(ctx context.Context, params api.DenyOAuthConsentParams) (*api.OAuthDecision, error) {
	return s.decideOAuthConsent(ctx, params.RequestId, s.oauthDeny)
}

func (s *Server) oauthApprove(ctx context.Context, id string) (string, error) {
	return s.oauth.Approve(ctx, id)
}

func (s *Server) oauthDeny(ctx context.Context, id string) (string, error) {
	return s.oauth.Deny(ctx, id)
}

// decideOAuthConsent runs one of the two decisions. Both answer the same shape,
// so one function maps both.
func (s *Server) decideOAuthConsent(
	ctx context.Context, requestID string, decide func(context.Context, string) (string, error),
) (*api.OAuthDecision, error) {
	if s.oauth == nil {
		return nil, errOAuthRequestNotFound
	}
	redirect, err := decide(ctx, requestID)
	if errors.Is(err, oauthserver.ErrRequestNotFound) {
		return nil, errOAuthRequestNotFound
	}
	if err != nil {
		return nil, err
	}
	return &api.OAuthDecision{RedirectUrl: redirect}, nil
}

// ListConnectedApps returns the signed-in user's Connected Apps. With the OAuth
// server off no app can exist, so the list is empty rather than a 404: the
// account page then shows its empty state instead of an error.
func (s *Server) ListConnectedApps(ctx context.Context) (*api.ConnectedApps, error) {
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
func (s *Server) RevokeConnectedApp(ctx context.Context, params api.RevokeConnectedAppParams) error {
	if s.oauth == nil {
		return errConnectedAppNotFound
	}
	err := s.oauth.Disconnect(ctx, int(params.ID))
	if errors.Is(err, oauthserver.ErrAppNotConnected) {
		return errConnectedAppNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
