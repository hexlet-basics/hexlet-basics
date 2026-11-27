---
status: accepted
---

# Typed state columns and one generated state writer

Amends the state-transition addenda of [ADR-0020](0020-contract-declared-permissions-and-audit-trail.md)
(blog posts; banners, reviews and landing pages; Course Readiness). The operations,
their routes, Permissions, the transition tables and the 409 they answer with are
unchanged. What changes is where a move is written and how a state is typed.

**Problem.** Five modules carried the same `move()`: read the row, `Check` the
transition, `UpdateOneID().SetState()` (`blog`, `sitebanner`, `coursereview`,
`landingpage`, and `catalog.Readiness` for the Course Readiness). Three of them
existed for that alone. The read side named "published" four more times, once per
package (`blog.StatePublished`, `landingpage.StatePublished`,
`coursereview.StatePublished`, a bare literal in `sitebanner`), none of them the
`transition` constant, and `landingpage.CourseCompleted` named "completed" again.
The columns were plain nullable strings, with `CurrentOf`/`ReadinessOf` reading a
NULL as a draft, and the contract had four identical enums (`BannerState`,
`BlogPostState`, `ReviewState`, `LandingPageState`), three of them nullable.
The check and the write were two statements, so two concurrent publishes both
passed the check and both left an Audit Record claiming draft to published, which
the addendum promises never happens.

**Decision.**

- **One vocabulary, typed end to end.** The state columns are ent enums whose Go
  type is the `transition` type (`field.Enum(...).GoType(transition.PublicationState(""))`,
  `transition.Readiness` for the Course), so a query or a write that names a
  state can only name one of the vocabulary's, and ent validates the value. The
  contract has one `PublicationState` enum (the GLOSSARY term) for the four
  resources and a required `CourseReadiness`; the frontend takes both from the
  generated client.
- **Required, a draft by default.** Migrations backfill the NULLs as `draft`, set
  the default, and make the columns NOT NULL (through a validated check, so no
  table is scanned under an exclusive lock). `CurrentOf`, `ReadinessOf` and the
  "a row with no state counts as a draft" tests are gone with the NULLs.
- **One writer, generated.** A state field is annotated `schema.StateMove`; the
  `statemove` template (`ent/template/statemove.tmpl`) generates `Move<Field>` on
  its client (`MoveState`, `MoveReadiness`), like `AdminCRUD` generates the admin
  helpers: the operation has a rule, but one rule shared by every resource, with
  one implementation. It is a single conditional `UpdateOne`
  (`WHERE state IN (move.From)`), so a racing move is refused; when it matches
  nothing it re-reads the row to answer not found or a `*transition.IllegalError`.
  The write is an `UpdateOne`, so the audit hook still records the old and the
  new state. The handlers call it through `moveState`; the five editors and their
  wiring are deleted.
- **Enforced.** `forbidigo` rejects the state setters (`SetState`,
  `SetReadiness` and their `SetNillable*`) on those builders outside the
  generated code; tests may arrange a starting state with them.
- **`domain` imports nothing of ours.** The generated code imports `transition`,
  which wraps `domain.ErrConflict`, and `domain` imported ent for its error
  translator: a cycle. The translator and the pager move to the client's package
  (`store.ToDomain`, `store.UniqueOn`, `store.Paged`); `domain` keeps the error
  vocabulary and `Page`.

**Consequences.** A new resource with a Publication State is a column typed by
`transition.PublicationState` with the annotation, three operations in the
contract, and three one-line handlers; it cannot write the state any other way.
A move that one day needs a side effect (an event, a job) moves into a use-case
module that calls `Move<Field>` and does the rest, as ADR-0020 already says for
side effects. The other state columns (Enrollment, Lesson Progress, versions,
books, users) are still plain strings with per-package constants: typing them the
same way needs their state types moved to a package ent may import, which is left
for when one of them is next changed.
