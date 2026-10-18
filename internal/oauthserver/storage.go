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
	db     *ent.Client
	tx     store.Transactor
	issuer string
	// resource is the canonical MCP URL every token is issued for.
	resource string
	now      func() time.Time
	signing  ed25519.PrivateKey
	// accessLifetime and refreshLifetime are how long the two token kinds live.
	accessLifetime  time.Duration
	refreshLifetime time.Duration
}

var _ op.Storage = (*storage)(nil)

// errUnsupported marks the OpenID features this server does not offer.
var errUnsupported = errors.New("oauthserver: not supported")

// randomToken returns n random bytes as unpadded base64url: an unguessable
// credential or id.
func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", errs.Wrapf(err, "read random bytes")
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// validTokenID reports whether id could be one randomToken minted: the base64url
// alphabet and a sane length. The library hands the storage ids it recovered by
// decrypting a string a client presented, and its decryption falls back to AES
// in CFB mode, which does not authenticate: a string that was never ours
// "decrypts" to arbitrary bytes, and when those happen to hold one colon they
// are split into an id and a subject. Such bytes are not valid UTF-8, so looking
// them up makes PostgreSQL fail (22021) and the request would answer 500, where
// revoking an unknown token is 200 (RFC 7009) and presenting one is 401.
func validTokenID(id string) bool {
	if id == "" || len(id) > 128 {
		return false
	}
	for i := range len(id) {
		c := id[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '-', c == '_':
		default:
			return false
		}
	}
	return true
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
// S256 is mandatory for every client (ADR-0022).
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
		return nil, errs.Wrapf(err, "create authorization request")
	}
	return &authRequest{row: row, resource: s.resource}, nil
}

// AuthRequestByID loads a request that has not expired.
func (s *storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	row, err := s.db.OAuthAuthRequest.Query().
		Where(oauthauthrequest.ID(id), oauthauthrequest.ExpiresAtGT(s.now())).
		Only(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "authorization request")
	}
	return &authRequest{row: row, resource: s.resource}, nil
}

// AuthRequestByCode finds the request an authorization code was issued for,
// by the code's hash.
func (s *storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	hash := digest(code)
	row, err := s.db.OAuthAuthRequest.Query().
		Where(oauthauthrequest.CodeHash(hash), oauthauthrequest.ExpiresAtGT(s.now())).
		Only(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "authorization code")
	}
	return &authRequest{row: row, resource: s.resource}, nil
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
		return errs.Wrapf(err, "save authorization code")
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
		return errs.Wrapf(err, "delete authorization request")
	}
	return nil
}

// authRequest adapts a stored request to op.AuthRequest.
type authRequest struct {
	row      *ent.OAuthAuthRequest
	resource string
}

var _ op.AuthRequest = (*authRequest)(nil)

func (a *authRequest) GetID() string       { return a.row.ID }
func (a *authRequest) GetACR() string      { return "" }
func (a *authRequest) GetAMR() []string    { return nil }
func (a *authRequest) GetClientID() string { return a.row.ClientID }

// GetAudience is the canonical MCP URL: the authorization request was admitted
// only for that resource (RFC 8707), so the token is bound to it.
func (a *authRequest) GetAudience() []string { return []string{a.resource} }
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
		return 0, errs.With("subject", req.GetSubject()).Wrapf(err, "token subject is not a user")
	}
	return id, nil
}

// CreateAccessToken records an access token for the code exchange. The value
// the client receives is built by the library from the returned id.
func (s *storage) CreateAccessToken(ctx context.Context, req op.TokenRequest) (string, time.Time, error) {
	owner, ok := req.(tokenOwner)
	if !ok {
		return "", time.Time{}, errs.With("request_type", fmt.Sprintf("%T", req)).Errorf("unsupported token request")
	}
	userID, err := userIDOf(req)
	if err != nil {
		return "", time.Time{}, err
	}
	// Without a refresh token there is nothing to rotate: the token is a family
	// of its own, which only an explicit revocation can end.
	family, err := randomToken(16)
	if err != nil {
		return "", time.Time{}, err
	}
	id, expires, err := s.insertAccessToken(ctx, s.db, owner, userID, "", family)
	return id, expires, err
}

// CreateAccessAndRefreshTokens issues a pair: for a first grant (no current
// refresh token, which starts a new family) or in exchange for a refresh token,
// which is marked used and replaced by a child in the same family.
//
// Spending a token is one conditional UPDATE (unused and unexpired), so of two
// parallel exchanges of the same token the database lets exactly one claim it.
// The other finds it used and is treated as a replay: the family, including the
// child the winner just issued, is revoked. That is deliberately strict, with no
// grace window (ADR-0022 addendum, #839).
func (s *storage) CreateAccessAndRefreshTokens(
	ctx context.Context, req op.TokenRequest, current string,
) (string, string, time.Time, error) {
	owner, ok := req.(tokenOwner)
	if !ok {
		return "", "", time.Time{}, errs.With("request_type", fmt.Sprintf("%T", req)).Errorf("unsupported token request")
	}
	userID, err := userIDOf(req)
	if err != nil {
		return "", "", time.Time{}, err
	}

	var accessID, refresh string
	var expires time.Time
	err = s.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		var family string
		if current == "" {
			family, err = randomToken(16)
			if err != nil {
				return err
			}
		} else {
			// Spend the presented token, then end the access token it was paired
			// with: only the new pair is alive afterwards.
			hash := digest(current)
			claimed, err := db.OAuthRefreshToken.Update().
				Where(oauthrefreshtoken.TokenHash(hash), oauthrefreshtoken.UsedAtIsNil(), oauthrefreshtoken.ExpiresAtGT(s.now())).
				SetUsedAt(s.now()).
				Save(ctx)
			if err != nil {
				return errs.Wrapf(err, "spend refresh token")
			}
			if claimed == 0 {
				return errRefreshNotSpendable
			}
			old, err := db.OAuthRefreshToken.Query().Where(oauthrefreshtoken.TokenHash(hash)).Only(ctx)
			if err != nil {
				return errs.Wrapf(err, "refresh token")
			}
			family = old.FamilyID
			if _, err := db.OAuthAccessToken.Delete().
				Where(oauthaccesstoken.RefreshTokenID(old.ID)).Exec(ctx); err != nil {
				return errs.Wrapf(err, "delete paired access token")
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
		accessID, expires, err = s.insertAccessToken(ctx, db, owner, userID, refreshID, family)
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
			SetFamilyID(family).
			SetExpiresAt(s.now().Add(s.refreshLifetime)).
			Exec(ctx)
	})
	if errors.Is(err, errRefreshNotSpendable) {
		// The spend was rolled back with the rest of the transaction, so a replay
		// is dealt with in a transaction of its own.
		return "", "", time.Time{}, s.refuseUnspendable(ctx, current)
	}
	if err != nil {
		return "", "", time.Time{}, err
	}
	return accessID, refresh, expires, nil
}

// TokenRequestByRefreshToken finds the grant a presented refresh token belongs to.
// A token that was already used is not a grant: it is a replay, and the family
// is revoked before the request is refused.
func (s *storage) TokenRequestByRefreshToken(ctx context.Context, refresh string) (op.RefreshTokenRequest, error) {
	row, err := s.db.OAuthRefreshToken.Query().
		Where(oauthrefreshtoken.TokenHash(digest(refresh)), oauthrefreshtoken.ExpiresAtGT(s.now())).
		Only(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "refresh token")
	}
	if row.UsedAt != nil {
		return nil, s.refuseReuse(ctx, row)
	}
	return &refreshRequest{row: row, scopes: row.Scopes, resource: s.resource}, nil
}

// errRefreshNotSpendable marks a refresh token that was not unused and live at
// the moment of the exchange.
var errRefreshNotSpendable = errors.New("refresh token cannot be spent")

// GetRefreshTokenInfo lets the revocation endpoint tell a refresh token from an
// access token.
func (s *storage) GetRefreshTokenInfo(ctx context.Context, _ string, token string) (string, string, error) {
	row, err := s.db.OAuthRefreshToken.Query().
		Where(oauthrefreshtoken.TokenHash(digest(token))).Only(ctx)
	if ent.IsNotFound(err) {
		return "", "", op.ErrInvalidRefreshToken
	}
	if err != nil {
		return "", "", errs.Wrapf(err, "refresh token info")
	}
	return strconv.Itoa(row.UserID), row.ID, nil
}

// RevokeToken ends the token with this id. An access token is deleted alone, so
// the client can still refresh; a refresh token ends its whole family, since
// the user withdrawing a connection means every token of it, the rotated ones
// included. A token that is already gone is not an error: the caller wanted it
// dead and it is. Only the client the token was issued to may revoke it.
func (s *storage) RevokeToken(ctx context.Context, id, _, clientID string) *oidc.Error {
	if !validTokenID(id) {
		// Nothing of ours has this id (see validTokenID), so there is nothing to
		// end, and RFC 7009 answers 200 for a token that is not there.
		return nil
	}
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
		return deleteFamily(ctx, db, refresh.FamilyID)
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

// TerminateSession ends every grant a user made to one client.
func (s *storage) TerminateSession(ctx context.Context, userID, clientID string) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return errs.With("subject", userID).Wrapf(err, "session subject is not a user")
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

func deleteFamily(ctx context.Context, db *ent.Client, familyID string) error {
	if _, err := db.OAuthAccessToken.Delete().Where(oauthaccesstoken.FamilyID(familyID)).Exec(ctx); err != nil {
		return errs.Wrapf(err, "delete family access tokens")
	}
	if _, err := db.OAuthRefreshToken.Delete().Where(oauthrefreshtoken.FamilyID(familyID)).Exec(ctx); err != nil {
		return errs.Wrapf(err, "delete family refresh tokens")
	}
	return nil
}

func (s *storage) SigningKey(context.Context) (op.SigningKey, error) {
	return signingKey{key: s.signing}, nil
}

// refreshRequest adapts a stored refresh token to op.RefreshTokenRequest.
type refreshRequest struct {
	row      *ent.OAuthRefreshToken
	scopes   []string
	resource string
}

var _ op.RefreshTokenRequest = (*refreshRequest)(nil)

func (r *refreshRequest) GetAMR() []string            { return nil }
func (r *refreshRequest) GetAudience() []string       { return []string{r.resource} }
func (r *refreshRequest) GetAuthTime() time.Time      { return r.row.AuthTime }
func (r *refreshRequest) GetClientID() string         { return r.row.ClientID }
func (r *refreshRequest) GetScopes() []string         { return r.scopes }
func (r *refreshRequest) GetSubject() string          { return strconv.Itoa(r.row.UserID) }
func (r *refreshRequest) SetCurrentScopes(s []string) { r.scopes = s }

func (s *storage) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	return []jose.SignatureAlgorithm{jose.EdDSA}, nil
}

func (s *storage) KeySet(context.Context) ([]op.Key, error) {
	return []op.Key{publicKey{key: s.signing.Public().(ed25519.PublicKey)}}, nil
}

var errWrongClient = errors.New("token belongs to another client")

// GetClientByClientID loads a client for the current request.
func (s *storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	row, err := s.db.OAuthClient.Query().Where(oauthclient.ClientID(clientID)).Only(ctx)
	if err != nil {
		return nil, errs.With("client_id", clientID).Wrapf(err, "load client")
	}
	return &client{row: row, loginBase: s.issuer + consentPagePath}, nil
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

// AuthorizeClientIDSecret always refuses: every client is public.
func (s *storage) AuthorizeClientIDSecret(context.Context, string, string) error {
	return errors.New("clients have no secret: authenticate with PKCE")
}

func (s *storage) Health(context.Context) error { return nil }

// SetUserinfoFromScopes fills the claims of an ID token. The library copies the
// userinfo it is given over the token's own claims, so leaving Subject empty
// here would blank the `sub` of every ID token. Nothing else about the user is
// disclosed: no profile or email claim is offered.
func (s *storage) SetUserinfoFromScopes(_ context.Context, info *oidc.UserInfo, userID, _ string, _ []string) error {
	info.Subject = userID
	return nil
}

// --- clients ---------------------------------------------------------------

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

// --- OpenID-only parts, stubbed --------------------------------------------

func (s *storage) GetPrivateClaimsFromScopes(context.Context, string, string, []string) (map[string]any, error) {
	return nil, nil //nolint:nilnil // op.Storage's contract: no private claims is a nil map, not an error
}

func (s *storage) GetKeyByIDAndClientID(context.Context, string, string) (*jose.JSONWebKey, error) {
	return nil, errUnsupported
}

func (s *storage) ValidateJWTProfileScopes(context.Context, string, []string) ([]string, error) {
	return nil, errUnsupported
}

func (s *storage) insertAccessToken(
	ctx context.Context, db *ent.Client, owner tokenOwner, userID int, refreshID, familyID string,
) (string, time.Time, error) {
	id, err := randomToken(16)
	if err != nil {
		return "", time.Time{}, err
	}
	expires := s.now().Add(s.accessLifetime)
	create := db.OAuthAccessToken.Create().
		SetID(id).
		SetClientID(owner.GetClientID()).
		SetUserID(userID).
		SetScopes(pq.StringArray(owner.GetScopes())).
		SetFamilyID(familyID).
		SetExpiresAt(expires)
	if refreshID != "" {
		create.SetRefreshTokenID(refreshID)
	}
	if err := create.Exec(ctx); err != nil {
		return "", time.Time{}, errs.Wrapf(err, "create access token")
	}
	return id, expires, nil
}

// refuseUnspendable answers an exchange that lost the race for a token: a token
// that is on file and used is a replay (or the loser of a parallel exchange), and
// revokes its family; anything else is simply unknown or expired.
func (s *storage) refuseUnspendable(ctx context.Context, refresh string) error {
	row, err := s.db.OAuthRefreshToken.Query().
		Where(oauthrefreshtoken.TokenHash(digest(refresh)), oauthrefreshtoken.ExpiresAtGT(s.now())).
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return oidc.ErrInvalidGrant().WithDescription("refresh token is invalid or expired")
	case err != nil:
		return errs.Wrapf(err, "refresh token")
	}
	return s.refuseReuse(ctx, row)
}

// refuseReuse revokes the family of a refresh token presented a second time and
// returns the refusal. A used token is evidence it has been copied: whoever holds
// the newer tokens, the thief and the owner alike, loses them, and the user must
// sign in again (ADR-0022, the MCP authorization spec).
func (s *storage) refuseReuse(ctx context.Context, row *ent.OAuthRefreshToken) error {
	if err := s.revokeFamily(ctx, row.FamilyID); err != nil {
		return oidc.ErrServerError().WithParent(err)
	}
	return oidc.ErrInvalidGrant().WithDescription("refresh token was already used: the authorization is revoked, sign in again")
}

// revokeFamily deletes every access and refresh token descended from one
// authorization, used refresh tokens included.
func (s *storage) revokeFamily(ctx context.Context, familyID string) error {
	return s.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		return deleteFamily(ctx, db, familyID)
	})
}

// PurgeExpired deletes what can no longer be used: expired access and refresh
// tokens (the used refresh tokens kept for replay detection among them, once
// their own lifetime is over) and authorization requests nobody completed. It
// is what keeps the tables from growing without bound, since a rotated refresh
// token is retained instead of deleted. It returns how many rows it removed.
func PurgeExpired(ctx context.Context, db *ent.Client, now time.Time) (int, error) {
	total := 0
	n, err := db.OAuthAccessToken.Delete().Where(oauthaccesstoken.ExpiresAtLTE(now)).Exec(ctx)
	if err != nil {
		return total, errs.Wrapf(err, "purge access tokens")
	}
	total += n
	n, err = db.OAuthRefreshToken.Delete().Where(oauthrefreshtoken.ExpiresAtLTE(now)).Exec(ctx)
	if err != nil {
		return total, errs.Wrapf(err, "purge refresh tokens")
	}
	total += n
	n, err = db.OAuthAuthRequest.Delete().Where(oauthauthrequest.ExpiresAtLTE(now)).Exec(ctx)
	if err != nil {
		return total, errs.Wrapf(err, "purge authorization requests")
	}
	return total + n, nil
}
