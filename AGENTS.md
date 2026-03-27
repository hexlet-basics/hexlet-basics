# AGENTS Guide

This file is for coding agents working in `hexlet-basics`.
Prefer these instructions over generic Rails or React defaults.

## Project Overview
- Backend: Ruby on Rails app in `app/`.
- Frontend: InertiaJS + React 19 + TypeScript + Vite in `app/javascript/`.
- Tests: Minitest for Ruby, Vitest available for frontend.
- Package manager: `pnpm`.
- Main orchestration: `make` targets in `Makefile`.

## Instruction Sources Checked
- Existing repo guidance found in `/Users/kirillmokevnin/projects/hexlet-basics/AGENTS.md` and replaced with this expanded version.
- No `.cursorrules` file found.
- No files found under `.cursor/rules/`.
- No `.github/copilot-instructions.md` file found.
- If any of those files are added later, merge their rules into this document instead of ignoring them.

## Repository Layout
- `app/controllers/`: Rails controllers, many render Inertia pages.
- `app/models/`: Active Record models and form objects.
- `app/services/`, `app/lib/`, `app/handlers/`: service and integration logic.
- `app/resources/`: serializer/resource objects returned to the frontend.
- `app/javascript/pages/`: route-level Inertia pages.
- `app/javascript/components/`: shared React components.
- `app/javascript/hooks/`: reusable frontend hooks.
- `app/javascript/types/`: shared TS types and generated serializer types.
- `config/routes.rb`: main route map; check it before adding endpoints or page paths.
- `test/`: Minitest unit, controller, job, mailer, and system tests.

## Setup And Dev Commands
- `make setup` - copy `.env.local`, run `bin/setup --skip-server`, load fixtures, install JS deps, install git hooks.
- `make dev` - run the local app stack via `overmind` and `Procfile.dev`.
- `make staging` - run the staging-like local stack.
- `make services-start` / `make services-stop` - start or stop local Postgres in Docker.
- `make db-reset` - reset DB and reload fixtures.
- `make sync` - sync locales, fixtures, TS types, routes, browserlist data.
- `make build-assets` - precompile Rails/Vite assets.

## Lint And Typecheck Commands
- `make lint` - full local lint gate; runs locale checks, Biome, TypeScript build, and RuboCop.
- `make lint-fix` - auto-fix RuboCop and Biome issues where possible.
- `pnpm exec biome check app/javascript/path/to/file.tsx` - lint a single frontend file.
- `pnpm exec biome check --fix app/javascript/path/to/file.tsx` - auto-fix one frontend file.
- `pnpm tsc --build` - run the TypeScript project build check.
- `bin/rubocop` - run Ruby linting.
- `bin/rubocop app/controllers/web/users_controller.rb` - lint one Ruby file.

## Test Commands
- `make test` - run the Rails test suite.
- `make test-system` - run Rails system tests.
- `make test-frontend` - run Vitest once.
- `make test-frontend-watch` - run Vitest in watch mode.
- `make test-all` - run backend, frontend, and system tests.

## Single Test Commands
- `bin/rails test test/controllers/web/sessions_controller_test.rb` - run one Ruby test file.
- `bin/rails test test/controllers/web/sessions_controller_test.rb:11` - run one Minitest by line number.
- `bin/rails test test/models/lead_test.rb:1` - use line targeting when the file has many tests.
- `bin/rails test:system test/system/pages_test.rb` - run one system test file.
- `pnpm vitest run app/javascript/path/to/file.test.ts` - run one Vitest file.
- `pnpm vitest run -t "test name"` - run Vitest tests matching a name.
- `pnpm vitest app/javascript/path/to/file.test.ts` - run one frontend file in watch mode.

## High-Value Agent Workflow
- Before editing, read nearby files and follow local conventions instead of introducing new patterns.
- Prefer focused changes; this codebase mixes Rails, Inertia, generated types, and event code.
- After changing Ruby code, run the narrowest relevant Rails test first.
- After changing TypeScript or React code, run Biome and `pnpm tsc --build` at minimum.
- Before finishing substantial work, run `make lint` plus the most relevant tests if time permits.
- Do not edit generated outputs casually; regenerate them with the provided sync tasks when needed.

## Ruby And Rails Style
- Follow RuboCop Rails Omakase; the repo inherits `rubocop-rails-omakase`.
- Ruby target version is 4.0.2.
- Most Ruby files use double-quoted strings and `# frozen_string_literal: true` when present; preserve existing style in each file.
- Use small controller actions that prepare data, set flash/meta, and delegate domain work to models, services, handlers, or events.
- Keep business logic out of views.
- Prefer expressive names like `signed_up_event_data` over short temporary names.
- Use guard clauses and small helper methods when actions become dense.
- Respect existing concern usage such as authentication, locale, browser, flash, and event concerns.

## Ruby Types And Contracts
- Sorbet is used in parts of the codebase, especially events, handlers, services, and some forms.
- Preserve existing `# typed:` sigils.
- Add `sig` blocks only where the surrounding code already uses Sorbet or where the file is clearly typed.
- Use `T.must` only when the invariant is genuinely guaranteed nearby.
- Do not start broad Sorbet migrations in unrelated work.

## Ruby Error Handling
- Rescue specific exceptions instead of broad `StandardError` unless the file already has a deliberate boundary catch.
- Keep rescue logic close to the failure boundary, as seen in controllers that recover from validation or uniqueness errors.
- Prefer validation errors, result objects, or redirects with flash over silent failure.
- Avoid swallowing exceptions without user-visible handling, logging, or a clear fallback.
- Preserve existing failure flows such as `flash.inertia`, redirects, and form rehydration.

## Rails Controller And Resource Conventions
- Many HTML routes render Inertia pages with `render inertia: true, props: { ... }`.
- Backend-to-frontend payloads often go through resource classes in `app/resources/`.
- When adding new page props, keep them serializable and shaped for TS consumers.
- Reuse named route helpers from Rails and generated JS routes rather than hardcoding paths.
- Check `config/routes.rb` before creating new controllers or page locations.

## Frontend Style
- TypeScript is `strict`; write code that passes without relying on implicit `any`.
- Use the `@/` alias for imports from `app/javascript/`.
- Keep `import type` statements separate for type-only imports when possible.
- Existing frontend formatting uses Biome defaults with spaces and single quotes in many files; run Biome instead of hand-formatting.
- Prefer functional React components.
- Route-level pages commonly use default exports; preserve file-local convention.
- Shared reusable components may use either named or default exports; follow the surrounding folder's pattern.
- Prefer existing hooks and helpers such as `useAppForm`, `useDataTableProps`, and route helpers before introducing new abstractions.
- Do not invent a second path alias system.

## Frontend Naming And Data Conventions
- Match current naming patterns: `Props`, `payload`, `form`, `grid`, `pagy`, `...Dto`, `...Resource`, `...Event`.
- Keep page component names aligned with route semantics such as `Index`, `Show`, `New`, and `Edit` when the folder already uses them.
- Prefer serializer/resource-backed types from `@/types` or `@/types/serializers` for server data.
- Reuse generated route helpers from `@/routes.js` instead of string URLs.
- Reuse existing i18n patterns with `useTranslation()` and selector callbacks like `t(($) => $.sessions.new.title)`.

## Frontend Error Handling
- Prefer form-library callbacks like `onError` and server-provided validation state over ad hoc try/catch in components.
- Surface request failures in the established UI flow, not only in console logs.
- Reset or sanitize sensitive fields such as passwords on failed auth forms when matching existing behavior.
- Keep client code aligned with Rails response shapes and Inertia redirects.

## Generated And Synced Files
- `app/javascript/routes.*` is generated; do not hand-edit unless the repo clearly expects it.
- `app/javascript/types/serializers/` contains generated serializer types.
- Some locale and type outputs are regenerated by `make sync` and `make sync-locales` / `make sync-types`.
- If a change affects routes, events, enums, or serializer-backed TS types, regenerate the relevant outputs.

## Testing Guidance
- Prefer targeted tests first, then broader suites.
- For controller or model changes, add or update a nearby `_test.rb` file under `test/`.
- Use fixtures from `test/fixtures/`; the suite loads all fixtures by default.
- System tests are reserved for user-visible flows that need browser coverage.
- Frontend tests are not widespread today; if you add Vitest coverage, keep it narrow and fast.
- Do not leave failing tests behind if the touched area can be tested locally.

## Git And Delivery Notes
- Pre-push hook runs `make lint`; expect pushes to require a clean lint state.
- Keep PRs and commits focused; mention migrations, generated artifacts, or infra changes explicitly.
- Never commit secrets from `.env.local`, credentials files, or external service configs.
- If you touch ops directories like `k8s/`, `ansible/`, or `terraform/`, call that out clearly in handoff notes.

## Practical Defaults For Agents
- Use `make` targets first when a task maps cleanly to them.
- Prefer minimal diffs over broad refactors.
- Preserve mixed typed and untyped Ruby; do not "normalize" unrelated files.
- Preserve Inertia + Rails patterns instead of replacing them with API-only patterns.
- When unsure whether something is generated, search for the producing task before editing.
