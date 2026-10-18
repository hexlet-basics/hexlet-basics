package oauthserver

import (
	"context"
	"errors"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"hexletbasics/ent"
	"hexletbasics/ent/oauthaccesstoken"
	"hexletbasics/ent/oauthclient"
	"hexletbasics/ent/oauthrefreshtoken"
)

// ErrAppNotConnected means the signed-in user has no Connected App with that
// id: it never existed, belongs to someone else, or has already lapsed or been
// revoked. One answer for all of them, so the id says nothing about other users.
var ErrAppNotConnected = errors.New("oauthserver: not a connected app of this user")

// ConnectedApp is one OAuth client a user has approved and that can still act
// for them (GLOSSARY: Connected App).
type ConnectedApp struct {
	// ID is the client row's id, the handle a revoke names. The client_id is not
	// used for it: a Client ID Metadata Document id is a URL.
	ID   int
	Name string
	// Scopes is what the newest live grant holds.
	Scopes []string
	// ConnectedAt is when the user approved the oldest grant to this client that
	// is still on file (a rotated token keeps the approval time of its family).
	ConnectedAt time.Time
	// ClientHost is the host of the client_id when it is a metadata document URL,
	// the one fact about such an app that its author does not choose through a
	// display name; empty for any other client.
	ClientHost string
}

// ConnectedApps lists the signed-in user's Connected Apps, newest connection
// first. An app is connected while the user holds a live grant to it: a refresh
// token that is unused and unexpired (what a client uses to keep access), or an
// unexpired access token (a grant made without offline_access, which lives out
// its hour). A used refresh token is a tombstone for replay detection, not a
// grant, so it never keeps an app on the list by itself; it does count towards
// ConnectedAt as part of the history of the connection.
//
// An app therefore disappears thirty days after its last refresh (the refresh
// lifetime), whether or not its MCP calls kept coming: "used" means refreshed.
func (s *Server) ConnectedApps(ctx context.Context) ([]ConnectedApp, error) {
	account, ok := s.signedIn(ctx)
	if !ok {
		return nil, ErrNotSignedIn
	}

	now := s.storage.now()
	db := s.storage.db

	refresh, err := db.OAuthRefreshToken.Query().
		Where(oauthrefreshtoken.UserID(account.ID), oauthrefreshtoken.ExpiresAtGT(now)).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list refresh tokens of user")
	}

	access, err := db.OAuthAccessToken.Query().
		Where(oauthaccesstoken.UserID(account.ID), oauthaccesstoken.ExpiresAtGT(now)).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "list access tokens of user")
	}

	type grants struct {
		live        bool
		connectedAt time.Time
		newest      time.Time
		scopes      []string
	}

	byClient := map[string]*grants{}
	// Times come from the token's own columns, never from created_at: that one
	// is stamped by the process clock (and without a zone), while expiry and
	// auth_time follow the injected clock the rest of the package runs on.
	see := func(clientID string, connected, issued time.Time, scopes []string, live bool) {
		g := byClient[clientID]
		if g == nil {
			g = &grants{connectedAt: connected}
			byClient[clientID] = g
		}

		if connected.Before(g.connectedAt) {
			g.connectedAt = connected
		}

		if live {
			g.live = true
			if issued.After(g.newest) {
				g.newest, g.scopes = issued, scopes
			}
		}
	}
	for _, r := range refresh {
		// auth_time is when the user approved; every rotation inherits it.
		see(r.ClientID, r.AuthTime, r.ExpiresAt.Add(-s.storage.refreshLifetime), r.Scopes, r.UsedAt == nil)
	}

	for _, a := range access {
		issued := a.ExpiresAt.Add(-s.storage.accessLifetime)
		see(a.ClientID, issued, issued, a.Scopes, true)
	}

	clientIDs := make([]string, 0, len(byClient))
	for id, g := range byClient {
		if g.live {
			clientIDs = append(clientIDs, id)
		}
	}

	if len(clientIDs) == 0 {
		return []ConnectedApp{}, nil
	}

	clients, err := db.OAuthClient.Query().Where(oauthclient.ClientIDIn(clientIDs...)).All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "load clients of connected apps")
	}

	apps := make([]ConnectedApp, 0, len(clients))
	for _, c := range clients {
		g := byClient[c.ClientID]
		apps = append(apps, ConnectedApp{
			ID: c.ID, Name: c.Name, Scopes: slices.Clone(g.scopes),
			ConnectedAt: g.connectedAt, ClientHost: documentHost(c.ClientID),
		})
	}

	slices.SortFunc(apps, func(a, b ConnectedApp) int {
		if c := b.ConnectedAt.Compare(a.ConnectedAt); c != 0 {
			return c
		}

		return b.ID - a.ID
	})

	return apps, nil
}

// documentHost is the host of a client_id that is an https URL (a Client ID
// Metadata Document), and empty for an opaque id.
func documentHost(clientID string) string {
	if !strings.HasPrefix(clientID, "https://") {
		return ""
	}

	u, err := url.Parse(clientID)
	if err != nil {
		return ""
	}

	return u.Host
}

// Disconnect revokes one Connected App of the signed-in user: every access
// token and every refresh token the user holds for that client, in every family
// and used tombstones included, goes in one transaction. It is the same
// primitive as terminating a session, so the effect is immediate: the access
// token is refused at its next use, and the refresh token at its next
// presentation. Other users' grants to the same client are untouched.
//
// An id that is not one of this user's Connected Apps is ErrAppNotConnected.
func (s *Server) Disconnect(ctx context.Context, appID int) error {
	account, ok := s.signedIn(ctx)
	if !ok {
		return ErrNotSignedIn
	}

	apps, err := s.ConnectedApps(ctx)
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(apps, func(a ConnectedApp) bool { return a.ID == appID }) {
		return ErrAppNotConnected
	}

	row, err := s.storage.db.OAuthClient.Get(ctx, appID)
	if ent.IsNotFound(err) {
		return ErrAppNotConnected
	}

	if err != nil {
		return errs.Wrapf(err, "load client of connected app")
	}

	return s.storage.TerminateSession(ctx, strconv.Itoa(account.ID), row.ClientID)
}
