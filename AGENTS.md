# Repository Guidelines

## Project Structure & Module Organization
- Backend: Ruby on Rails in `app/` (models, controllers, jobs, views, policies).
- Frontend: InertiaJS/React/TypeScript via Vite in `app/javascript/`.
- Tests: Ruby tests in `test/` (unit, integration, system); frontend tests use Vitest.
- Config & ops: `config/`, `k8s/`, `ansible/`, `terraform/`, `Procfile.*`.
- Assets/Public: `public/`; built assets via Vite and Rails assets pipeline.

## Build, Test, and Development Commands
- `make setup`: bootstrap app (env, gems, DB fixtures, npm, hooks).
- `make dev`: run the app locally (Overmind/Hivemind/Foreman + `Procfile.dev`).
- `make test`: run Rails test suite. `make test-system`: system tests.
- `make test-frontend`: run Vitest once. `make test-frontend-watch`: watch mode.
- `make lint`: Biome + TypeScript build checks + RuboCop. `make lint-fix`: autofix.
- `make db-reset`: reset DB and reload fixtures. `make sync`: sync i18n, fixtures, types.
- `make services-start|services-stop`: start/stop local Postgres via Docker.
- `make build-assets`: precompile assets. `make coverage-open`: open coverage report.

## Coding Style & Naming Conventions
- Ruby: RuboCop Rails Omakase.
- JS/TS: Biome formatter (single quotes, spaces); ESM modules.

## Testing Guidelines
- Rails: Minitest in `test/` with `_test.rb` suffix; use fixtures in `test/fixtures/`.
- Frontend: Vitest for TS/React. Prefer fast unit tests; mock network and Cable.
- Aim to cover new logic and critical paths; add system tests for user‑visible flows.
- Run `make test` and `make test-frontend` before opening a PR.

## Commit & Pull Request Guidelines
- CI: ensure `make lint` and all tests pass locally (pre‑push hook runs lint).
- Keep PRs small and focused; note migrations and any ops changes.

## Security & Configuration Tips
- Copy envs from `.env.sample` → `.env`; never commit secrets.
- Local domain: `https://code-basics.localhost`. Use `make setup-macos|setup-ubuntu` for system deps.
