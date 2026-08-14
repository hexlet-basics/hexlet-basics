# Hexlet Basics — Go API + TypeSpec contract + TanStack (React) frontend.
# Rails lives in legacy/ with its own Makefile. Run these from the repo root.

DATABASE_PORT ?= 54330
AMOCRM_API_REF := 5a24b11242c66c5b64e4a0d1c107a45ee9dfbab3
AMOCRM_OPENAPI := https://raw.githubusercontent.com/Hexlet/amocrm-api/$(AMOCRM_API_REF)/tsp-output/schema/openapi.yaml

.DEFAULT_GOAL := help

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'

# ---------------------------------------------------------------------------
# Setup
# ---------------------------------------------------------------------------

## prepare: install the mise-pinned toolchain (go, golangci-lint, kiota, atlas)
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
# Services (development Postgres in Docker; tests use testcontainers-go)
# ---------------------------------------------------------------------------

## services-start: start local Postgres
services-start:
	docker run -d -it --rm \
		-p $(DATABASE_PORT):5432 \
		--name code_basics_postgres \
		-e POSTGRES_DB=code_basics_development \
		-e POSTGRES_PASSWORD=postgres \
		-v code_basics_pgdata18:/var/lib/postgresql \
		postgres:18
	@until docker exec code_basics_postgres pg_isready -U postgres >/dev/null 2>&1; do sleep 1; done
	$(MAKE) dev-db-create

## services-stop: stop local Postgres
services-stop:
	docker stop code_basics_postgres

# ---------------------------------------------------------------------------
# Dev (run the stack)
# ---------------------------------------------------------------------------

## dev: run API + async worker (air live-reload) + frontend; Ctrl-C stops all
dev:
	pnpm dev:all

## dev-api: run only the Go API server with live-reload (:3001)
dev-api:
	go tool air

## dev-worker: run only the async Go worker with live-reload
dev-worker:
	go tool air -c .air.worker.toml

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

## gen-locales: extract the English backend catalog from the message registry
gen-locales:
	go generate ./internal/localization

## gen-amocrm: generate the typed amoCRM client from the pinned upstream contract
gen-amocrm:
	kiota generate \
		--language Go \
		--class-name APIClient \
		--namespace-name hexletbasics/internal/amocrm/generated \
		--openapi $(AMOCRM_OPENAPI) \
		--output internal/amocrm/generated \
		--include-path '/api/v4/leads/unsorted/forms#POST' \
		--serializer github.com/microsoft/kiota-serialization-json-go.JsonSerializationWriterFactory \
		--deserializer github.com/microsoft/kiota-serialization-json-go.JsonParseNodeFactory \
		--structured-mime-types application/json \
		--structured-mime-types application/hal+json \
		--structured-mime-types application/problem+json \
		--exclude-backward-compatible \
		--clean-output
	gofmt -w internal/amocrm/generated
	pnpm exec oxfmt internal/amocrm/generated/kiota-lock.json

## gen-ent: regenerate the ent ORM from ent/schema
gen-ent:
	go generate ./ent

## fixtures-import: convert legacy Rails fixtures -> Go testfixtures YAML in
## fixtures/ (non-destructive; FORCE=1 overwrites, ONLY=t1,t2 scopes). Borrows
## the Rails runtime via `rails runner`; adds no files to legacy/.
fixtures-import:
	cd legacy && RAILS_ENV=test bin/rails runner ../scripts/export_fixtures.rb

## gen-all: regenerate everything (ent + contracts + locales + external clients)
gen-all: gen-ent gen gen-locales gen-amocrm

## tidy: prune/sync go.mod after generation or dep changes
tidy:
	go mod tidy

# ---------------------------------------------------------------------------
# Lint & typecheck
# ---------------------------------------------------------------------------

## lint: lint Go + frontend
lint: lint-go lint-web

## lint-go: locale catalogs, gofmt, go vet, golangci-lint (skips generated files)
lint-go: lint-locales
	@test -z "$$(gofmt -l cmd internal ent/schema)" || (echo "gofmt needed:"; gofmt -l cmd internal ent/schema; exit 1)
	go vet ./...
	golangci-lint run

## lint-locales: verify generated source messages and complete translations
lint-locales:
	sh scripts/lint_locales.sh

## lint-web: oxlint, oxfmt check, and TypeScript project build
lint-web:
	pnpm lint
	pnpm format:check
	pnpm check

## lint-fix: auto-fix Go and frontend formatting and lint findings
lint-fix:
	gofmt -w cmd internal ent/schema
	-golangci-lint run --fix
	pnpm lint:fix
	pnpm format

# ---------------------------------------------------------------------------
# Test
# ---------------------------------------------------------------------------

# Development DB coordinates.
DB_HOST ?= 127.0.0.1
DB_PORT ?= $(DATABASE_PORT)
DEV_DB_NAME ?= code_basics_development
DEV_DB_URL  ?= postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DEV_DB_NAME)?sslmode=disable
DB_ADMIN_URL ?= postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/postgres?sslmode=disable

## dev-db-create: idempotently create the development database.
dev-db-create:
	@if ! psql "$(DB_ADMIN_URL)" -tAc "SELECT 1 FROM pg_database WHERE datname = '$(DEV_DB_NAME)'" | grep -q 1; then \
		createdb --maintenance-db="$(DB_ADMIN_URL)" "$(DEV_DB_NAME)"; \
	fi

## migrate-new: scaffold a new empty migration to hand-author, e.g.
## `make migrate-new NAME=add_widgets`. Edit it, then `atlas migrate hash`.
migrate-new:
	atlas migrate new $(NAME) --dir "file://migrations"
	atlas migrate hash --dir "file://migrations"

## dev-migrate: apply atlas migrations to the development DB.
dev-migrate: dev-db-create
	atlas migrate apply --env local --url "$(DEV_DB_URL)&search_path=public"

## dev-prepare: create and migrate the development DB without replacing its data.
dev-prepare: dev-db-create dev-migrate

## test: run the Go suite; DB packages start disposable Postgres 18 containers
test:
	go test ./...

## test-web: run the frontend test suite (Vitest Browser Mode via Playwright).
## Needs a browser — run `pnpm exec playwright install chromium` once locally.
test-web:
	pnpm test

# ---------------------------------------------------------------------------
# Build
# ---------------------------------------------------------------------------

## build: build the API, async worker, and frontend bundle
build: build-api build-worker build-web

## build-api: compile the Go server to bin/server
build-api:
	go build -o bin/server ./cmd/server

## build-worker: compile the async runtime to bin/worker
build-worker:
	go build -o bin/worker ./cmd/worker

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

.PHONY: help prepare install setup services-start services-stop dev-db-create migrate-new dev-migrate dev-prepare dev dev-api dev-worker dev-web dev-spec \
	gen gen-spec gen-api gen-client gen-amocrm gen-ent gen-all tidy \
	lint lint-go lint-web lint-fix test build build-api build-worker build-web clean \
	deps-update update-skills
