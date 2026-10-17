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
away works at once. A Connected App's Scope can only narrow it (ADR-0019).
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

## Addendum: the matrix is retired

Settled by the slice that moved every remaining resource, after blog posts.

**What moved.** Banners, reviews, leads, messages, course categories, courses,
course lessons, course lesson reviews, lesson progress and course landing pages
declare their Permissions the same way as blog posts, named with the
[GLOSSARY](../../GLOSSARY.md) terms (`courses.read`, `course_categories.update`,
`lesson_progress.read`, `course_landing_pages.create`); the legacy
`language_*` values stop being identifiers. Actions follow the old bits: `read`
(list, get), `create`, `update` (update and every operation the old policy let
`can_update` open, which is the review of a course or a lesson) and `delete`. A
bit that opened no operation (a delete on courses, a create on lessons) has no
Permission. One migration per resource group turns each bit into the grant for
exactly those operations, and a test runs every migration over each combination
of a resource's four bits and asks the policy about each operation, so nobody
gains or loses one.

**Operations nothing delegates yet.** An admin operation that declares no
Permission declares an `x-access` extension instead, `administrators` for users,
staff members, roles and role permissions, question-and-answer items, course
versions and the permissions catalog, `any-staff` for the attachment upload. It is
read from the same parsed contract as `x-permission`, so it is visible in
`openapi.yaml`, and no Go table holds it. `TestEveryAdminOperationIsClassified`
fails for an operation that declares neither; declaring both is a parse error.
Replacing the marker on an operation with a Permission is a one-line change in the
TypeSpec, which is how the next slice makes these operations delegable.

**Departures from the text above: none, with one decision.** The old
create/read/update/delete storage is dropped, not kept empty: the table
`staff_member_role_permissions`, the ent entity, the `PermissionResource` enum, the
matrix fields of the roles contract and the matrix rows of the roles screen are
gone, and `PUT /admin/management/role_permissions/{roleId}` takes only `grants`.
Leaving an empty copy would invite a second place to hold a right. Rows the
migrations did not convert named a resource no operation used, so they carried no
access. The migration tests recreate the legacy table inside their own transaction
to feed the conversions.

## Addendum: how an Audit Record joins the business transaction

Settled by the audit tracer slice, blog posts. It closes the "Not designed yet"
item under Consequences. It meets the ADR as written: a change and its record
share one transaction, from one place, shared by the website and MCP; no
departure.

**The problem.** The generated server hands a handler a context, not a
transaction, and handlers query through `s.db`, an autocommit client. A record
written "in the same transaction" needs the handler's own writes to run on a
transaction the recorder also holds.

**One unit of work, selected by the contract.** `handlers.NewAuditMiddleware`
is an ogen middleware (`api.WithMiddleware`, wired in `di` and in the test
harness). It runs after security and body decoding, and only for an operation
whose declared Permission (`x-permission`) has kind `write`; everything else
passes through. For such an operation it opens one transaction through
`store.Transactor`, runs the handler, and on success writes the Audit Record
through the same transaction before commit. A handler error rolls everything
back, so a change that failed leaves no record; a record that cannot be written
fails the operation, so a change never commits without one. A resource that
declares its Permissions in TypeSpec is audited with no code of its own, which
is the whole mechanism later tickets inherit: #845 only has to declare the
remaining admin writes.

**How the handler joins without knowing.** `store.WithAmbientTx` puts the
transaction in the context, and two things honour it:

- the ent client built by `store.NewClient` is wrapped in a routing driver that
  sends a statement to the context's transaction when there is one and to the
  pool when there is not (ent hands the caller's context to every `Exec`,
  `Query` and `Tx`). A handler keeps writing `s.db.BlogPost.Create()`, and its
  read-after-write (`getAdminBlogPost`) sees the write;
- `Store.WithinTx` called with such a context joins it instead of opening a
  second transaction, so `blogposts.Remover`, `relatedcourses.Replacer` and the
  other modules that own an atomic step become part of the operation's commit.
  They keep their own atomicity when called outside one (a worker, a test).

Considered and not chosen: a `s.dbFor(ctx)` helper replacing `s.db` in every
write handler (per-handler code, and a forgotten call silently writes outside
the transaction); a decorator around `api.Handler` (needs a method per
operation, regenerated with the contract); an ent hook alone (it can observe a
mutation but cannot make the write share a transaction).

**What a record holds, and how changes are captured.** `audit.Hook` is an ent
hook installed by `store.NewClient` and `store.NewTxClient`; it does nothing
unless an `audit.Recorder` is in the context. With one, it notes every
mutation: entity type, id, operation and, for each field, the old and new value
as JSON. An update of one row reads the old values first (ent allows that for
`UpdateOne`) and leaves out a field set to what it already was; `created_at`
and `updated_at` are ent's own and are skipped. A delete of one row records
entity and id, not old values, because ent offers no generic read of the row
being deleted; a bulk update or delete (a `Where`) is recorded without id. The
record is filed under the first single-row change (the post, not the likes
deleted before it), and one operation is always one record however many rows it
touched. An operation that changes no row still leaves its record.

**Secrets.** The marker is ent's own `Sensitive()`, not a new annotation: it
already means "keep this out of output", `User.password_digest` and
`confirmation_token` carry it, and a second marker could disagree with it. An
ent template (`ent/template/auditsecret.tmpl`) generates `ent.AuditSecretFields`
from the schema, so marking a field is enough. A secret field appears in
`changes` as `{"name": ..., "secret": true}` with no values, and only when it
changed.

**Who, where, how.** The actor is the authenticated user of the context; the
address is `WithClientIP`'s; user agent and request id come from the HTTP
request (the request id is `X-Request-ID` if sent, else the OpenTelemetry trace
id, else a random one); the locale is the `Accept-Language` match, unchanged;
the channel is `web`, from `audit.SourceFrom(ctx)`. The MCP server sets
`audit.WithSource` (channel `mcp`, Connected App and token ids) and runs its
tool calls through the same middleware, so it records the same way; the columns
for both ids exist and are nullable. The names in code are Audit Record, never
"event".

**The table.** `audit_records` (migration `audit_records`): no foreign keys,
since a record outlives the user, app and entity it names; primary key
`(id, occurred_at)` so the table can become monthly-partitioned (#846) without
a rewrite, which this slice does not do, nor retention. A
row-level `BEFORE UPDATE OR DELETE` trigger raises an error, in the migration
because production Go may not issue raw SQL (`forbidigo`); a row trigger is
cloned onto partitions, a statement trigger is not. Retention drops whole
partitions, which the trigger does not see. A test proves the database refuses
both statements.

**Left for the next tickets, and not designed out.** Refusals never reach the
middleware (they are rejected in `HandleAdminSession`), and a handler failure
rolls back its own record; `audit.Write` already takes an outcome and an
`Entry` without a user, so #845 writes `forbidden` and `failed` records from
the security handler and from the middleware's error path, each in its own short
transaction. Reads of personal data (#850) go through Watermill, not this path.

## Addendum: monthly partitions and retention of the audit table

Settled by the partitioning ticket (#846). It meets the ADR as written (retention
by dropping whole partitions, the application unable to touch a single record);
no departure.

**The table.** A second migration (`partition_audit_records`) turns
`audit_records` into a table partitioned by month on `occurred_at`: a new
partitioned parent with the same columns, `(id, occurred_at)` key, checks,
indexes and row-level append-only trigger, the data copied over, the old table
dropped. The first migration is untouched. The ent schema and the insert path
are unchanged: ent inserts into the parent and PostgreSQL routes the row. A
`audit_records_default` partition is the safety net, so an insert never fails
for a missing month; it is never dropped by retention.

**Where the DDL lives.** `forbidigo` bans raw SQL in production Go and ent has
no partition DDL, so the DDL is two SQL functions created by the migration,
`audit_records_ensure_partition(month)` and
`audit_records_drop_partitions_ended_before(cutoff)`. `audit.Partitions` calls
each with one fixed statement (a single `//nolint:forbidigo` with the reason).
Both take an advisory lock, so two workers cannot race. Considered and not
chosen: enabling ent's `sql/execquery` feature (regenerates the whole client to
expose a general raw-SQL door), and DDL assembled as Go strings.

**The job.** `maintain_audit_partitions` is a periodic River job (daily, and on
start) in the worker process: it creates the current month and the next three
(`audit.PartitionsAhead`), then drops expired partitions. If rows already sit in
the default partition for a month being created (the job did not run in time),
the function moves them into the new partition: PostgreSQL cannot create an
overlapping partition, so it detaches the default partition (which removes its
cloned trigger), moves the rows unchanged and reattaches it (which clones the
trigger back). Only this function can do that; the application has no such
statement.

**Retention.** `AUDIT_RETENTION_MONTHS` (default 24). A partition is dropped
only when its whole month ended at or before `now - retention`, so the month
holding the cutoff is kept whole: a record lives at least the retention and up
to a month longer, and the month in progress is never dropped. `DROP TABLE` does
not fire the trigger; it is the one permitted way a record leaves the table.
Tests prove the trigger still refuses UPDATE and DELETE on rows in a partition.

## Addendum: delegable Permissions, review, and capability Permissions

Settled by the delegation ticket (#843). It meets the ADR as written; no departure.

**Administrators-only is gone.** Every operation that declared
`x-access: administrators` declares an `x-permission` now: `users.*` (users and
the management users), `staff_members.*`, `staff_roles.*` (roles, the Permission
catalog and a role's grants), `course_landing_page_qna_items.*`,
`course_category_qna_items.*` and `course_versions.create`. The parser accepts
only `any-staff` as an `x-access` (the attachment upload); `administrators` is a
parse error. An operation that declares nothing is still closed to staff, so a
new one is closed until classified. Only administrators held these operations
before, so no grant is migrated and no role gains anything; a test shows the
seeded editor role is refused every one of them, and an administrator keeps them.
Delegating `staff_roles.update` lets its holder edit any role's grants: that is
the point of delegating it, and it is a Permission to hand out with care.

**Review is separate from update.** The contract has no publish operation (blog
posts publish through the update), so only review is split: `courses.review`
(`AdminReviewCourse`) and `course_lessons.review` (`AdminReviewCourseLesson`).
The migration `review_permissions` grants each to every role holding the matching
update, so nobody gains or loses access; `course_lessons.update`, which opened
only the review, is replaced by the review grant instead of left as a right that
names nothing. A test runs the migration over each combination of the two update
grants and asks the policy about update and review. A publish Permission arrives
with the first publish operation.

**How capability Permissions are declared.** A Permission nothing requires yet
cannot be read off an operation, so the TypeSpec service lists them under a
root-level extension on the `HexletBasics` namespace, emitted at the root of
`openapi.yaml`:

```
@extension("x-capabilities", #[#{ name: "mcp.read", kind: "read" }, ...])
```

`permissions.Parse` reads it into the same catalog as the operation Permissions
(kind checked, a duplicate or a clash with an operation's kind is a parse error),
so the roles screen and the grant validation need nothing new. The catalog's
`Permission` carries `capability: true` for them, which the screen shows as an
"explicit grant only" badge. The four are `mcp.read`, `mcp.write`, `mcp.manage`
and `audit.read`; #854 (MCP) and #853 (audit page) require them.

**The administrator flag does not imply a capability.** `Policy.Allow` lets an
administrator through every operation Permission, but an operation that declares
a capability Permission (the audit page later) needs the grant on the user's
staff role even from an administrator, who then needs a staff record to hold it.
`Policy.HasCapability(ctx, user, capability)` is the check the MCP layer calls: it
looks only at the role's grants, never at the flag, and a name that is not a
capability is an error of the caller, not a refusal. No operation requires a
capability yet, so that branch of `Allow` is exercised through `HasCapability`
alone until #853 declares one.

## Addendum: every admin write, and every refusal or failure (#845)

Settled by the ticket that closed the audit trail. One departure from the text
above, in the last section: an operation that calls an external service is not
wrapped in the transaction its record shares.

**Which operations are admin writes.** `permissions.IsAdminWrite`, read from the
parsed contract, no Go list. An operation is on the admin surface when it
requires the `AdminSession` security scheme or declares an `x-permission` or
`x-access`. It is a write when its declared Permission has kind `write`; an
operation that declares only an `x-access` marker (or, until it is classified,
nothing) is a write when its HTTP method is not GET, HEAD or OPTIONS. The
audit middleware acts on exactly that set, so an operation behind an
`x-access` marker (only the any-staff upload remains since #843 gave the
administrators-only ones Permissions) is recorded, as those were before the
change, because the rule never looked at the marker. `TestEveryAdminWriteOperationOfTheContractIsRecorded` runs the
middleware once for each name `permissions.AdminWrites()` returns and fails if
one leaves no record, so a new operation is covered by declaring it in TypeSpec.

**Refusals.** A refusal is rejected in `AuthHandler.HandleAdminSession`, before
any handler or transaction exists. There, a `staffroles.ErrForbidden` writes an
Audit Record with outcome `forbidden` as one insert on the pooled client, which
commits on its own, and then answers 403 as before. This covers every admin
operation, reads included, because a refused read of users or leads is as
worth seeing as a refused change. A signed-out call (401) is not a refusal by
the permission check and leaves nothing: it has no actor. A record that cannot
be written is logged and does not change the answer. The request's user agent
and request id, which a security handler cannot read from a context, now ride
in it next to the client address (`WithClientIP`).

**Failures.** When the handler returns an error, or the record cannot be written,
the transaction rolls back and the middleware then writes an Audit Record with
outcome `failed` through `store.Transactor` in a transaction of its own, on a
context stripped of cancellation so a client that went away does not stop the
record of what was tried. A handler answer of `*api.ValidationError` is a
response, not an error, and changes nothing; it rolls back and is recorded as
`failed` too, not as a change that succeeded. The caller still gets that answer.

**What a refusal or failure record holds.** Who, channel, address, user agent,
request id, locale, operation and outcome; never `changes`, entity or error text.
The changes did not happen, and a failed call may be one the caller was not
allowed to make, so values they were not allowed to read cannot reach the record
through it. A test pins `changes`, `entity_type` and `entity_id` empty for both
outcomes, with a body the caller sent that must not appear.

**Deletes keep their old values.** ent offers no generic read of a row about to
be deleted, so an ent template (`ent/template/auditrow.tmpl`) generates
`(*Client).AuditRow(ctx, typeName, id)`, which loads one row of any schema type
and returns its columns by name. The hook (`audit.Hook(loader)`, installed with
the client that owns the hook) calls it before a single-row delete, with the
delete's own context, so inside the same transaction; each non-empty field
becomes an `old` value, a secret by name only (`AuditSecretFields`), ent's own
timestamps skipped. A delete of a row that is not there records nothing. A
bulk update or delete (`Where`) still has no row to load and records no
values, which is why the dependents of a post (its likes) show as entity and
operation only while the post itself is complete.

**Transactions around external calls.** The only admin write that does I/O in
its handler is the attachment upload (blob storage). The AI related-courses
suggestion, declared `blog_posts.update`, does not call the model in the handler:
it enqueues a River job, and the worker calls it later with no request
transaction. So the middleware's single transaction was held across a storage
call for the upload alone. The smallest fix is a contract declaration,
`x-unit-of-work: standalone` (read by `permissions.Standalone`), on that
operation. A standalone operation runs with no transaction: its writes commit as
they happen (the handler already compensates the blob on a failed insert), the
recorder still notes them, and the Audit Record is written afterwards in its own
short transaction. **This departs from "a change and its record are written in
the same transaction" for these operations**: a crash between the change and the
record loses the record. It is accepted for uploads, where the row is metadata
for a file nobody can reference before the response returns, and the failure
and refusal paths still record them; every other admin write keeps one
transaction. Any later operation that calls a mailer, a model or storage inside
its handler declares the same marker rather than holding a connection across the
call.

**Jobs.** `jobs.Enqueuer.Enqueue` goes through `Store.WithinTx`, which joins the
ambient transaction, so a job enqueued by an admin write commits with it and is
discarded with it. A test makes the record fail after the suggestion job was
enqueued and finds no job.

**Left as is.** Reads of personal data (#850) are not on this path. The stored
`failed` record carries no reason; adding an error class would need a column and
a decision on what may be named.
