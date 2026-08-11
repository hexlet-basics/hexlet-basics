# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Before exploring, read these

- **`CONTEXT.md`** at the repo root, or
- **`CONTEXT-MAP.md`** at the repo root if it exists — it points at one `CONTEXT.md` per context. Read each one relevant to the topic.
- **`docs/adr/`** — read ADRs that touch the area you're about to work in. In multi-context repos, also check `src/<context>/docs/adr/` for context-scoped decisions.

If any of these files don't exist, **proceed silently**. Don't flag their absence; don't suggest creating them upfront. The `/domain-modeling` skill (reached via `/grill-with-docs` and `/improve-codebase-architecture`) creates them lazily when terms or decisions actually get resolved.

## File structure

This repo is **single-context**: one glossary and one ADR set at the root.

```
/
├── CONTEXT.md                         ← domain glossary
├── docs/adr/                          ← 0001-… binding architecture decisions
└── internal/, ent/, api-spec/, src/
```

`legacy/` is the Rails app being replaced. It carries its own `AGENTS.md` and is
out of scope for these domain docs — don't mine it for vocabulary, and don't add
a second `CONTEXT.md` there.

If this ever becomes multi-context, the marker is a `CONTEXT-MAP.md` at the root
pointing at one `CONTEXT.md` per context, with context-scoped `docs/adr/`
alongside each.

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, a test name), use the term as defined in `CONTEXT.md`.

`CONTEXT.md` opens with a **Naming rule** that says where the term binds: everywhere the name is ours to choose — Go and TypeScript identifiers, contract models and fields, ent schema and edge names, i18n keys, docs, tests, open issues — with three permanent exceptions (table names, column names, public URLs) and one that expires at the cutover (`PermissionResource` values). Read that section before proposing a rename; the exceptions are what stops a sweep from moving storage.

If the concept you need isn't in the glossary yet, that's a signal — either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/domain-modeling`).

### This repo's glossary format differs from the skill's

`/domain-modeling` ships a `CONTEXT-FORMAT.md` that prescribes an `_Avoid_` line per entry. This repo deliberately dropped those: the Naming rule replaced them, and an entry that only warns about a retired word is a warning, not a definition. Don't reintroduce `_Avoid_` lines when editing `CONTEXT.md`.

Two additions the skill's format doesn't have, both load-bearing:

- **`_ru_` / `_es_` lines** on terms that reach a screen, so the same concept is not translated two ways across the three locales.
- **A storage mapping table** at the end, recording which legacy table each concept kept. It is the only place that mapping exists for concepts whose ent schema hasn't been written yet.

## Flag ADR conflicts

If your output contradicts an existing ADR, surface it explicitly rather than silently overriding:

> _Contradicts ADR-0001 (contract-first pipeline with ogen) — but worth
> reopening because…_
