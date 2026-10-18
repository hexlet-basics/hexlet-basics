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
//   2. AI assistant streaming. Answered as a plain text stream from the
//      generated operation (see `AssistantAnswerStream`).
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

/**
 * A completed check: the exercise outcome, plus the progress cookie a guest's
 * advanced position travels back in.
 *
 * The cookie is absent for a signed-in learner, whose progress is rows rather
 * than a cookie — which is the only place the two paths differ, and the client
 * never has to notice: it reads the same body either way.
 */
model LessonCheckingResult {
  @header("Set-Cookie") setCookie?: string[];
  @body response: LessonCheckingResponse;
}

@route("/lessons/{id}/check")
@tag("Lesson Player")
interface LessonCheck {
  /**
   * Run a submitted solution and record progress. Synchronous to match legacy;
   * revisit as submit + stream for the Docker-job model (see file header).
   *
   * Public on purpose: a guest must be able to submit, and what they have
   * finished travels in the signed cookie rather than in rows. Any lesson may
   * be submitted, in any order.
   */
  @post
  @operationId("checkLesson")
  check(@path id: int32, @body body: CheckLessonInput):
    | LessonCheckingResult
    | ApiError;
}

// ---- Starting a lesson -----------------------------------------------------

@route("/lessons/{id}/start")
@tag("Lesson Player")
interface LessonStart {
  /**
   * Start a lesson: enroll the learner in its course if they are not enrolled
   * yet, and mark the lesson started. This is the only way progress begins —
   * loading a page never has this effect, because the frontend preloads routes
   * on hover.
   *
   * Idempotent: starting an already-started or already-finished lesson
   * succeeds and changes nothing. Any lesson may be started, in any order.
   *
   * Public, because a guest progresses the same way. For them it stores
   * nothing: their state is the signed cookie, and only a check moves it. It
   * answers with nothing: where the learner stands is the course and lesson
   * reads' business, and the page they go to next reads it from there.
   */
  @post
  @operationId("startLesson")
  start(@path id: int32): NoContent | ApiError;
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

/**
 * The learner's in-lesson chat, plus whether they may ask again today. The
 * quota flag rides along with the history so the panel renders its disabled
 * state on first paint instead of discovering it from a refused question.
 */
model LessonAssistantChat {
  /** User and assistant turns, oldest first; empty before the first question. */
  messages: LessonAssistantMessage[];

  /** The learner has used up today's questions (UTC day). */
  quotaExceeded: boolean;
}

/** A prompt to the in-lesson assistant, carrying the editor context. */
model AssistantMessageInput {
  @minLength(1) message: string;
  output: string | null;
  userCode: string | null;
}

/**
 * The assistant's answer, streamed as plain text while the model produces it.
 *
 * ogen generates a `text/plain` string body as an `io.Reader` the encoder copies
 * straight to the socket, which is what lets tokens reach the learner as they
 * arrive. (`bytes` would not do: TypeSpec emits it as `format: byte` for a text
 * content type, and ogen then base64-encodes the stream.) SSE is out of reach — ogen
 * cannot serve it yet (ogen#1742) — and the client does not need it: the AI SDK
 * reads a plain text stream directly.
 */
model AssistantAnswerStream {
  @header contentType: "text/plain";
  @body body: string;
}

/** The learner has used up today's assistant questions. */
@error
model AssistantQuotaExceededError {
  @statusCode statusCode: 429;
  @header contentType: "application/problem+json";
  @body body: ProblemDetails;
}

@route("/ai/lessons/{lessonId}/messages")
@tag("Lesson Player")
@useAuth(UserSession)
interface LessonAssistant {
  /**
   * The assistant chat history for the current user in this lesson, and their
   * remaining quota. A lesson they never asked about answers an empty chat.
   */
  @get
  @operationId("listAssistantMessages")
  list(@path lessonId: int32): LessonAssistantChat | ApiError;

  /**
   * Ask the assistant and stream its answer back. Both turns are stored once
   * the answer is complete, with the answer's token usage.
   *
   * Asking is doing the lesson, so it starts the lesson when the learner has
   * not yet. 429 once today's questions are used up — decided before anything
   * streams.
   */
  @post
  @useAuth(UserSession)
  @operationId("createAssistantMessage")
  create(@path lessonId: int32, @body body: AssistantMessageInput):
    | AssistantAnswerStream
    | AssistantQuotaExceededError
    | ApiError;
}
