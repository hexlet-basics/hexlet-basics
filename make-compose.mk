compose:
	docker-compose up -d

compose-build:
	docker-compose build

compose-logs:
	docker-compose logs -f

compose-down:
	docker-compose down || true

compose-stop:
	docker-compose stop || true

compose-restart:
	docker-compose restart

compose-setup: compose-down compose-build app-setup

compose-ci-build:
	docker-compose -f docker-compose.yml build

compose-ci:
	docker-compose --file docker-compose.yml up --abort-on-container-exit

env-prepare:
	cp -n .env.populate .env || true
