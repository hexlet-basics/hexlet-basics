compose-build:
	docker-compose build

compose:
	docker-compose up -d

compose-down:
	docker-compose down -v || true

compose-restart:
	docker-compose restart

web-ci-test:
	make compose-setup
	make app-test

compose-setup: env-prepare compose-down compose-build app-setup

env-prepare:
	cp -n .env.example .env || true
