compose-build:
	docker-compose build

compose:
	docker-compose up -d

compose-down:
	docker-compose down -v || true

compose-restart:
	docker-compose restart

compose-console:
	docker-compose run web /bin/bash

web-install:
	docker-compose run web bundle install
	docker-compose run web yarn install --check-files

ci-test:
	cp -n .env.example .env || true
	make compose-setup
	docker-compose run web make lint
	docker-compose run web make test

compose-setup: compose-down compose-build web-install
