include k8s/Makefile

setup:
	# brew install vips
	cp -n .env.sample .env || exit 0
	bin/setup --skip-server
	bin/rails db:fixtures:load
	pnpm install
	pnpm simple-git-hooks

test-all: test test-frontend test-system

test:
	bin/rails test

test-frontend:
	pnpm vitest run

test-frontend-watch:
	pnpm vitest

test-system:
	bin/rails test:system

db-reset:
	bin/rails db:reset
	bin/rails db:fixtures:load

dev:
	# rm -rf dist
	bin/dev

log-mails:
	 tail -n 30 log/mailer.log | base64 --decode

build-assets:
	bin/rails assets:precompile

analyze-bundle:
	bin/vite build
	open stats.html

setup-staging:
	bin/vite clobber
	# VISUAL="code --wait" bin/rails credentials:edit
	RAILS_ENV=staging bin/rails db:drop
	RAILS_ENV=staging bin/rails db:prepare
	RAILS_ENV=staging bin/rails db:fixtures:load
	# RAILS_ENV=staging bin/rails db:prepare db:fixtures:load
	# bin/vite build --ssr
	RAILS_ENV=staging make build-assets

staging:
	DEBUG=vite-plugin-ruby:* NODE_ENV=development RAILS_ENV=staging overmind start -f Procfile.staging

dev-ssr:
	bin/vite ssr

sync-i18n:
	bundle exec i18n export

sync-fixtures:
	bin/rails db:fixtures:load

sync-browserlist:
	bundle exec browserslist generate

editor-setup:
	-bin/tapioca gems --verify
	-bin/tapioca dsl --verify
	-bin/tapioca require
	-bin/tapioca gem

sync-types:
	ENABLE_TYPELIZER=1 bin/rails typelizer:generate:refresh
	bin/rails js:routes:typescript
	bin/rails app:export_events_to_ts
	bin/rails app:export_enums_to_ts


sync: sync-i18n sync-fixtures sync-types sync-browserlist

coverage-open:
	open coverage/index.html

app-lint-staged:
	echo 'disabled'

language-load:
	bin/rails app:load_exercises[${L}]

check-types:
	# bundle exec srb tc
	pnpm run check

lint:
	pnpm exec biome check
	pnpm tsc --build
	bin/rubocop

lint-fix:
	bin/rubocop -x
	pnpm exec biome check --fix

clear:
	rm -rf ./.overmind.sock
	bin/rails log:clear
	bin/rails tmp:clear
	bin/rails assets:clobber

docker-build:
	docker build . -t hexlet-basics/hexlet-basics

docker-staging: docker-build
	# TODO: implement starting server

tag:
	git tag $(TAG) && git push --tags --no-verify

next-tag:
	@old_version=$(shell git fetch --tags && git tag -l 'v[0-9]*' | sort -V | tail -n 1 | sed 's/^v//'); \
		new_version=$$((old_version + 1)); \
		make tag TAG=v$$new_version

services-frontend-run:
	bin/vite dev

services-frontend-ssr-run:
	bin/vite ssr

services-app-run:
	bin/rails s -p 3000

services-jobs-run:
	bin/jobs

services-webserver-run:
	caddy run # --config ./services/webserver/caddy/conf/Caddyfile --envfile=.env

services-cable-run:
	bundle exec puma -p 28080 cable/config.ru

services-db-start:
	docker run -d -it --rm \
		-p 5432:5432 \
		--name code_basics_postgres \
		-e POSTGRES_DB=code_basics_development \
		-e POSTGRES_PASSWORD=postgres \
		-v code_basics_pgdata:/var/lib/postgresql/data \
		postgres:17

services-db-stop:
	docker stop code_basics_postgres

services-remove: services-db-remove

services-db-remove: services-db-stop
	docker rm code_basics_postgres || true
	docker volume remove code_basics_pgdata

services-db-down: service-db-remove
	docker volume remove code_basics_pgdata

services-stop: services-db-stop

services-start: services-db-start

# sync-analytics:
# 	bin/rails analytics:refresh_user_survey_pivot

setup-macos:
	brew install caddy libpq vips watchman
	brew link --force libpq

setup-ubuntu:
	curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
	curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
	sudo apt-get update && sudo apt-get install -yq caddy libpq-dev libvips
	gem install overmind
	# run caddy on port 443 without sudo
	# sudo setcap CAP_NET_BIND_SERVICE=+ep $(which caddy)

.PHONY: test
