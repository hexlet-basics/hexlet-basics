compose:
	docker-compose up -d

compose-build:
	docker-compose build

compose-down:
	docker-compose down || true

compose-restart:
	docker-compose restart

compose-setup: compose-down compose-build app-setup

compose-ci-build:
	docker-compose -f docker-compose.yml build

compose-ci:
	docker-compose -f docker-compose.yml up

env-prepare:
	cp -n .env.populate .env || true
