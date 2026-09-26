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
// Auth is a JWT in an httpOnly cookie (ADR-0003/0008): login/signup/magic-link
// endpoints set the cookie server-side, so the contract just returns the
// `User`. Legacy URLs preserved. Phone sign-in and passkeys are out of the
// cutover (ADR-0015) and have no operations here.
// ---------------------------------------------------------------------------

// ---- Inputs ---------------------------------------------------------------

/** Sign-up form (legacy: `UserSignUpForm`). */
model SignUpInput {
  firstName: string | null;
  @format("email") email: string;
  @minLength(6) password: string;
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
  @minLength(6) password: string;
}

model EmailInput {
  @format("email") email: string;
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
  startedEnrollments: Enrollment[];
  finishedEnrollments: Enrollment[];

  /** Landing page (catalog item) keyed by course id, for each enrollment. */
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
  consume(@path token: string): SessionResponse | NotFoundError | ApiError;
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

  /** Set a new password using a reset token; signs the user in. */
  @patch
  @route("/{token}")
  @operationId("updatePassword")
  update(
    @path token: string,
    @body body: ResetPasswordInput,
  ): SessionResponse | NotFoundError | ValidationError | ApiError;
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
