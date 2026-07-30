import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";

using Http;
using Rest;
using OpenAPI;

namespace HexletBasics;

// ---------------------------------------------------------------------------
// Authentication and the signed-in user area.
//
// Auth is a JWT in an httpOnly cookie (ADR-0003/0008): login/signup/magic-link/
// passkey endpoints set the cookie server-side, so the contract just returns
// the `User`; the Set-Cookie header is not modeled. Legacy URLs preserved.
// ---------------------------------------------------------------------------

// ---- Inputs ---------------------------------------------------------------

/** Sign-up form (legacy: `UserSignUpForm`). */
model SignUpInput {
  firstName: string | null;
  @format("email") email: string;
  @minLength(1) password: string;
}

/** Email + password login. */
model SessionInput {
  @format("email") email: string;
  @minLength(1) password: string;
}

/** Profile edit form (legacy: `UserProfileForm`). */
model ProfileInput {
  firstName: string | null;
  lastName: string | null;
}

/** Password reset submission (legacy: `UserPassword`). */
model ResetPasswordInput {
  @minLength(1) password: string;
}

model EmailInput {
  @format("email") email: string;
}

model PhoneInput {
  @minLength(1) phone: string;
}

model PhoneConfirmInput {
  @minLength(1) phone: string;
  @minLength(1) code: string;
}

/**
 * A WebAuthn ceremony payload. The challenge/options and the client response
 * are opaque JSON owned by `go-webauthn`; the contract carries them as strings
 * so the browser API round-trips them verbatim.
 */
model PasskeyChallenge {
  /** JSON-encoded PublicKeyCredentialCreationOptions / RequestOptions. */
  options: string;
}

model PasskeyRegistrationInput {
  /** JSON-encoded registration credential from `navigator.credentials.create`. */
  credential: string;
  nickname: string | null;
}

model PasskeyAssertionInput {
  /** JSON-encoded assertion credential from `navigator.credentials.get`. */
  credential: string;
}

/** The current user resolved from the session cookie (null when anonymous). */
model CurrentUser {
  user: User | null;
}

/**
 * A successful sign-in. Cookie headers are part of the HTTP contract so the
 * implementation stays behind the generated ogen seam.
 */
model SessionResponse {
  @header("Set-Cookie") setCookie: string[];
  @body user: User;
}

/** A successful sign-out with the expired auth cookies. */
model SessionDeleted {
  @statusCode statusCode: 204;
  @header("Set-Cookie") setCookie: string[];
}

/** A newly-created user and the auth cookies issued for that user. */
model UserCreated {
  @statusCode statusCode: 201;
  @header("Set-Cookie") setCookie: string[];
  @body user: User;
}

/** The signed-in user's dashboard (legacy: `MyController#show`). */
model MyDashboard {
  startedCourseMembers: CourseMember[];
  finishedCourseMembers: CourseMember[];

  /** Landing page (catalog item) keyed by course id, for each membership. */
  landingPagesByCourseId: Record<CourseCatalogItem>;
}

// ---- Session / current user ----------------------------------------------

@route("/session")
@tag("Auth")
interface Session {
  /** Log in with email + password; sets the JWT cookie. */
  @post
  @operationId("createSession")
  create(@body body: SessionInput): SessionResponse | ValidationError | ApiError;

  /** Log out; clears the JWT cookie. */
  @delete
  @useAuth([UserSession, XsrfToken])
  @operationId("deleteSession")
  delete(): SessionDeleted | UnauthorizedError | ApiError;
}

@route("/me")
@tag("Auth")
interface Me {
  /** Resolve the current user from the session cookie (for SSR). */
  @get
  @operationId("getCurrentUser")
  get(@header("Cookie") cookie?: string): CurrentUser | ApiError;
}

@route("/my")
@tag("User Area")
@useAuth(UserSession)
interface My {
  /** The signed-in user's course dashboard. */
  @get
  @operationId("getMyDashboard")
  get(): MyDashboard | UnauthorizedError | ApiError;
}

// ---- Registration & recovery ---------------------------------------------

@route("/users")
@tag("Auth")
interface Registration {
  /** Sign up; sets the JWT cookie on success. */
  @post
  @operationId("createUser")
  create(@body body: SignUpInput): UserCreated | ValidationError | ApiError;
}

@route("/magic_links")
@tag("Auth")
interface MagicLinks {
  /** Email a one-time sign-in link. */
  @post
  @operationId("createMagicLink")
  create(@body body: EmailInput): NoContent | ValidationError | ApiError;

  /** Consume a magic link token; sets the JWT cookie. */
  @get
  @route("/{token}")
  @operationId("consumeMagicLink")
  consume(@path token: string): User | NotFoundError | ApiError;
}

@route("/remind_password")
@tag("Auth")
interface RemindPassword {
  /** Email a password-reset link. */
  @post
  @operationId("createPasswordReminder")
  create(@body body: EmailInput): NoContent | ValidationError | ApiError;
}

@route("/password")
@tag("Auth")
interface PasswordReset {
  /** Validate a reset token before showing the form. */
  @get
  @route("/{token}/edit")
  @operationId("checkPasswordResetToken")
  check(@path token: string): NoContent | NotFoundError | ApiError;

  /** Set a new password using a reset token. */
  @patch
  @route("/{token}")
  @operationId("updatePassword")
  update(@path token: string, @body body: ResetPasswordInput): NoContent | ValidationError | ApiError;
}

@route("/phone_auth")
@tag("Auth")
interface PhoneAuth {
  /** Send an SMS verification code. */
  @post
  @operationId("createPhoneAuth")
  create(@body body: PhoneInput): NoContent | ValidationError | ApiError;

  /** Confirm the SMS code; sets the JWT cookie. */
  @post
  @route("/confirm")
  @operationId("confirmPhoneAuth")
  confirm(@body body: PhoneConfirmInput): User | ValidationError | ApiError;
}

// ---- Passkeys / WebAuthn --------------------------------------------------

@route("/passkey_session")
@tag("Auth")
interface PasskeySession {
  /** Begin a passkey login ceremony. */
  @get
  @route("/new")
  @operationId("newPasskeySession")
  new(): PasskeyChallenge | ApiError;

  /** Complete a passkey login; sets the JWT cookie. */
  @post
  @operationId("createPasskeySession")
  create(@body body: PasskeyAssertionInput): User | ValidationError | ApiError;
}

@route("/account/passkeys")
@tag("User Area")
@useAuth(UserSession)
interface AccountPasskeys {
  /** List the current user's passkeys. */
  @get
  @operationId("listPasskeys")
  list(): UserCredential[] | UnauthorizedError | ApiError;

  /** Begin a passkey registration ceremony. */
  @get
  @route("/new")
  @operationId("newPasskey")
  new(): PasskeyChallenge | UnauthorizedError | ApiError;

  /** Register a new passkey for the current user. */
  @post
  @useAuth([UserSession, XsrfToken])
  @operationId("createPasskey")
  create(@body body: PasskeyRegistrationInput): {
    @statusCode statusCode: 201;
    @body credential: UserCredential;
  } | ValidationError | UnauthorizedError | ApiError;

  /** Remove a passkey. */
  @delete
  @useAuth([UserSession, XsrfToken])
  @route("/{id}")
  @operationId("deletePasskey")
  delete(@path id: int32): NoContent | NotFoundError | UnauthorizedError | ApiError;
}

// ---- Account profile ------------------------------------------------------

@route("/account/profile")
@tag("User Area")
@useAuth(UserSession)
interface AccountProfile {
  /** The editable profile of the current user. */
  @get
  @route("/edit")
  @operationId("getProfile")
  get(): User | UnauthorizedError | ApiError;

  /** Update the current user's profile. */
  @patch
  @useAuth([UserSession, XsrfToken])
  @operationId("updateProfile")
  update(@body body: ProfileInput): User | ValidationError | UnauthorizedError | ApiError;

  /** Delete the current user's account. */
  @delete
  @useAuth([UserSession, XsrfToken])
  @operationId("deleteAccount")
  delete(): NoContent | UnauthorizedError | ApiError;
}

// ---- Locale ---------------------------------------------------------------

@route("/locale")
@tag("Auth")
interface LocaleSwitch {
  /** Persist the preferred UI locale on the session. */
  @get
  @route("/switch")
  @operationId("switchLocale")
  switch(@query locale: Locale): NoContent | ApiError;
}
