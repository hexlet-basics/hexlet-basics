# Contributing to Hexlet Basics

Thanks for your interest in contributing! Hexlet Basics is a platform of free
introductory programming courses — each course teaches one programming language
through a sequence of in-browser coding lessons.

This guide covers how to set up the project, the conventions we follow, and how
to get a change merged. It complements two other docs you should skim:

- [`README.md`](./README.md) — quick setup reference.
- [`AGENTS.md`](./AGENTS.md) — the detailed engineering conventions (also used
  by coding agents). It is the source of truth for style and command details.
- [`CONTEXT.md`](./CONTEXT.md) — the project's domain glossary (Course, Lesson,
  Enrollment, …). Read it before naming things or discussing the model.

Developer-facing docs, code comments, and commit messages are written in
**English**.

## Tech stack

- **Backend**: Ruby on Rails (`app/`). Ruby `4.0.6` (see `.ruby-version`),
  Rails `8.1`, PostgreSQL.
- **Frontend**: Inertia.js + React 19 + TypeScript + Vite (`app/javascript/`),
  Mantine UI, Monaco editor.
- **Package manager**: `pnpm` (Node `25.2.1`, see `.node-version`). Do not use
  `npm`/`npx`.
- **Tests**: Minitest (Ruby), Vitest (frontend).
- **Lint/format**: RuboCop (Rails Omakase) + Sorbet for Ruby, Biome +
  TypeScript for frontend.
- **Course content** lives in separate `exercises-<slug>` repositories under
  the [hexlet-basics org](https://github.com/hexlet-basics); this repo imports
  and renders it.

## Getting started

Requirements: `make`, `docker`, Ruby `4.0.6`, Node `25.2.1`.

1. Install system dependencies:

   ```bash
   make setup-ubuntu   # Ubuntu
   make setup-macos    # macOS
   ```

2. Start local services (PostgreSQL in Docker):

   ```bash
   make services-start
   ```

   The local Postgres is published on port `54330` so it does not conflict with
   a database already running on `5432`.

3. Set up the app (copies `.env`, installs deps, loads fixtures, installs git
   hooks):

   ```bash
   make setup
   ```

4. Run the dev stack:

   ```bash
   make dev
   ```

   Open <https://code-basics.localhost>. Port `3000` is the internal Puma port —
   do not open `https://localhost:3000` directly; local HTTPS is served through
   Caddy at `https://code-basics.localhost`.

   Admin login for local development: `full@test.io` / `password`.

5. Load a course locally when you need real lesson content:

   ```bash
   make language-load L=php
   ```

## Development workflow

### Everyday commands

- `make dev` — run the app (via `overmind` and `Procfile.dev`).
- `make test` — run the Rails test suite.
- `make test-frontend` — run Vitest once (`make test-frontend-watch` to watch).
- `make test-system` — run Rails system (browser) tests.
- `make test-all` — backend + frontend + system tests.
- `make lint` — the full local gate: locale checks, Biome, TypeScript build,
  RuboCop, and Sorbet. **A push must pass this** (the pre-push hook runs it).
- `make lint-fix` — auto-fix RuboCop and Biome issues where possible.
- `make sync` — regenerate synced/generated outputs (locales, TS types, JS
  routes, event exports, Sorbet RBIs). Run it when your change touches routes,
  events, enums, locales, or serializer-backed types.

Run a single test:

```bash
bin/rails test test/controllers/web/sessions_controller_test.rb      # one file
bin/rails test test/controllers/web/sessions_controller_test.rb:11   # one test
pnpm vitest run app/javascript/path/to/file.test.ts                  # one Vitest file
```

### Making a change

- Read nearby files first and follow local conventions rather than introducing
  new patterns. This codebase deliberately mixes Rails, Inertia, generated
  types, and event code.
- Keep changes focused; prefer minimal diffs over broad refactors.
- After Ruby changes, run the narrowest relevant Rails test, then `bin/rubocop`.
- After TypeScript/React changes, run `pnpm exec biome check` and
  `pnpm tsc --build` at minimum.
- Do not hand-edit generated files (e.g. `app/javascript/routes.*`,
  `app/javascript/types/serializers/`). Regenerate them with `make sync`.
- Before finishing substantial work, run `make lint` plus the relevant tests.

### Style

Ruby and Rails follow RuboCop Rails Omakase. Frontend is TypeScript `strict`
with Biome formatting (use the `@/` import alias for `app/javascript/`). See
[`AGENTS.md`](./AGENTS.md) for the full style guide — Ruby error handling,
Sorbet usage, controller/resource conventions, and frontend naming/data
conventions.

## Commit messages

Every commit follows [Conventional Commits](https://www.conventionalcommits.org/):
`<type>(<scope>?): <description>`. `commitlint` enforces the allowed types via a
git hook, and messages are written in English. The body should explain *why*,
not *what* (the diff already shows what).

Common types and their release impact (releases are cut by
[release-please](https://github.com/googleapis/release-please) from commit
history — see `release-please-config.json`):

| Type | Meaning | Changelog |
| --- | --- | --- |
| `feat` | User-visible feature | Features |
| `fix` | User-visible bug fix | Bug Fixes |
| `perf` | User-visible speedup | Performance Improvements |
| `refactor` | Code change, no behavior change | — |
| `test` | Tests only | — |
| `docs` | Documentation only | — |
| `chore` / `build` / `ci` / `style` | Routine work | Miscellaneous / Build System |

Mark breaking changes with `type!:` or a `BREAKING CHANGE:` body line.

Examples:

```
feat(lessons): let users resize the theory and editor panes
fix(blog): sync URL symmetrically while scrolling blog posts
```

## Opening a pull request

1. Make sure `make lint` and the relevant tests pass locally.
2. Keep the PR focused on a single concern. Explicitly call out migrations,
   generated artifacts, or infra changes (`k8s/`, `ansible/`, `terraform/`) in
   the description.
3. Reference the issue the PR addresses.
4. Never commit secrets — `.env.local`, credentials files, or external service
   configuration.

Releases (version bumps in `package.json`, `k8s/app-chart/Chart.yaml`, and
`CHANGELOG.md`) are produced automatically by release-please when its release PR
is merged; you do not need to bump versions or create tags by hand.

## Reporting issues

Open an issue on [GitHub](https://github.com/hexlet-basics/hexlet-basics/issues).
For bugs, include steps to reproduce, what you expected, and what happened.
Before filing, search existing issues to avoid duplicates.

Note that **lesson content** (theory text, exercises, tests) is not in this
repository — it lives in the per-course `exercises-<slug>` repos under the
[hexlet-basics org](https://github.com/hexlet-basics). Content fixes belong
there; this repo hosts the platform that imports and renders that content.

## License

By contributing, you agree that your contributions will be licensed under the
project's [GNU General Public License v3.0](./LICENSE).
