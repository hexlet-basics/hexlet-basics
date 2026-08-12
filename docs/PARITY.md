# Parity gap: legacy Rails → Go rewrite

Snapshot of what the legacy Rails app (`legacy/`) still does that the Go stack
at the root does not. The plan of record is a hard cutover at parity
(ADR-0002), so everything below is cutover-blocking unless marked *dropped*.

Concepts are named in the `CONTEXT.md` vocabulary (Enrollment, Lesson Progress,
Completion); table, column and legacy job names are quoted as they are, because
storage keeps the old words.

Method: `operationId`s in `api-spec/dist/openapi.yaml` (107) diffed against
handler methods in `internal/`, plus `legacy/config/routes.rb` diffed against
`api-spec/*.tsp` and `src/routes/`, plus `legacy/db/schema.rb` tables diffed
against `ent/schema/`.

## Where it stands

| Surface | Contract | Go handlers | Frontend |
|---|---|---|---|
| Admin (all resources) | done | done (66 ops) | done (17 screens) |
| Public catalog `/languages` | done | `listCourses` only | catalog page only |
| Course landing | done | done, with progress | stub page |
| Lesson player | done (first cut) | done (read + check + runner) | none |
| Learner progress (start, check, dashboard, guest) | done | done | none |
| Blog, categories, reviews, pages, sitemap | done | none | none |
| Auth (login/signup) | done | done | done |
| Auth (magic link, password reset, phone, passkeys) | done | none | none |
| User area (`/my`) | done | done | dashboard only |
| Profile, locale, account deletion | done | none | none |
| Cases, book download, feeds, error pages | **absent** | none | none |

76 of 108 contract operations have handlers. The 32 without one are listed
per blocker below (`adminUploadAttachment` is the exception — it is
implemented outside the generated layer in `internal/handlers/attachments.go`
because ogen cannot generate multipart, so the real count is 31).

## Blockers, not tickets

The unimplemented operations cluster under a handful of missing foundations.
Each foundation unlocks its whole group; sequencing should follow this graph,
not the operation list.

### 1. Enrollment and progress — **done** (ADR-0012)

The learner spine is built: `internal/progress` owns sequential progression,
`ent/schema/enrollment.go` maps `language_members` behind a unique (user,
course) index, and the four progress events have publishers. `startLesson`,
`checkLesson`, `getCourse` and `getMyDashboard` are implemented, for a guest
carrying a signed cookie as well as for a signed-in learner, and completion is
re-evaluated when a version is promoted rather than by the unported
`finish_language_members_job`.

What is left of it is frontend: the course landing page is still a stub and the
dashboard has no screen.

### 2. Code runner — **done** (ADR-0013)

The check is synchronous, as legacy had it, and runs in Docker through
`internal/exerciserunner`: isolation tightened past legacy with every limit
configurable, grading pinned per course version, and the exit-code contract
(0 / 124 / other) carried over unchanged. `getCourseLesson`, the player
payload, is implemented too.

Two follow-ups are open rather than done: pre-pulling images on promotion
(#767, so the first check after a deploy does not pay for a pull) and the
canary that would let the two off-by-default limits be turned on (#768).

What is left here is the frontend: the whole
`src/routes/…/languages/$slug/lessons/$slug` player, which now has a complete
contract to build against.

### 3. No mailer (ADR-0006 unimplemented)

No Postbox/SES code exists. Blocks: `createMagicLink`, `consumeMagicLink`,
`createPasswordReminder`, `checkPasswordResetToken`, `updatePassword`. Legacy
equivalents are `UserMailer#magic_link` / `#reset_password`.

### 4. No SMS sender

Blocks `createPhoneAuth`, `confirmPhoneAuth`. Legacy: `SendSmsJob`.

### 5. No `user_credentials` schema

Blocks all four passkey operations (`newPasskey`, `createPasskey`,
`listPasskeys`, `deletePasskey`) and `newPasskeySession` /
`createPasskeySession`. go-webauthn is chosen (ADR-0003) but unwired.

### 6. No `language_category_items` schema

Category → landing pages is a many-to-many through that table in legacy.
Without it `getPublicCourseCategory` (`CourseCategoryView.landingPages`) and
`listPublicCourseCategories` cannot be served.

### 7. No `book_requests` schema

Blocks `createBookRequest`; the `bookRequested` event has no publisher.

### 8. In-lesson AI chat

`internal/assistant` is only the OpenAI completion adapter (used by the
admin lesson-review worker). `listAssistantMessages` / `createAssistantMessage`
are unimplemented, and the SSE event schema for streaming tokens is undesigned.
Legacy: `Ai::Lessons::MessagesController` + `Assistants::RunJob`.
(Commit 37ea6a9e implemented the *admin* messages list — a different operation.)

### 9. Static pages have no content source

`getPage` returns `bodyHtml`, but legacy `/pages/:slug` renders ERB views
(about, authors, privacy, tos, cookie). Where the HTML comes from in the
rewrite is an open decision.

### 10. Remaining public reads with no blocker

Straight ports, only unwritten: `listBlogPosts`, `getBlogPost`,
`getNextBlogPost`, `likeBlogPost`, `listPublicReviews`, `getSitemap`,
`getProfile`, `updateProfile`, `deleteAccount`, `switchLocale`, `createLead`
(the public lead form — everything downstream of it exists: the admin list, the
`leadCreated` consumer, the river job and the amoCRM client. Only the write
endpoint that publishes the event is missing).

## Legacy routes with no contract operation at all

Invisible to the operation diff — these are not in TypeSpec yet:

- `/cases` and `/cases/for_teachers`
- `GET /book/download` (only `POST /book/create_request` is modeled)
- `/api/feeds/yandex_courses` (Yandex course feed)
- `GET /languages/:slug/success` (course completion page)
- `/:code` error pages (403/404/500)
- Google one-tap / OAuth (`google_auth`, `auth#request|callback` — commented
  out in legacy routes; confirm whether it is live in production)
- `admin/flipper` (Flipper UI) and the SolidQueue monitor

## Async / jobs

`internal/jobs` has `exercise_loader`, `review_lesson`, `amocrm_lead`.
Unported legacy jobs:

- `finish_language_members_job` — **not being ported**: completion is
  re-evaluated when a version is promoted (ADR-0012), which is why the Go stack
  needs no nightly sweep for it
- `find_related_courses_for_blog_post_job` — admin sets related courses
  manually today; the automatic suggestion is missing
- `send_sms_job` (#4)
- `assistants/run_job` (#8)
- `reap_stuck_version_builds_job` — and there is **no recurring/cron mechanism
  at all** in the Go app; legacy schedules it hourly in `config/recurring.yml`

Also missing with no replacement decided: **feature flags** (Flipper) and
product analytics (`ahoy`).

## Tables with no ent schema

Classified so the raw count does not mislead:

- **Genuinely open subsystems, no contract, no plan:** `surveys`,
  `survey_items`, `survey_answers`, `survey_scenario*`, `user_survey_pivots`,
  `ahoy_events`, `ahoy_visits`, `taggings`, `tags`, `uploads`, `ai_models`,
  `ai_tool_calls`.
- **Blocking known work:** `user_credentials` (#5), `user_accounts` (OAuth),
  `language_category_items` (#6), `book_requests` (#7). (`language_members` is
  mapped now — see #1.)
- **Verify whether needed:** `language_version_infos`,
  `language_module_descriptions` (module/version info schemas exist for lessons
  and modules but not these).
- **Deliberately dropped:** `action_text_rich_texts` (blog stores plain HTML).
- **Infra replaced:** `event_store_events*` → watermill, `solid_queue_*` →
  river, `solid_cache_*` / `solid_cable_*` → n/a, `flipper_*` → *no
  replacement chosen*.

## Suggested order

1. ~~Enrollment (`language_members`) + progress events~~ — done.
2. ~~Build the exercise runner~~ — done. Next is the lesson player frontend on
   top of the contract that now exists, plus the course page: the backend of
   the product's core loop is complete and nothing renders it.
3. Mailer, then the four email-based auth flows.
4. The blocker-free public reads (#10) and their pages — blog, reviews,
   categories, sitemap, static pages.
5. In-lesson assistant (needs the SSE schema designed first).
6. Passkeys, phone auth, cases/book/feeds, error pages.
7. Recurring-task scheduler + feature flags.

## Note

`CLAUDE.md` still points at `docs/STACK.md`, deleted in commit 890ed208.
