# Hexlet Basics — Go API + TypeSpec contract + TanStack (React) frontend.
# Rails lives in legacy/ with its own Makefile. Run these from the repo root.

DATABASE_PORT ?= 54330

.DEFAULT_GOAL := help

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'

# ---------------------------------------------------------------------------
# Setup
# ---------------------------------------------------------------------------

## prepare: install the mise-pinned toolchain (go, golangci-lint, atlas, testfixtures)
prepare:
	mise install

## install: install project dependencies (go modules + frontend/api-spec via root package.json)
install:
	go mod download
	pnpm install

## setup: full local bootstrap (toolchain + seed .env + install deps)
setup: prepare install
	cp -n .env.example .env || true

# ---------------------------------------------------------------------------
# Services (local Postgres in Docker; schema is still owned by legacy Rails)
# ---------------------------------------------------------------------------

## services-start: start local Postgres
services-start:
	docker run -d -it --rm \
		-p $(DATABASE_PORT):5432 \
		--name code_basics_postgres \
		-e POSTGRES_DB=code_basics_development \
		-e POSTGRES_PASSWORD=postgres \
		-v code_basics_pgdata:/var/lib/postgresql/data \
		postgres:17

## services-stop: stop local Postgres
services-stop:
	docker stop code_basics_postgres

# ---------------------------------------------------------------------------
# Dev (run the stack)
# ---------------------------------------------------------------------------

## dev: run API (air live-reload) + frontend (vite) together; Ctrl-C stops both
dev:
	pnpm dev:all

## dev-api: run only the Go API server with live-reload (:3001)
dev-api:
	air

## dev-web: run only the Vite frontend
dev-web:
	pnpm dev

## dev-spec: watch TypeSpec and re-emit OpenAPI on change
dev-spec:
	pnpm spec:watch

# ---------------------------------------------------------------------------
# Code generation (contract is the source of truth — see ADR-0001)
# ---------------------------------------------------------------------------

## gen: full contract pipeline — TypeSpec -> OpenAPI -> ogen server + hey-api client
gen: gen-spec gen-api gen-client

## gen-spec: compile TypeSpec to api-spec/dist/openapi.yaml
gen-spec:
	pnpm spec

## gen-api: generate the ogen Go server from OpenAPI
gen-api:
	go generate ./internal/apigen

## gen-client: generate the hey-api TS client + Query hooks from OpenAPI
gen-client:
	pnpm generate

## gen-ent: regenerate the ent ORM from ent/schema
gen-ent:
	go generate ./ent

## fixtures-import: convert legacy Rails fixtures -> Go testfixtures YAML in
## fixtures/ (non-destructive; FORCE=1 overwrites, ONLY=t1,t2 scopes). Borrows
## the Rails runtime via `rails runner`; adds no files to legacy/.
fixtures-import:
	cd legacy && RAILS_ENV=test bin/rails runner ../scripts/export_fixtures.rb

## gen-all: regenerate everything (ent + full contract pipeline)
gen-all: gen-ent gen

## tidy: prune/sync go.mod after generation or dep changes
tidy:
	go mod tidy

# ---------------------------------------------------------------------------
# Lint & typecheck
# ---------------------------------------------------------------------------

## lint: lint Go + frontend
lint: lint-go lint-web

## lint-go: gofmt check, go vet, golangci-lint (skips generated files)
lint-go:
	@test -z "$$(gofmt -l cmd internal ent/schema)" || (echo "gofmt needed:"; gofmt -l cmd internal ent/schema; exit 1)
	go vet ./...
	golangci-lint run

## lint-web: biome + tsc project build
lint-web:
	pnpm lint
	pnpm check

## lint-fix: auto-fix Go formatting and frontend lint
lint-fix:
	gofmt -w cmd internal ent/schema
	-golangci-lint run --fix
	pnpm lint:fix

# ---------------------------------------------------------------------------
# Test
# ---------------------------------------------------------------------------

# Test DB coordinates (override to target a scratch DB, e.g. DB_NAME=some_test).
DB_HOST ?= 127.0.0.1
DB_PORT ?= 54330
DB_NAME ?= code_basics_test
DB_URL  ?= postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

## migrate-new: scaffold a new empty migration to hand-author, e.g.
## `make migrate-new NAME=add_widgets`. Edit it, then `atlas migrate hash`.
migrate-new:
	atlas migrate new $(NAME) --dir "file://migrations"
	atlas migrate hash --dir "file://migrations"

## test-migrate: apply the atlas migrations (schema) to the test DB. Atlas owns
## the schema now (baseline = the retired Rails schema); this replaces the old
## structure.sql load.
test-migrate:
	atlas migrate apply --env local --url "$(DB_URL)&search_path=public"

## test-load-fixtures: load the committed fixtures/ snapshot into the test DB via
## the testfixtures CLI (provided by mise). sslmode=disable because the CLI's
## lib/pq driver defaults to requiring SSL, which the local/CI Postgres doesn't serve.
test-load-fixtures:
	testfixtures -d postgres -D fixtures -c "$(DB_URL)"

## test-prepare: ready the test DB for `make test` — apply migrations, then load
## the fixtures baseline. Run once before `make test` (CI does exactly this).
test-prepare: test-migrate test-load-fixtures

## test: run the Go test suite
test:
	go test ./...

# ---------------------------------------------------------------------------
# Build
# ---------------------------------------------------------------------------

## build: build the API binary and the frontend bundle
build: build-api build-web

## build-api: compile the Go server to bin/server
build-api:
	go build -o bin/server ./cmd/server

## build-web: build the production frontend bundle
build-web:
	pnpm build

## clean: remove build artifacts
clean:
	rm -rf bin dist

# ---------------------------------------------------------------------------
# Maintenance
# ---------------------------------------------------------------------------

## deps-update: bump Go + frontend (incl. api-spec) dependencies to latest
deps-update:
	go get -u ./...
	go mod tidy
	npx --yes npm-check-updates -u
	pnpm install
	@echo ">> deps bumped — run 'make gen-all' then 'make lint test' to verify"

## update-skills: update the project's agent skills
update-skills:
	npx --yes skills update --project --yes

.PHONY: help prepare install setup services-start services-stop migrate-new test-migrate test-load-fixtures test-prepare dev dev-api dev-web dev-spec \
	gen gen-spec gen-api gen-client gen-ent gen-all tidy \
	lint lint-go lint-web lint-fix test build build-api build-web clean \
	deps-update update-skills
