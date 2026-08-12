import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";

using Http;
using Rest;
using OpenAPI;

namespace HexletBasics;

// ===========================================================================
// LESSON PLAYER — the hard, NON-CRUD surface (FRONTEND_PLAN Wave 4).
//
// This is a first cut so the player can be built against a mock; it is NOT a
// finished design. The two pieces that need their own pass:
//
//   1. Code run / solution check. Legacy does this synchronously
//      (`POST /api/lessons/:id/check` blocks on the Docker run). In the Go
//      rewrite the run is a river job shelling to Docker (STACK.md), so this
//      should likely become submit -> poll/stream, not one blocking POST.
//      Modeled synchronously below to match legacy; revisit before Wave 4.
//
//   2. AI assistant streaming. Legacy enqueues a job and streams tokens back.
//      hey-api can emit an SSE client (`serverSentEvents.gen.ts`), but the
//      event schema (token deltas, done, error) is not designed here yet.
// ===========================================================================

/** Result of running a submitted solution (legacy: `LessonCheckingResponse`). */
model LessonCheckingResponse {
  passed: boolean;
  output: string;
  result: "passed" | "failed" | "failed-infinity";
  status: int32;
  lessonHasBeenFinished: boolean;
  courseHasBeenFinished: boolean;
}

/** A submitted solution to run against the lesson's tests. */
model CheckLessonInput {
  @minLength(1) code: string;

  /** The lesson version the code was written against. */
  versionId: int32;
}

@route("/lessons/{id}/check")
@tag("Lesson Player")
interface LessonCheck {
  /**
   * Run a submitted solution and record progress. Synchronous to match legacy;
   * revisit as submit + stream for the Docker-job model (see file header).
   */
  @post
  @operationId("checkLesson")
  check(@path id: int32, @body body: CheckLessonInput): LessonCheckingResponse | ValidationError | ApiError;
}

// ---- Starting a lesson -----------------------------------------------------

/**
 * The learner has not reached this lesson yet. Sequential progression
 * (ADR-0012) makes a lesson available only when its position in the course's
 * current version is at most one past the furthest lesson they have finished.
 *
 * Distinct from 403: the caller's permissions are fine, their progress is not.
 */
@error
model LessonNotAvailableError {
  @statusCode statusCode: 409;
  @header contentType: "application/problem+json";
  @body body: ProblemDetails;
}

@route("/lessons/{id}/start")
@tag("Lesson Player")
@useAuth(UserSession)
interface LessonStart {
  /**
   * Start a lesson: enroll the learner in its course if they are not enrolled
   * yet, and mark the lesson started. This is the only way progress begins —
   * loading a page never has this effect, because the frontend preloads routes
   * on hover.
   *
   * Idempotent: starting an already-started or already-finished lesson
   * succeeds and changes nothing. 409 when the lesson is beyond the gate.
   */
  @post
  @useAuth([UserSession, XsrfToken])
  @operationId("startLesson")
  start(@path id: int32): NoContent | LessonNotAvailableError | NotFoundError | UnauthorizedError | ApiError;
}

// ---- AI assistant (in-lesson chat) ----------------------------------------

/** One message in the in-lesson AI chat (legacy: `AiMessage`). */
model LessonAssistantMessage {
  id: int32;
  role: string;
  userId: int32 | null;
  content: string;
  courseSlug: string;
  courseLessonSlug: string;
  courseLessonName: string;
  createdAt: utcDateTime;
}

/** A prompt to the in-lesson assistant, carrying the editor context. */
model AssistantMessageInput {
  @minLength(1) message: string;
  output: string | null;
  userCode: string | null;
}

@route("/ai/lessons/{lessonId}/messages")
@tag("Lesson Player")
@useAuth(UserSession)
interface LessonAssistant {
  /** The assistant chat history for the current user in this lesson. */
  @get
  @operationId("listAssistantMessages")
  list(@path lessonId: int32): LessonAssistantMessage[] | UnauthorizedError | ApiError;

  /**
   * Ask the assistant. Enqueues generation (river job); the reply is delivered
   * out-of-band. 202 = accepted, 429 = rate limited. SSE streaming is TBD.
   */
  @post
  @useAuth([UserSession, XsrfToken])
  @operationId("createAssistantMessage")
  create(@path lessonId: int32, @body body: AssistantMessageInput): {
    @statusCode statusCode: 202;
  } | UnauthorizedError | ApiError;
}
