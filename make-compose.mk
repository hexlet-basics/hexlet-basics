compose:
	docker-compose up

compose-build:
	docker-compose build

compose-down:
	docker-compose down -v || true

compose-restart:
	docker-compose restart

compose-setup: compose-down compose-build app-setup

compose-ci-build:
	docker-compose -f docker-compose.yml build

compose-ci:
	docker-compose -f docker-compose.yml up


# env-prepare:
# 	cp -n .env.example .env || true
