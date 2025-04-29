include k8s/Makefile

setup:
	# brew install vips
	cp -n .env.sample .env || exit 0
	bin/setup --skip-server
	bin/rails db:fixtures:load
	npm install
	npx simple-git-hooks

test:
	bin/rails test

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

setup-staging:
	bin/vite clobber
	# VISUAL="code --wait" bin/rails credentials:edit
	RAILS_ENV=staging bin/rails db:drop
	RAILS_ENV=staging bin/rails db:prepare
	RAILS_ENV=staging bin/rails db:fixtures:load
	# RAILS_ENV=staging bin/rails db:prepare db:fixtures:load
	# bin/vite build --ssr
	NODE_ENV=development RAILS_ENV=staging bin/rails assets:precompile

staging:
	DEBUG=vite-plugin-ruby:* NODE_ENV=development RAILS_ENV=staging overmind start -f Procfile.staging

dev-ssr:
	bin/vite ssr

i18n-export:
	bundle exec i18n export

sync-fixtures:
	bin/rails db:fixtures:load

editor-setup:
	-bin/tapioca gems --verify
	-bin/tapioca dsl --verify
	-bin/tapioca require
	-bin/tapioca gem

sync: i18n-export sync-fixtures
	ENABLE_TYPELIZER=1 bin/rails typelizer:generate:refresh
	bin/rails js:routes:typescript

coverage-open:
	open coverage/index.html

app-lint-staged:
	echo 'disabled'

language-load:
	bin/rails exercises:load[${L}]

check-types:
	# bundle exec srb tc
	npm run check

lint:
	npx tsc --build
	bin/rubocop
	npx oxlint

lint-fix:
	bin/rubocop -x
	npx oxlint --fix --fix-suggestions

docker-build:
	docker build . -t hexlet-basics/hexlet-basics

docker-staging: docker-build
	# TODO: implement starting server

ansible-generate-env:
	docker run --rm -e RUNNER_PLAYBOOK=ansible/development.yml \
		-v $(CURDIR)/ansible/development:/runner/inventory \
		-v $(CURDIR):/runner/project \
		quay.io/ansible/ansible-runner

ansible-terraform-vars-generate:
	docker run --rm -e RUNNER_PLAYBOOK=ansible/terraform.yml \
		-v $(CURDIR)/ansible/production:/runner/inventory \
		-v $(CURDIR):/runner/project \
		-e ANSIBLE_VAULT_PASSWORD_FILE=tmp/ansible-vault-password \
		quay.io/ansible/ansible-runner

ansible-vaults-edit:
	docker run -it --rm \
		-v $(CURDIR):/runner/project \
		quay.io/ansible/ansible-runner ansible-vault edit --vault-password-file project/tmp/ansible-vault-password project/ansible/production/group_vars/all/vault.yml

tag:
	git tag $(TAG) && git push --tags --no-verify

next-tag:
	@old_version=$(shell git fetch --tags && git tag -l 'v[0-9]*' | sort -V | tail -n 1 | sed 's/^v//'); \
		new_version=$$((old_version + 1)); \
		make tag TAG=v$$new_version

services-frontend-run:
	bin/vite dev

services-app-run:
	bin/rails s -p 3000

services-jobs-run:
	bin/jobs

services-webserver-run:
	caddy run # --config ./services/webserver/caddy/conf/Caddyfile --envfile=.env

services-db-start:
	docker run -d -it --rm \
		-p 5432:5432 \
		--name code_basics_postgres \
		-e POSTGRES_DB=code_basics_development \
		-e POSTGRES_PASSWORD=postgres \
		-v code_basics_pgdata:/var/lib/postgresql/data \
		postgres

services-db-stop:
	docker stop code_basics_postgres

services-db-remove: services-db-stop
	docker volume remove codebasics_pgdata

services-stop: services-db-stop

services-start: services-db-start

sync-analytics:
	bin/rails analytics:refresh_user_survey_pivot

.PHONY: test
