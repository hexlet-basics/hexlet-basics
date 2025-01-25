include k8s/Makefile

setup:
	# brew install vips
	bundle install
	bin/rails db:prepare
	bin/rails db:fixtures:load

test:
	bin/rails test

db-reset:
	bin/rails db:reset
	bin/rails db:fixtures:load

dev:
	bin/dev

log-mails:
	 tail -n 30 log/mailer.log | base64 --decode

staging:
	bin/vite clobber
	# VISUAL="code --wait" bin/rails credentials:edit
	RAILS_ENV=staging bin/rails db:prepare db:fixtures:load
	# bin/vite build --ssr
	NODE_ENV=development RAILS_ENV=staging bin/rails assets:precompile
	DEBUG=vite-plugin-ruby:* NODE_ENV=development RAILS_ENV=staging bundle exec foreman start -f Procfile.staging

dev-ssr:
	bin/vite ssr

i18n-export:
	bundle exec i18n export

sync-fixtures:
	bin/rails db:fixtures:load

editor-setup:
	bin/tapioca dsl

sync: i18n-export sync-fixtures
	bin/rails typelizer:generate:refresh

coverage-open:
	open coverage/index.html

app-lint-staged:
	echo 'disabled'

language-load:
	bin/rails exercises:load[${L}]

lint:
	bin/rubocop -x
	# npx @biomejs/biome check

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

# tag:
# 	git tag $(TAG) && git push --tags --no-verify
#
# next-tag:
# 	make tag TAG=$(shell bin/generate_next_tag)

.PHONY: test
