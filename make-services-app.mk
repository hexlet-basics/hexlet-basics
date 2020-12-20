app-bash:
	docker-compose run app /bin/bash

app-language-load:
	docker-compose run app make language-load L=${L}

app-test: app-db-prepare
	docker-compose run app make test
	docker-compose run app make lint

app-rails-console:
	docker-compose run app bin/rails c

app-setup: app-install app-db-prepare

app-install:
	docker-compose run app make setup

app-db-prepare:
	docker-compose run app make db-prepare
