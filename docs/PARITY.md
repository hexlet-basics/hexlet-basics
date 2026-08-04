# Parity gap: legacy Rails → Go rewrite

Snapshot of what the legacy Rails app (`legacy/`) still does that the Go stack
at the root does not. The plan of record is a hard cutover at parity
(ADR-0002), so everything below is cutover-blocking unless marked *dropped*.

Method: `operationId`s in `api-spec/dist/openapi.yaml` (107) diffed against
handler methods in `internal/`, plus `legacy/config/routes.rb` diffed against
`api-spec/*.tsp` and `src/routes/`, plus `legacy/db/schema.rb` tables diffed
against `ent/schema/`.

## Where it stands

| Surface | Contract | Go handlers | Frontend |
|---|---|---|---|
| Admin (all resources) | done | done (66 ops) | done (17 screens) |
| Public catalog `/languages` | done | `listCourses` only | catalog page only |
| Course landing / lesson player | done (first cut) | none | stub page |
| Blog, categories, reviews, pages, sitemap | done | none | none |
| Auth (login/signup) | done | done | done |
| Auth (magic link, password reset, phone, passkeys) | done | none | none |
| User area (`/my`, profile) | done | none | none |
| Cases, book download, feeds, error pages | **absent** | none | none |

73 of 107 contract operations have handlers. The 34 without one are listed
per blocker below (`adminUploadAttachment` is the exception — it is
implemented outside the generated layer in `internal/handlers/attachments.go`
because ogen cannot generate multipart, so the real count is 33).

## Blockers, not tickets

The unimplemented operations cluster under a handful of missing foundations.
Each foundation unlocks its whole group; sequencing should follow this graph,
not the operation list.

### 1. No course membership (`language_members`) — biggest one

`CourseMember` exists **only** in the contract: no `ent/schema`, no
`apiconv` mapping, no table access. This sits under the entire learner spine:

- `getMyDashboard` (`/my`)
- `CourseView.member` — the signed-in state of the course landing page
- lesson progress, course start/finish
- the `courseStarted` / `courseFinished` / `lessonStarted` / `lessonFinished`
  events, which are *defined* in `internal/events` but have no publisher
- `finish_language_members_job` (unported)

### 2. No code runner → no lesson player

`checkLesson` is unimplemented and nothing in the repo shells out to Docker
(the only Docker references are `courseloader` parsing and testcontainers).
Legacy runs the check synchronously inside the request; the rewrite plans a
river job (ADR-0004), which means the contract shape is **still undecided** —
`api-spec/lesson-player.tsp` says so in its own header. Decide sync vs
submit+poll/stream before building.

Also blocked here: `getCourseLesson` (the player payload) and the whole
`src/routes/…/languages/$slug/lessons/$slug` frontend.

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

Straight ports, only unwritten: `getCourse`, `listBlogPosts`, `getBlogPost`,
`getNextBlogPost`, `likeBlogPost`, `listPublicReviews`, `getSitemap`,
`getProfile`, `updateProfile`, `deleteAccount`, `switchLocale`.

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

- `finish_language_members_job` (blocked by #1)
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
- **Blocking known work:** `language_members` (#1), `user_credentials` (#5),
  `user_accounts` (OAuth), `language_category_items` (#6), `book_requests` (#7).
- **Verify whether needed:** `language_version_infos`,
  `language_module_descriptions` (module/version info schemas exist for lessons
  and modules but not these).
- **Deliberately dropped:** `action_text_rich_texts` (blog stores plain HTML).
- **Infra replaced:** `event_store_events*` → watermill, `solid_queue_*` →
  river, `solid_cache_*` / `solid_cable_*` → n/a, `flipper_*` → *no
  replacement chosen*.

## Suggested order

1. `language_members` + progress events (unlocks `/my`, course landing
   signed-in state, lesson progress).
2. Decide the lesson-check contract, then build the runner (the product's core
   loop and the only surface with an unsettled API design).
3. Mailer, then the four email-based auth flows.
4. The blocker-free public reads (#10) and their pages — blog, reviews,
   categories, sitemap, static pages.
5. In-lesson assistant (needs the SSE schema designed first).
6. Passkeys, phone auth, cases/book/feeds, error pages.
7. Recurring-task scheduler + feature flags.

## Note

`CLAUDE.md` still points at `docs/STACK.md`, deleted in commit 890ed208.
