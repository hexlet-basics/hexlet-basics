compose-build:
	docker-compose build

compose:
	docker-compose up -d

compose-down:
	docker-compose down -v || true

compose-restart:
	docker-compose restart

web-install:
	docker-compose run web bundle install
	docker-compose run web yarn install --check-files

web-ci-test:
	cp -n .env.example .env || true
	make compose-setup
	make app-test

compose-setup: compose-down compose-build web-install

app-bash:
	docker-compose run web /bin/bash

app-test:
	docker-compose run web make test-db-prepare
	docker-compose run web /bin/bash -c "make lint && echo fs.inotify.max_user_watches=524288 | tee -a /etc/sysctl.conf && sysctl -p
 &&  make test"

app-rails-console:
	docker-compose run web bin/rails c

