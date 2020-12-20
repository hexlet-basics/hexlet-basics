app-bash:
	docker-compose run web /bin/bash

app-language-load:
	docker-compose run web make language-load L=${L}

app-test: app-db-prepare
	docker-compose run web make test
	docker-compose run web make lint

app-rails-console:
	docker-compose run web bin/rails c

app-setup: app-install app-db-prepare

app-install:
	docker-compose run web bundle install
	docker-compose run web yarn install --check-files

app-db-prepare:
	docker-compose run web make db-prepare
