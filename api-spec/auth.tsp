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
// Auth is a server-side session behind an httpOnly cookie (ADR-0017):
// login/signup/magic-link endpoints start it, and the server writes its cookie
// itself, so the contract just returns the `User`. Legacy URLs preserved. Phone sign-in and passkeys are out of the
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

/**
 * A name on the profile (legacy `User` validations on first/last name): at
 * most 40 characters, none of the characters in the pattern; blank is allowed.
 */
@maxLength(40)
@pattern("^[^`!@#$%\\^&*+=]*$")
scalar ProfileName extends string;

/** Profile edit form (legacy: `UserProfileForm`). */
model ProfileInput {
  firstName: ProfileName | null;
  lastName: ProfileName | null;
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
 * A successful sign-in. The session cookie is not part of the contract: the
 * server's session middleware writes it. `Set-Cookie` here carries only the
 * cookie that ends the visitor's guest progress, once it has been credited to
 * the account that just signed in.
 */
model SessionResponse {
  @header("Set-Cookie") setCookie?: string[];
  @body user: User;
}

/** A successful sign-out; the session cookie is expired by the server. */
model SessionDeleted {
  @statusCode statusCode: 204;
}

/** A remembered locale choice and the cookie that carries it. */
model LocaleSwitched {
  @statusCode statusCode: 204;
  @header("Set-Cookie") setCookie: string[];
}

/** A newly-created, signed-in user; `Set-Cookie` as in SessionResponse. */
model UserCreated {
  @statusCode statusCode: 201;
  @header("Set-Cookie") setCookie?: string[];
  @body user: User;
}

/** One course on the dashboard: the enrollment, where the learner stands, and the card to draw. */
model DashboardCourse {
  enrollment: Enrollment;
  progress: CourseProgress;

  /** The course's main catalog entry in the request locale. */
  card: CourseCatalogItem;
}

/** The signed-in user's dashboard (legacy: `MyController#show`). */
model MyDashboard {
  /** Courses in progress, newest enrollment first. */
  started: DashboardCourse[];

  /** Finished courses, newest enrollment first. */
  finished: DashboardCourse[];
}

// ---- Session / current user ----------------------------------------------

@route("/session")
@tag("Auth")
interface Session {
  /** Log in with email + password; starts the session. */
  @post
  @operationId("createSession")
  create(@body body: SessionInput): SessionResponse | ValidationError | ApiError;

  /** Log out; ends the session. */
  @delete
  @useAuth(UserSession)
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
  /** Sign up; starts the session on success. */
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

  /** Consume a magic link token; starts the session. */
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
  @useAuth(UserSession)
  @operationId("updateProfile")
  update(@body body: ProfileInput): User | ValidationError | UnauthorizedError | ApiError;

  /** Delete the current user's account. */
  @delete
  @useAuth(UserSession)
  @operationId("deleteAccount")
  delete(): SessionDeleted | UnauthorizedError | ApiError;
}

// ---- Locale ---------------------------------------------------------------

@route("/locale")
@tag("Auth")
interface LocaleSwitch {
  /**
   * Remember the chosen UI locale: on the signed-in user, and in a cookie the
   * site root reads to pick the locale it redirects to. The page navigates
   * itself afterwards.
   */
  @get
  @route("/switch")
  @operationId("switchLocale")
  switch(@query locale: Locale): LocaleSwitched | ApiError;
}
