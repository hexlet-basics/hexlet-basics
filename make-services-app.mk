app-bash:
	docker-compose run web /bin/bash

app-test:
	docker-compose run web make test-db-prepare
	docker-compose run web make lint
	docker-compose run web make test

app-rails-console:
	docker-compose run web bin/rails c

app-setup: app-install app-db-prepare

app-install:
	docker-compose run web bundle install
	docker-compose run web yarn install --check-files

app-db-prepare:
	docker-compose run web bin/rails db:drop || true
	docker-compose run web bin/rails db:create || true
	docker-compose run web bin/rails db:schema:load || true
	docker-compose run web bin/rails db:migrate || true
	docker-compose run web bin/rails db:migrate || true
	docker-compose run web bin/rails db:fixtures:load || true
