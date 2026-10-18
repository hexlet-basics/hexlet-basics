package oauthserver

import (
	"context"
	"errors"
	"net/url"

	"hexletbasics/ent"
	"hexletbasics/ent/oauthauthrequest"
	"hexletbasics/ent/oauthclient"
)

var (
	// ErrNotSignedIn means the consent operation arrived without a session.
	ErrNotSignedIn = errors.New("oauthserver: no signed-in user")
	// ErrRequestNotFound covers every way an authorization request can be
	// unusable: never existed, expired, or already decided. The consent screen
	// shows one message for all of them, and the distinction would only tell a
	// guesser which ids exist.
	ErrRequestNotFound = errors.New("oauthserver: authorization request not found")
)

// Consent is what the consent screen shows: who is asking, where they will be
// sent back to, what they ask for, and which account would grant it.
type Consent struct {
	ClientName string
	// RedirectHost is the host the user returns to after deciding, so a
	// look-alike application name still shows its real destination.
	RedirectHost string
	Scopes       []string
	Account      Account
}

// Pending describes the authorization request the browser arrived with.
func (s *Server) Pending(ctx context.Context, requestID string) (Consent, error) {
	account, ok := s.app.SignedIn(ctx)
	if !ok {
		return Consent{}, ErrNotSignedIn
	}
	req, err := s.openRequest(ctx, requestID)
	if err != nil {
		return Consent{}, err
	}
	app, err := s.storage.db.OAuthClient.Query().Where(oauthclient.ClientID(req.ClientID)).Only(ctx)
	if err != nil {
		return Consent{}, errs.Wrapf(err, "load client of authorization request")
	}
	host := ""
	if u, err := url.Parse(req.RedirectURI); err == nil {
		host = u.Host
	}
	return Consent{ClientName: app.Name, RedirectHost: host, Scopes: req.Scopes, Account: account}, nil
}

// Approve records that the signed-in user grants the request and returns where
// the browser goes next: the authorization endpoint's callback, which issues the
// code and sends the user back to the client. The callback URL is built from the
// configured issuer, not from the incoming request.
func (s *Server) Approve(ctx context.Context, requestID string) (string, error) {
	account, ok := s.app.SignedIn(ctx)
	if !ok {
		return "", ErrNotSignedIn
	}
	now := s.storage.now()
	// One conditional update decides it: a request is approved at most once, and
	// never after it expired, whatever the two browser tabs do.
	n, err := s.storage.db.OAuthAuthRequest.Update().
		Where(
			oauthauthrequest.ID(requestID),
			oauthauthrequest.DoneEQ(false),
			oauthauthrequest.ExpiresAtGT(now),
		).
		SetDone(true).
		SetUserID(account.ID).
		SetAuthTime(now).
		Save(ctx)
	if err != nil {
		return "", errs.Wrapf(err, "approve authorization request")
	}
	if n == 0 {
		return "", ErrRequestNotFound
	}
	return s.issuer + callbackPath + "?id=" + url.QueryEscape(requestID), nil
}

// Deny refuses the request and returns the client's redirect URI carrying
// error=access_denied and the client's state (RFC 6749 section 4.1.2.1). The
// request is deleted, so nothing it held can be approved afterwards.
func (s *Server) Deny(ctx context.Context, requestID string) (string, error) {
	if _, ok := s.app.SignedIn(ctx); !ok {
		return "", ErrNotSignedIn
	}
	req, err := s.openRequest(ctx, requestID)
	if err != nil {
		return "", err
	}
	if err := s.storage.DeleteAuthRequest(ctx, requestID); err != nil {
		return "", err
	}
	target, err := url.Parse(req.RedirectURI)
	if err != nil {
		return "", errs.Wrapf(err, "parse redirect uri of authorization request")
	}
	query := target.Query()
	query.Set("error", "access_denied")
	query.Set("error_description", "The user denied the request.")
	query.Set("iss", s.issuer)
	if req.State != "" {
		query.Set("state", req.State)
	}
	target.RawQuery = query.Encode()
	return target.String(), nil
}

// openRequest loads a request that is still awaiting a decision.
func (s *Server) openRequest(ctx context.Context, requestID string) (*ent.OAuthAuthRequest, error) {
	req, err := s.storage.db.OAuthAuthRequest.Query().
		Where(
			oauthauthrequest.ID(requestID),
			oauthauthrequest.DoneEQ(false),
			oauthauthrequest.ExpiresAtGT(s.storage.now()),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrRequestNotFound
	}
	if err != nil {
		return nil, errs.Wrapf(err, "load authorization request")
	}
	return req, nil
}
