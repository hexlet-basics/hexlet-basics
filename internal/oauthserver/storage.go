package oauthserver

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/lib/pq"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"hexletbasics/ent"
	"hexletbasics/ent/oauthaccesstoken"
	"hexletbasics/ent/oauthauthrequest"
	"hexletbasics/ent/oauthclient"
	"hexletbasics/ent/oauthrefreshtoken"
	"hexletbasics/internal/store"
)

// storage is the op.Storage behind the provider, on ent. Only the slice an
// authorization-code flow with refresh needs is real; the OpenID-only parts
// (user info, JWT profile, introspection) are deliberate stubs: the MCP endpoint
// verifies a token by looking it up (VerifyAccessToken), not through them.
type storage struct {
	db      *ent.Client
	tx      store.Transactor
	issuer  string
	now     func() time.Time
	signing ed25519.PrivateKey
}

var _ op.Storage = (*storage)(nil)

// errUnsupported marks the OpenID features this server does not offer.
var errUnsupported = errors.New("oauthserver: not supported")

// randomToken returns n random bytes as unpadded base64url: an unguessable
// credential or id.
func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// digest is what the database keeps in place of a code or refresh token: enough
// to recognise one that is presented, never enough to reproduce it. Both are
// 128-bit-plus random values, so a fast hash is sufficient.
func digest(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

// --- authorization requests ----------------------------------------------

// CreateAuthRequest stores the request the client sent the user over with and
// hands the library back the row it will redirect to the consent screen by.
//
// PKCE is enforced here, where the raw parameters are still visible: the library
// treats a missing challenge as acceptable and an empty method as `plain`, but
// S256 is mandatory for every client (ADR-0016).
func (s *storage) CreateAuthRequest(ctx context.Context, req *oidc.AuthRequest, _ string) (op.AuthRequest, error) {
	if req.CodeChallenge == "" || req.CodeChallengeMethod != oidc.CodeChallengeMethodS256 {
		return nil, oidc.ErrInvalidRequest().WithDescription("PKCE is required: send code_challenge with code_challenge_method=S256")
	}
	if len(req.Prompt) == 1 && req.Prompt[0] == oidc.PromptNone {
		// Nobody can sign in without a screen, so there is no silent grant.
		return nil, oidc.ErrLoginRequired()
	}
	id, err := randomToken(16)
	if err != nil {
		return nil, err
	}
	row, err := s.db.OAuthAuthRequest.Create().
		SetID(id).
		SetClientID(req.ClientID).
		SetRedirectURI(req.RedirectURI).
		SetState(req.State).
		SetScopes(pq.StringArray(req.Scopes)).
		SetCodeChallenge(req.CodeChallenge).
		SetCodeChallengeMethod(string(req.CodeChallengeMethod)).
		SetNonce(req.Nonce).
		SetResponseType(string(req.ResponseType)).
		SetExpiresAt(s.now().Add(authRequestLifetime)).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create authorization request: %w", err)
	}
	return &authRequest{row: row}, nil
}

// AuthRequestByID loads a request that has not expired.
func (s *storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	row, err := s.db.OAuthAuthRequest.Query().
		Where(oauthauthrequest.ID(id), oauthauthrequest.ExpiresAtGT(s.now())).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("authorization request: %w", err)
	}
	return &authRequest{row: row}, nil
}

// AuthRequestByCode finds the request an authorization code was issued for,
// by the code's hash.
func (s *storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	hash := digest(code)
	row, err := s.db.OAuthAuthRequest.Query().
		Where(oauthauthrequest.CodeHash(hash), oauthauthrequest.ExpiresAtGT(s.now())).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("authorization code: %w", err)
	}
	return &authRequest{row: row}, nil
}

// SaveAuthCode records the hash of the code issued for a request. A request
// yields one code: a second attempt to issue (the callback URL opened twice) is
// refused instead of replacing the first, so an id seen in a log cannot be used
// to mint a code after the user has been redirected.
func (s *storage) SaveAuthCode(ctx context.Context, id, code string) error {
	n, err := s.db.OAuthAuthRequest.Update().
		Where(oauthauthrequest.ID(id), oauthauthrequest.DoneEQ(true), oauthauthrequest.CodeHashIsNil()).
		SetCodeHash(digest(code)).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("save authorization code: %w", err)
	}
	if n == 0 {
		return errors.New("authorization request is not awaiting a code")
	}
	return nil
}

// DeleteAuthRequest consumes a request: after the token exchange its code is
// dead, which is what makes a code single-use.
func (s *storage) DeleteAuthRequest(ctx context.Context, id string) error {
	if _, err := s.db.OAuthAuthRequest.Delete().Where(oauthauthrequest.ID(id)).Exec(ctx); err != nil {
		return fmt.Errorf("delete authorization request: %w", err)
	}
	return nil
}

// authRequest adapts a stored request to op.AuthRequest.
type authRequest struct{ row *ent.OAuthAuthRequest }

var _ op.AuthRequest = (*authRequest)(nil)

func (a *authRequest) GetID() string       { return a.row.ID }
func (a *authRequest) GetACR() string      { return "" }
func (a *authRequest) GetAMR() []string    { return nil }
func (a *authRequest) GetClientID() string { return a.row.ClientID }

// GetAudience is the client for now; binding it to the MCP endpoint's canonical
// URL through the `resource` parameter is a later step of ADR-0016.
func (a *authRequest) GetAudience() []string { return []string{a.row.ClientID} }
func (a *authRequest) GetAuthTime() time.Time {
	if a.row.AuthTime == nil {
		return time.Time{}
	}
	return *a.row.AuthTime
}
func (a *authRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return &oidc.CodeChallenge{
		Challenge: a.row.CodeChallenge,
		Method:    oidc.CodeChallengeMethod(a.row.CodeChallengeMethod),
	}
}
func (a *authRequest) GetNonce() string       { return a.row.Nonce }
func (a *authRequest) GetRedirectURI() string { return a.row.RedirectURI }
func (a *authRequest) GetResponseType() oidc.ResponseType {
	return oidc.ResponseType(a.row.ResponseType)
}
func (a *authRequest) GetResponseMode() oidc.ResponseMode { return "" }
func (a *authRequest) GetScopes() []string                { return a.row.Scopes }
func (a *authRequest) GetState() string                   { return a.row.State }
func (a *authRequest) Done() bool                         { return a.row.Done }
func (a *authRequest) GetSubject() string {
	if a.row.UserID == nil {
		return ""
	}
	return strconv.Itoa(*a.row.UserID)
}

// --- tokens ---------------------------------------------------------------

// tokenOwner is what both an authorization request and a refresh request know.
type tokenOwner interface {
	op.TokenRequest
	GetClientID() string
	GetAuthTime() time.Time
}

func userIDOf(req op.TokenRequest) (int, error) {
	id, err := strconv.Atoi(req.GetSubject())
	if err != nil {
		return 0, fmt.Errorf("token subject %q is not a user: %w", req.GetSubject(), err)
	}
	return id, nil
}

// CreateAccessToken records an access token for the code exchange. The value
// the client receives is built by the library from the returned id.
func (s *storage) CreateAccessToken(ctx context.Context, req op.TokenRequest) (string, time.Time, error) {
	owner, ok := req.(tokenOwner)
	if !ok {
		return "", time.Time{}, fmt.Errorf("unsupported token request %T", req)
	}
	userID, err := userIDOf(req)
	if err != nil {
		return "", time.Time{}, err
	}
	id, expires, err := s.insertAccessToken(ctx, s.db, owner, userID, "")
	return id, expires, err
}

func (s *storage) insertAccessToken(
	ctx context.Context, db *ent.Client, owner tokenOwner, userID int, refreshID string,
) (string, time.Time, error) {
	id, err := randomToken(16)
	if err != nil {
		return "", time.Time{}, err
	}
	expires := s.now().Add(accessTokenLifetime)
	create := db.OAuthAccessToken.Create().
		SetID(id).
		SetClientID(owner.GetClientID()).
		SetUserID(userID).
		SetScopes(pq.StringArray(owner.GetScopes())).
		SetExpiresAt(expires)
	if refreshID != "" {
		create.SetRefreshTokenID(refreshID)
	}
	if err := create.Exec(ctx); err != nil {
		return "", time.Time{}, fmt.Errorf("create access token: %w", err)
	}
	return id, expires, nil
}

// CreateAccessAndRefreshTokens issues a pair: for a first grant (no current
// refresh token) or in exchange for a refresh token, which is spent and replaced.
// Detecting a spent token being presented again is a later step (ADR-0016); here
// the old token is simply gone, so a replay is refused as unknown.
func (s *storage) CreateAccessAndRefreshTokens(
	ctx context.Context, req op.TokenRequest, current string,
) (string, string, time.Time, error) {
	owner, ok := req.(tokenOwner)
	if !ok {
		return "", "", time.Time{}, fmt.Errorf("unsupported token request %T", req)
	}
	userID, err := userIDOf(req)
	if err != nil {
		return "", "", time.Time{}, err
	}

	var accessID, refresh string
	var expires time.Time
	err = s.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		if current != "" {
			// Spend the presented token and the access token it was paired with.
			old, err := db.OAuthRefreshToken.Query().
				Where(oauthrefreshtoken.TokenHash(digest(current)), oauthrefreshtoken.ExpiresAtGT(s.now())).
				Only(ctx)
			if err != nil {
				return fmt.Errorf("refresh token: %w", err)
			}
			if _, err := db.OAuthAccessToken.Delete().
				Where(oauthaccesstoken.RefreshTokenID(old.ID)).Exec(ctx); err != nil {
				return fmt.Errorf("delete paired access token: %w", err)
			}
			if err := db.OAuthRefreshToken.DeleteOne(old).Exec(ctx); err != nil {
				return fmt.Errorf("delete spent refresh token: %w", err)
			}
		}

		refreshID, err := randomToken(16)
		if err != nil {
			return err
		}
		refresh, err = randomToken(32)
		if err != nil {
			return err
		}
		accessID, expires, err = s.insertAccessToken(ctx, db, owner, userID, refreshID)
		if err != nil {
			return err
		}
		authTime := owner.GetAuthTime()
		if authTime.IsZero() {
			authTime = s.now()
		}
		return db.OAuthRefreshToken.Create().
			SetID(refreshID).
			SetTokenHash(digest(refresh)).
			SetClientID(owner.GetClientID()).
			SetUserID(userID).
			SetScopes(pq.StringArray(owner.GetScopes())).
			SetAuthTime(authTime).
			SetAccessTokenID(accessID).
			SetExpiresAt(s.now().Add(refreshTokenLifetime)).
			Exec(ctx)
	})
	if err != nil {
		return "", "", time.Time{}, err
	}
	return accessID, refresh, expires, nil
}

// TokenRequestByRefreshToken finds the grant a presented refresh token belongs to.
func (s *storage) TokenRequestByRefreshToken(ctx context.Context, refresh string) (op.RefreshTokenRequest, error) {
	row, err := s.db.OAuthRefreshToken.Query().
		Where(oauthrefreshtoken.TokenHash(digest(refresh)), oauthrefreshtoken.ExpiresAtGT(s.now())).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}
	return &refreshRequest{row: row, scopes: row.Scopes}, nil
}

// refreshRequest adapts a stored refresh token to op.RefreshTokenRequest.
type refreshRequest struct {
	row    *ent.OAuthRefreshToken
	scopes []string
}

var _ op.RefreshTokenRequest = (*refreshRequest)(nil)

func (r *refreshRequest) GetAMR() []string            { return nil }
func (r *refreshRequest) GetAudience() []string       { return []string{r.row.ClientID} }
func (r *refreshRequest) GetAuthTime() time.Time      { return r.row.AuthTime }
func (r *refreshRequest) GetClientID() string         { return r.row.ClientID }
func (r *refreshRequest) GetScopes() []string         { return r.scopes }
func (r *refreshRequest) GetSubject() string          { return strconv.Itoa(r.row.UserID) }
func (r *refreshRequest) SetCurrentScopes(s []string) { r.scopes = s }

// GetRefreshTokenInfo lets the revocation endpoint tell a refresh token from an
// access token.
func (s *storage) GetRefreshTokenInfo(ctx context.Context, _ string, token string) (string, string, error) {
	row, err := s.db.OAuthRefreshToken.Query().
		Where(oauthrefreshtoken.TokenHash(digest(token))).Only(ctx)
	if ent.IsNotFound(err) {
		return "", "", op.ErrInvalidRefreshToken
	}
	if err != nil {
		return "", "", fmt.Errorf("refresh token info: %w", err)
	}
	return strconv.Itoa(row.UserID), row.ID, nil
}

// RevokeToken deletes the token with this id, an access token or a refresh token
// (with the access token paired to it). A token that is already gone is not an
// error: the caller wanted it dead and it is. Only the client the token was
// issued to may revoke it.
func (s *storage) RevokeToken(ctx context.Context, id, _, clientID string) *oidc.Error {
	err := s.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		access, err := db.OAuthAccessToken.Get(ctx, id)
		switch {
		case err == nil:
			if access.ClientID != clientID {
				return errWrongClient
			}
			return db.OAuthAccessToken.DeleteOne(access).Exec(ctx)
		case !ent.IsNotFound(err):
			return err
		}

		refresh, err := db.OAuthRefreshToken.Get(ctx, id)
		if ent.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if refresh.ClientID != clientID {
			return errWrongClient
		}
		if _, err := db.OAuthAccessToken.Delete().
			Where(oauthaccesstoken.RefreshTokenID(refresh.ID)).Exec(ctx); err != nil {
			return err
		}
		return db.OAuthRefreshToken.DeleteOne(refresh).Exec(ctx)
	})
	switch {
	case err == nil:
		return nil
	case errors.Is(err, errWrongClient):
		return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
	default:
		return oidc.ErrServerError().WithParent(err)
	}
}

var errWrongClient = errors.New("token belongs to another client")

// TerminateSession ends every grant a user made to one client.
func (s *storage) TerminateSession(ctx context.Context, userID, clientID string) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("session subject %q is not a user: %w", userID, err)
	}
	return s.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		if _, err := db.OAuthAccessToken.Delete().
			Where(oauthaccesstoken.UserID(id), oauthaccesstoken.ClientID(clientID)).Exec(ctx); err != nil {
			return err
		}
		_, err := db.OAuthRefreshToken.Delete().
			Where(oauthrefreshtoken.UserID(id), oauthrefreshtoken.ClientID(clientID)).Exec(ctx)
		return err
	})
}

// --- keys ------------------------------------------------------------------

type signingKey struct{ key ed25519.PrivateKey }

func (k signingKey) SignatureAlgorithm() jose.SignatureAlgorithm { return jose.EdDSA }
func (k signingKey) Key() any                                    { return k.key }
func (k signingKey) ID() string                                  { return keyID(k.key) }

type publicKey struct{ key ed25519.PublicKey }

func (k publicKey) ID() string                       { return keyID(k.key) }
func (publicKey) Algorithm() jose.SignatureAlgorithm { return jose.EdDSA }
func (publicKey) Use() string                        { return "sig" }
func (k publicKey) Key() any                         { return k.key }

// keyID names a key by a hash of its public half, so a rotation (a new
// encryption secret) shows up as a new id in the key set.
func keyID(key any) string {
	var pub ed25519.PublicKey
	switch k := key.(type) {
	case ed25519.PrivateKey:
		pub = k.Public().(ed25519.PublicKey)
	case ed25519.PublicKey:
		pub = k
	}
	return digest(string(pub))[:16]
}

func (s *storage) SigningKey(context.Context) (op.SigningKey, error) {
	return signingKey{key: s.signing}, nil
}

func (s *storage) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	return []jose.SignatureAlgorithm{jose.EdDSA}, nil
}

func (s *storage) KeySet(context.Context) ([]op.Key, error) {
	return []op.Key{publicKey{key: s.signing.Public().(ed25519.PublicKey)}}, nil
}

// --- clients ---------------------------------------------------------------

// GetClientByClientID loads a client for the current request.
func (s *storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	row, err := s.db.OAuthClient.Query().Where(oauthclient.ClientID(clientID)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("client %q: %w", clientID, err)
	}
	return &client{row: row, loginBase: s.issuer + consentPagePath}, nil
}

// AuthorizeClientIDSecret always refuses: every client is public.
func (s *storage) AuthorizeClientIDSecret(context.Context, string, string) error {
	return errors.New("clients have no secret: authenticate with PKCE")
}

// --- OpenID-only parts, stubbed --------------------------------------------

func (s *storage) Health(context.Context) error { return nil }

func (s *storage) SetUserinfoFromScopes(context.Context, *oidc.UserInfo, string, string, []string) error {
	return nil
}

// SetUserinfoFromToken answers the userinfo endpoint with the subject only,
// for a token that is still on file.
func (s *storage) SetUserinfoFromToken(ctx context.Context, info *oidc.UserInfo, tokenID, subject, _ string) error {
	if _, err := s.liveAccessToken(ctx, tokenID, subject); err != nil {
		return err
	}
	info.Subject = subject
	return nil
}

func (s *storage) SetIntrospectionFromToken(context.Context, *oidc.IntrospectionResponse, string, string, string) error {
	return errUnsupported
}

func (s *storage) GetPrivateClaimsFromScopes(context.Context, string, string, []string) (map[string]any, error) {
	return nil, nil
}

func (s *storage) GetKeyByIDAndClientID(context.Context, string, string) (*jose.JSONWebKey, error) {
	return nil, errUnsupported
}

func (s *storage) ValidateJWTProfileScopes(context.Context, string, []string) ([]string, error) {
	return nil, errUnsupported
}
