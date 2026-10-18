---
status: accepted
---

# Singular package names, and interfaces declared only where a seam needs one

Part of the backend clean-up (spec #905, tickets #926 and #932). Replaces the
package-naming and interface rules of the "Use-case modules" section of
`AGENTS.md`, which made package names plural and had handlers declare a small
interface per module.

**Problem.** Plural names were chosen so that a package would not collide with
the singular package ent generates per entity (`ent/blogpost`). That avoided the
collision by accident of spelling, not by a rule, and plural names are not how Go
names packages. Meanwhile `handlers` declared about 25 interfaces, each with one
implementation and no test substitute: the tests run the real modules over a
rolled-back transaction (ADR-0019), so the interfaces only hid which type a
handler really called.

**Decision.**

- A package under `internal/` is named for its concept in the singular, and is
  never equal to an `ent/<x>` package, so no file needs an import alias.
  `mise run lint:packages` enforces it. Exceptions, recorded here: `events`,
  `jobs` and `handlers` stay plural because the singular is the commonest local
  variable name (`event`, `job`, `handler`); mass nouns stay as they are
  (`progress`, `dashboard`).
- Renames, where the old name collided or said little: `blogposts` to `blog`,
  `staffroles` to `permission`, `sessions` to `auth`, `courses` and
  `coursecategories` merged into `catalog`, `courselessons` merged into
  `curriculum`, `leads` to `leadcapture`, `reviews` to `coursereview`, `banners`
  to `sitebanner`. Every other plural is singularized (`account`, `book`,
  `coursepage`, `landingpage`, `lessonreview`, `qnaitem`, `relatedcourse`,
  `emailtoken`, `versionbuild`, `accountemail`, `feed` and the rest).
- Type names say the role (`Editor`, `Directory`, `Registrar`), not `Service`,
  `Manager` or `Handler`.
- A module exports no interface and its constructor returns the concrete type.
  `handlers.Server` holds those concrete types.
- An interface is declared by the adapter that needs it, and only for a test
  substitute or a second real implementation: the LLM, the exercise runner, the
  mailer and the course fetcher (their moq mocks live in `testsupport`), and in
  `mcpserver` the `Verifier` and `Catalogue`, which tests replace with fakes.
- Where the dependency direction forbids naming the concrete type, the adapter
  takes functions, not an interface: `handlers` imports `mcpserver`, so
  `mcpserver.Deps` carries `MyDashboard`, `MyCourseProgress` and
  `HeldCapabilities` as function fields.

**Consequences.** A handler test no longer has a type to substitute: it builds
the real module through `testsupport`, as every test already did. Adding a
collaborator is a field of a concrete type in `handlers.Deps` and one line in
`di`.
