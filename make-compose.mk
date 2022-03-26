compose:
	docker compose up -d

compose-production:
	docker compose --file docker-compose.yml run production

compose-build:
	docker compose build

compose-logs:
	docker compose logs -f

compose-down:
	docker compose down || true

compose-clear:
	docker compose down -v || true

compose-stop:
	docker compose stop || true

compose-restart:
	docker compose restart

compose-setup: compose-down compose-build app-setup

compose-ci-build:
	docker build -f Dockerfile.production -t hexletbasics/services-app:cached .
	docker compose -f docker-compose.yml build

compose-ci: compose-ci-build
	docker compose --file docker-compose.yml up --abort-on-container-exit

env-prepare:
	cp -n .env.populate .env || true
