# ADR-0012: Sequential progression, and guest progress in a signed cookie

**Status:** Accepted

## Context

The legacy application let a learner take a Course's Lessons in any order, and
created the Enrollment as a side effect of *viewing* a Lesson page. A visitor
without an account progressed through the same Lessons, with their passed
Lessons kept in the Rails server session and claimed on sign-up.

Neither behaviour survives the move to the Go stack, for reasons that are
structural rather than matters of taste.

**Enrollment on view is unsafe in the new frontend.** TanStack Router preloads
routes on hover, so a route loader runs for Lessons the learner never opened.
If viewing enrolled, pointing at a Course's Lesson list would create Lesson
Progress rows across it and publish a `lessonStarted` event for each —
inflating the very funnel those events feed.

**There is no server session to put guest progress in.** Authentication is a
stateless JWT in an httpOnly cookie (ADR-0003); the Go stack keeps no
server-side session store, and introducing one for this alone would add a
stateful dependency to an otherwise stateless HTTP tier.

Both replacements are expensive to reverse once learners are living under them:
the first changes what the marketing funnel measures, the second writes rows
into real accounts from state a client held.

## Decision

### Sequential progression

A Lesson is **available** when its position within the Course's current Version
is at most one greater than the highest position the learner has finished.

Availability is computed from the highest finished position, not from the first
unfinished one. A learner who has finished Lesson 20 keeps moving forward even
if Lesson 3 is unfinished. Gaps come from both directions — legacy data written
when order was free, and new Versions inserting Lessons — and neither should
push a learner backwards. A gap does not block progress, but a Course is
finished only when every one of its Current Lessons is finished, so a gap keeps
the Course open until the learner closes it.

**Position is only defined relative to a Version.** A Lesson Progress row
records a Lesson, while the ordering key lives on the Lesson Version, so
positions are resolved by joining the learner's finished Lessons to the current
Version's Lesson Versions. A finished Lesson absent from the current Version
contributes no position: it neither raises the gate nor counts toward
completion. This is what makes a shrinking Course behave sanely — dropping a
Lesson lowers the total without stranding the learner, and dropping the
learner's furthest Lesson moves the gate back to their furthest surviving one.

Progress advances only through deliberate actions — starting a Lesson, or
submitting a solution — never by loading a page. Starting is an explicit
contract operation and is idempotent. A check does not require an existing
Lesson Progress row: its precondition is the gate, and it creates the
Enrollment and the Lesson Progress lazily when they are missing, so arriving by
deep link and submitting still records progress.

The gate is enforced server-side on both the start command and the check. The
client renders from server-provided state and implements no rule of its own: the
per-Lesson read model carries availability alongside completion precisely so
that the rule exists once, on the server, with no client copy free to drift.

Theory stays fully public and indexed. On a Lesson the visitor cannot yet take,
the page renders in full and the editor stays usable; its actions lead to
sign-in, and an already-signed-in learner is sent to their dashboard.

### Guest progress in a signed cookie

A visitor without an account progresses under the same rule. Their entire state
is one entry per Course: the **slug of the furthest Lesson they finished**.
Sequential progression makes gaps unrepresentable for a guest, so that single
value is exact — everything before it is finished, everything after it is not.

The slug is stored rather than a bare position, because a position is
meaningless across Versions: promoting a Version that inserts or removes Lessons
renumbers everything, and a stored number would silently come to denote a
different Lesson. Storing the Lesson's identity and resolving its position
against the current Version at read time carries the guest correctly across
changes elsewhere in the Course. If the stored Lesson is no longer part of the
current Version there is nothing to resolve against — unlike the database case
there are no other finished Lessons to fall back to — so that Course's guest
position resets to the beginning.

The state is carried in a **signed, httpOnly cookie**, using the same secret as
the token service already signs authentication cookies with (ADR-0003). Signing
is not optional: this state is later converted into database rows and domain
events, so an unsigned cookie would let a visitor mint fabricated completions.
The cookie is not readable by the client — the server derives progress from it
and returns it in the response body, which is what keeps the client free of any
guest/learner branch. It is written by the server through a `Set-Cookie` on the
check response; the contract already models response cookies for the session
endpoints. Lifetime is one year, and if the cookie approaches its size limit
the least recently touched Courses are evicted.

The check therefore stays a public operation that writes progress. The contract
already declares it without a session, because guests must be able to submit,
and the cookie's signature is what stands in for the session it cannot have —
which is why the signature cannot be dropped.

Guests get no dashboard. The dashboard stays authenticated-only; it is where the
incentive to create an account lives.

Guest progress is merged into the account on **both sign-up and sign-in**,
inside the same transaction that issues the session, so the resulting domain
events leave through the existing transactional outbox together with the
sign-up or sign-in event. The merge takes the higher of the two positions per
Course, both resolved against the current Version, and clears the cookie. The
transferred set is the **prefix** `1..N`, not an explicit set of Lesson ids:
under sequential progression the prefix *is* the set. This is the one place
where rows are written for Lessons whose completion is implied by the rule
rather than individually observed, and it is a deliberate departure from the
legacy transfer, which moved an explicit — and possibly non-contiguous —
selection.

## Consequences

Two of these are hard to reverse, and both change numbers the marketing funnel
reports:

- **Enrollment now requires a deliberate action.** Enrollment and Lesson
  Progress counts, and the `courseStarted` / `lessonStarted` events, stop
  counting page views and start counting intent. Volumes drop against the
  legacy baseline without anything being broken.
- **Guests can no longer complete Lessons out of order.** Pre-signup activity is
  a prefix of the Course, so the shape of what a guest transfers on sign-up
  changes.

Further consequences:

- Handlers never evaluate positions, and nothing outside learns whether a
  learner's state came from the database or from a cookie. How that is packaged
  is an ordinary design choice and is deliberately not decided here.
- Both departures from the legacy behaviour are intentional. ADR-0002 keeps the
  cutover compatible on bcrypt passwords and public URLs; it does not promise
  behavioural parity, and these two decisions spend some.
- The public Course and Lesson reads become personalized: they accept the
  request cookie and return progress. Anonymous visitors with no cookie get the
  responses they get today, and CDN caching of those reads stays off — the
  legacy application never cached them either.
- Completion percentage is computed by counting finished Current Lessons, not
  from the denormalized counter on the Enrollment. That counter includes
  Lessons dropped by later Versions, which is why the legacy serializer had to
  clamp it at 100%. The public enrollment count on Courses stays denormalized
  (`enrollmentsCount`, reading the legacy `members_count` column); it is a
  marketing figure on cached catalogue pages, not a progress figure.
- Course completion is re-evaluated when a new Course Version is promoted,
  driven by the course-loading module. The legacy nightly sweep over started
  enrollments is not ported: it existed because completion was only ever
  evaluated during a check, and there is now an explicit event for the moment a
  Course's Lesson set changes.
- Sign-in and sign-up redirect an already-authenticated visitor to the
  dashboard. That is a general rule about the authentication pages, but the
  locked-Lesson call to action depends on it.
- A guest's forged or tampered cookie is rejected, and a cookie naming a Lesson
  absent from the current Version resets that Course's position. Neither can
  corrupt an account, because the merge takes the higher of the two positions.
