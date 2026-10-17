package oauthserver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"hexletbasics/ent"
	"hexletbasics/ent/oauthaccesstoken"
	"hexletbasics/ent/oauthclient"
)

// ErrInvalidToken is the one answer for every token that must not be accepted:
// not ours, tampered with, unknown, revoked or expired. A caller learns nothing
// about which.
var ErrInvalidToken = errors.New("oauthserver: invalid access token")

// Token is a verified access token.
type Token struct {
	ID       string
	ClientID string
	// ConnectedAppID is the id of the oauth_clients row behind ClientID: the
	// Connected App the Audit Record names (ADR-0020). The public client_id is
	// the application's own label, the row id is what the trail keeps.
	ConnectedAppID int
	UserID         int
	Scopes         []string
	// Audience is the resource the token is bound to: the canonical MCP URL. A
	// resource server must check that it is itself in it. It is not stored per
	// token, since this server has exactly one resource and admits no other
	// (RFC 8707), so every token on file is for it.
	Audience []string
}

// VerifyAccessToken is how a resource server (the MCP endpoint) accepts a bearer
// token. The token the client holds is the encrypted pair (id, subject); the
// database holds only the id. So verification decrypts, then looks the id up, and
// accepts only a row that is unexpired and belongs to that subject: deleting the
// row revokes the token at once, and a database leak yields ids that decrypt to
// nothing and cannot be presented as tokens.
func (s *Server) VerifyAccessToken(ctx context.Context, bearer string) (Token, error) {
	plain, err := s.provider.Crypto().Decrypt(bearer)
	if err != nil {
		return Token{}, ErrInvalidToken
	}
	id, subject, ok := strings.Cut(plain, ":")
	if !ok {
		return Token{}, ErrInvalidToken
	}
	userID, err := strconv.Atoi(subject)
	if err != nil {
		return Token{}, ErrInvalidToken
	}
	row, err := s.storage.liveAccessToken(ctx, id, subject)
	if err != nil {
		return Token{}, err
	}
	// A token whose client was removed is no longer a Connected App's: refused,
	// the same as a revoked token, rather than recorded against nothing.
	app, err := s.storage.db.OAuthClient.Query().Where(oauthclient.ClientID(row.ClientID)).OnlyID(ctx)
	if ent.IsNotFound(err) {
		return Token{}, ErrInvalidToken
	}
	if err != nil {
		return Token{}, fmt.Errorf("look up connected app: %w", err)
	}
	return Token{
		ID: row.ID, ClientID: row.ClientID, ConnectedAppID: app, UserID: userID,
		Scopes: row.Scopes, Audience: []string{s.resource},
	}, nil
}

// liveAccessToken loads an access token that is on file, unexpired and issued to
// the subject; every other case is ErrInvalidToken.
func (s *storage) liveAccessToken(ctx context.Context, id, subject string) (*ent.OAuthAccessToken, error) {
	userID, err := strconv.Atoi(subject)
	if err != nil || !validTokenID(id) {
		return nil, ErrInvalidToken
	}
	row, err := s.db.OAuthAccessToken.Query().
		Where(
			oauthaccesstoken.ID(id),
			oauthaccesstoken.UserID(userID),
			oauthaccesstoken.ExpiresAtGT(s.now()),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("look up access token: %w", err)
	}
	return row, nil
}
