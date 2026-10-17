---
status: accepted
---

# Lessons in any order

ADR-0012 made progression sequential: a Lesson was available only up to one past
the highest finished position. We drop that gate. Any Lesson may be started,
asked about and submitted in any order, which is what legacy did. Everything else
ADR-0012 decided stands: enrollment on a deliberate action, signed-cookie guest
progress, Version-relative positions and Completion.

The gate cost more than it earned. It put a lock on every Lesson ahead of the
learner while the page it led to opened in full and the editor stayed usable, so
the interface said "unavailable" and then let the learner in. It was enforced on
three server paths (start, check and the assistant) and shipped in every read
only to render that lock. Product-wise nobody asked for it: the legacy
application, which production runs, has no such rule, so removing it restores
parity rather than adding a feature.

**What goes away**

- The gate: `ErrLessonNotAvailable`, the 409 `LessonNotAvailableError` on
  `startLesson`, `checkLesson` and `createAssistantMessage`, and
  `lessons[].available` and `furthestFinishedPosition` in the read model. A field
  that is always true is not kept for symmetry.
- The lock mark and its string in the player and the Course page.

**What stays**

- Theory is public and indexed, as before. Starting is still an explicit command
  and a read still never enrolls (ADR-0012).
- **Next Lesson** is still the first unfinished Lesson in course order. It was
  always defined that way and never was the gate, so a learner who took the third
  Lesson first is sent back to the first gap, which is where "continue" should go.
- A Course is finished when every Current Lesson is finished, in whatever order.

**Guest progress becomes a set**

A guest's state was one entry per Course, the slug of the furthest Lesson they
finished, and was exact only because the gate made gaps unrepresentable. Without
the gate a visitor who passes the twentieth Lesson first has not passed the
nineteen before it, so the entry now names the **set of finished Lesson slugs**.
Everything about the cookie that was not about the prefix is unchanged: signed
with the authentication secret, httpOnly, a year long, slugs resolved against the
current Version at read time, the least recently touched Courses evicted first.

The merge into an account on sign-up and sign-in is now a **union**. It credits
exactly the Lessons the cookie names that the account lacks and never takes
anything away, where it used to credit the whole prefix and take the higher of
two positions. The facts it publishes are the same pair per credited Lesson, in
course order.

No cookie in the old format exists in the wild: the Go stack has not been cut
over, and legacy's guest progress lived in the Rails session. An old-format
cookie, were one presented, fails to decode into the new shape and is treated as
no cookie.

**Consequences**

The server no longer refuses a Lesson for progress reasons, so the client has
nothing to branch on and nothing to explain. A visitor can open and submit any
Lesson, as can a learner. What a learner has finished is a set of rows and a
guest's is a set of slugs, the same shape either way, which is simpler than the
prefix rule it replaces.
