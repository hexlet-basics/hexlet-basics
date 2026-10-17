---
status: accepted
---

# Permissions declared in the contract, and an audit trail

Amends [ADR-0011](0011-contract-declared-authentication.md).

ADR-0011 declared *who is signed in* in the contract but left *what a staff
member may do* to a hand-kept table in `internal/staffroles/rules.go`, keyed by
operation name, with every operation that fits no resource marked
"administrators only". An MCP client acting for a staff member makes that table
the only thing between an agent and the whole admin surface, and nothing
records what was done. We move the requirement into the contract, model
permissions in the database as grants, and record staff actions in an audit
trail.

**Permissions in the contract.** Every admin operation declares in TypeSpec the
Permission it requires and whether it reads or writes. The catalog of
Permissions is generated from those declarations, so the requirement is visible
in OpenAPI and in the generated clients, as ADR-0011 did for sessions. An
operation without a declared Permission must fail `mise run gen:check` or a
test; the mechanism is designed in the implementation, not fixed here.
`rules.go` goes away.

**Permissions in the database.** A Staff Role holds a set of grants,
`(role, permission)`, one row for each Permission it has. This replaces the
resource-by-CRUD-bits matrix. A Staff Member keeps one Staff Role and the list
of locales they work in; several roles per person and locales per resource are
deferred until a real case needs them. Administrators, the current
`users.admin` value read from the database on every call (ADR-0011), pass every
content Permission. The administrator flag does **not** imply the `mcp.*`
Permissions or `audit.read`: an administrator's role must be given them
explicitly, like anyone's, so any administrator can be denied MCP access or
limited to reading. ADR-0002 rules out running Rails beside Go, so nothing constrains
the new shape: the existing roles and rows are carried into it by the cutover
migration.

The catalog covers what the matrix could not:

- the operations that were administrators-only (question-and-answer items,
  course versions) become Permissions a role can be given;
- actions beyond create, read, update and delete (publishing, reviewing) are
  Permissions of their own;
- using the site over MCP is a Permission like any other: `mcp.read`,
  `mcp.write` and `mcp.manage`, together with `audit.read`.

**Effective rights.** This applies to admin operations only. An admin
operation called over MCP is allowed when the user's role holds the operation's
Permission *and* the `mcp.*` Permission matching its kind (`mcp.read` for a
read, `mcp.write` for a write, `mcp.manage` for operations on users, staff and
roles). It is evaluated from the database on every call, so taking a Permission
away works at once. A Connected App's Scope can only narrow it (ADR-0016).
Catalogue and learner operations need no Staff Role: they need only the Scope
(`catalog:read`, `progress:read`) and a user as the token's subject.

**Self-protection.** Whoever holds the right to manage staff must not be able
to remove their own footing. At the operation layer, so that the website and
MCP are bound alike: a user cannot change their own role or Staff Member
record, cannot clear the administrator flag of the last administrator, and
cannot delete themselves.

**Audit trail.** Every change an administrator or staff member makes through
the admin surface, from the website or over MCP, and every read of personal
data (users, leads, messages, learners' progress), leaves an Audit Record.
Reads of the catalogue and of content are not recorded.

- A record says who acted, through which channel (website or MCP, with the
  Connected App and token for MCP), the request id, address and user agent,
  which operation, on which entity, in which locale, and the outcome
  (succeeded, forbidden, failed).
- For a change it also keeps the old and new value of each changed field;
  fields marked secret in the ent schema keep only their name.
- The table is append-only: a database trigger refuses `UPDATE` and `DELETE`.
  It is partitioned by month and records are kept two years, configurable.
- A change and its record are written in the same transaction, from one place
  in the operation layer shared by the website and MCP, so a change never lacks
  its record and a record never describes a change that did not happen. A read
  of personal data is recorded asynchronously through Watermill (ADR-0004):
  losing one such record on a crash is acceptable, slowing the response is not.
- Records are read on an admin page and by an MCP tool, both requiring
  `audit.read` (the tool also `mcp.read`). Reading the trail is itself recorded
  as a read of personal data, because records carry addresses and old values.
- It is called an Audit Record, not an event: "event" already means a Watermill
  domain event.

## Considered Options

- **Keep the matrix and add columns.** Cheapest, but each new action or
  capability is a schema change, and the operation-to-resource map stays
  hand-kept and outside the contract.
- **Casbin.** Replaces one table with another and adds a policy layer between
  ent and the code, without the contract-declared requirement we want.
- **OpenFGA or SpiceDB.** Relationship-based, worth it only for per-entity
  rights (an editor for one course), which we are not building.

## Consequences

- This reaches every admin operation and the roles screen, and has to land
  before any MCP admin tool (ADR-0018).
- Not designed yet: how the single place that writes a record inside the
  business transaction is wired through ogen, which hands handlers a context
  and not a transaction. It needs its own design before the audit work starts.
- Existing role rows must be mapped onto the new catalog by the cutover
  migration, and the legacy resource names (`language_lessons` and the like)
  stop being identifiers.

## Addendum: how a Permission is declared and enforced

Settled by the first slice, blog posts, which every later resource follows.

**Declaring.** An operation carries one extension in the TypeSpec:
`@extension("x-permission", #{ name: "blog_posts.read", kind: "read" })`. The
name is `<resource>.<action>`; the kind is `read` or `write`. A Permission has
exactly one kind, so a role can be given the right to read without the right to
change; parsing the contract fails if one name is declared with two kinds. The
extension lands on the operation in `api-spec/dist/openapi.yaml`, which is how it
shows in the contract. Blog posts use four Permissions, `read` (list, get),
`create`, `update` (update, related-courses set, AI suggestion) and `delete`,
not `read` and `write` only, because the old matrix bits opened exactly those
groups and the migration must not widen anyone's access. "Publish" and
"review" become Permissions of their own when their resources migrate.

**Reading it in Go.** ogen ignores vendor extensions, so `internal/api` does not
carry the declaration. `internal/permissions` reads it from the generated
OpenAPI document, which the small package `api-spec` embeds (`go:embed` cannot
reach up out of a package directory). Considered and not chosen: generating a Go
catalog file as a new `gen` step. It would be a second generated artifact to
keep in step with the first, with its own `gen:check` surface, for no gain: the
embedded document is already regenerated by `mise run gen:spec`, is covered by
`gen:check`, and cannot be hand-edited. The cost is parsing the document once at
start-up (a few hundred kilobytes of YAML); a document that cannot be read panics
there, because opening or closing every operation silently is worse.
`permissions.For(operationName)` answers "what does this operation declare", and
`permissions.Catalog()` lists every declared Permission once; the roles screen
gets it from `GET /admin/management/permissions`, so its list cannot drift.

**Enforcing.** Unchanged in where it happens: `AuthHandler.HandleAdminSession`
runs `staffroles.Policy.Allow` before the request body is decoded. `Allow`
passes an administrator, requires a staff record with the request's locale, and
then, if the operation declares a Permission, requires a row in
`staff_member_role_grants` for the role and that Permission; otherwise it falls
back to the retired-resource table and matrix, which keeps serving the
resources not yet migrated. An operation that declares nothing and is in no
table is still administrators-only. `TestEveryAdminOperationIsClassified`
accepts a declaration as a classification, and a contract test pins that every
blog post operation declares a Permission and a kind.

**Storing.** `staff_member_role_grants(role_id, permission)`, unique per pair,
no foreign key from the permission name to anything: the catalog lives in the
contract, not in a table. The roles screen saves matrix and grants in one
transaction; grants are validated against the catalog and replaced as a set.
The cutover migration of a resource converts each matrix row into grants and
deletes the row; the fixtures, which load after migrations, are rewritten by
hand to the same shape.

Departures from the text above: none. The locale is still the `Accept-Language`
match, not a request argument; ADR-0018's locale argument is not part of this
slice.
