app-bash:
	docker-compose run web /bin/bash

app-test: app-db-prepare
	docker-compose run web make test
	docker-compose run web make lint

app-rails-console:
	docker-compose run web bin/rails c

app-setup: app-install app-db-prepare app-load-php app-load-javascript app-load-css app-load-html app-load-elixir app-load-python

app-install:
	docker-compose run web bundle install
	docker-compose run web yarn install --check-files

app-db-prepare:
	docker-compose run web make db-prepare

app-load-php:
	docker-compose run web bin/rails "exercises:load[php]"

app-load-javascript:
	docker-compose run web bin/rails "exercises:load[javascript]"

app-load-css:
	docker-compose run web bin/rails "exercises:load[css]"

app-load-html:
	docker-compose run web bin/rails "exercises:load[html]"

app-load-elixir:
	docker-compose run web bin/rails "exercises:load[elixir]"

app-load-python:
	docker-compose run web bin/rails "exercises:load[python]"
