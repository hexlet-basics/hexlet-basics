compose-build:
	docker-compose build

compose:
	docker-compose up -d

compose-down:
	docker-compose down -v || true

compose-restart:
	docker-compose restart

web-ci-test:
	cp -n .env.example .env || true
	make compose-setup
	make app-test

compose-setup: compose-down compose-build app-setup
