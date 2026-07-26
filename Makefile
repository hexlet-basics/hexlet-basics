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

## setup: install all deps (go, frontend, api-spec)
setup:
	go mod download
	pnpm install
	cd api-spec && pnpm install

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
	cd api-spec && pnpm watch

# ---------------------------------------------------------------------------
# Code generation (contract is the source of truth — see ADR-0001)
# ---------------------------------------------------------------------------

## gen: full contract pipeline — TypeSpec -> OpenAPI -> ogen server + hey-api client
gen: gen-spec gen-api gen-client

## gen-spec: compile TypeSpec to api-spec/dist/openapi.yaml
gen-spec:
	cd api-spec && pnpm build

## gen-api: generate the ogen Go server from OpenAPI
gen-api:
	go generate ./internal/apigen

## gen-client: generate the hey-api TS client + Query hooks from OpenAPI
gen-client:
	pnpm generate

## gen-ent: regenerate the ent ORM from ent/schema
gen-ent:
	go generate ./ent

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

## deps-update: bump Go + frontend + api-spec dependencies to latest
deps-update:
	go get -u ./...
	go mod tidy
	npx --yes npm-check-updates -u
	pnpm install
	cd api-spec && npx --yes npm-check-updates -u && pnpm install
	@echo ">> deps bumped — run 'make gen-all' then 'make lint test' to verify"

## update-skills: update the project's agent skills
update-skills:
	npx --yes skills update --project --yes

.PHONY: help setup services-start services-stop dev dev-api dev-web dev-spec \
	gen gen-spec gen-api gen-client gen-ent gen-all tidy \
	lint lint-go lint-web lint-fix test build build-api build-web clean \
	deps-update update-skills
