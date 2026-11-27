---
status: accepted
---

# Self-protection in the modules that write, and one meaning of removing an account

Amends the self-protection addendum (#844) of
[ADR-0020](0020-contract-declared-permissions-and-audit-trail.md) and the matching
paragraph of [ADR-0021](0021-mcp-server.md). The rules, their status codes and their
localized refusals are unchanged.

**Problem.** The self-protection rules were a `selfprotection.Guard` that seven
handlers called before the write it guarded, and the write itself was elsewhere:
`account.Editor` for users (which held the last-administrator rule on its own), the
generated `AdminCreate`/`AdminUpdate`/`AdminDelete` for Staff Member records, and
`permission.Permissions` for a role's grants. The order "guard, then write" lived in
each handler, nothing stopped a new caller from writing without the guard, and the
Staff Member writes went through generated helpers although they carry a rule, which
the use-case module convention forbids. Separately, removing an account had two
meanings: the profile's own removal kept the row and erased what identifies the
person (legacy `mark_as_removed!`), while `adminDeleteUser` deleted the row. Legacy
had no admin deletion at all, and a deleted row takes the person's progress,
reviews and memberships with it, or fails on their foreign keys.

**Decision.**

- **The rule sits with the write.** Each writer takes the acting user as a
  `domain.Actor` (id, administrator) and refuses before writing:
  `account.Editor.Create`/`Update` (the administrator flag only by an
  administrator, never taken from the last one), `account.Remover.Remove` (not the
  actor's own account, not the last administrator), `permission.StaffEditor`
  (`Create`/`Update`/`Delete`: not the actor's own record, nor one moved onto them),
  `permission.Permissions.Set` (not a non-administrator's own role, then the
  contract's grants, then the write). Each refusal is the module's own error, a
  `domain.ErrForbidden` or (the last administrator) a `domain.ErrConflict`. The
  `selfprotection` package is deleted.
- **Generated reads, module writes.** `AdminCRUD{ReadOnly: true}` generates only
  `AdminList` and `AdminGet`; Staff Member uses it, so no handler can write a record
  past the editor.
- **One removal.** `adminDeleteUser` removes the account the way the profile does:
  the row stays in the removed state with names, email, phone, password and
  linked accounts erased (`account.Remover.Remove`, the person's own being
  `RemoveOwn`). The operation keeps its name and its 204.

**Consequences.** The MCP admin tools still dispatch through the handlers, so the
rules bind them as before (ADR-0021). A handler of these operations passes the actor
and maps the refusals; a new writer of users, staff records or grants gets the rules
by calling the module, and has no generated helper to bypass them with. An account
an administrator removed can register again with the same email, like one its owner
removed.
