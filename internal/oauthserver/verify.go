package oauthserver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"hexletbasics/ent"
	"hexletbasics/ent/oauthaccesstoken"
)

// ErrInvalidToken is the one answer for every token that must not be accepted:
// not ours, tampered with, unknown, revoked or expired. A caller learns nothing
// about which.
var ErrInvalidToken = errors.New("oauthserver: invalid access token")

// Token is a verified access token.
type Token struct {
	ID       string
	ClientID string
	UserID   int
	Scopes   []string
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
	return Token{ID: row.ID, ClientID: row.ClientID, UserID: userID, Scopes: row.Scopes}, nil
}

// liveAccessToken loads an access token that is on file, unexpired and issued to
// the subject; every other case is ErrInvalidToken.
func (s *storage) liveAccessToken(ctx context.Context, id, subject string) (*ent.OAuthAccessToken, error) {
	userID, err := strconv.Atoi(subject)
	if err != nil {
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
