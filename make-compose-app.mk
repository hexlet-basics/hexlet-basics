app-bash:
	docker-compose run web /bin/bash

app-language-load:
	docker-compose run web make language-load L=${L}

app-test:
	docker-compose run web make test
	docker-compose run web make lint

app-test-file:
	docker-compose run web make test ${T}

app-rails-console:
	docker-compose run web bin/rails c

app-setup: app-install app-db-prepare
	make app-language-load L='javascript'
	make app-language-load L='php'
	make app-language-load L='python'
	make app-language-load L='css'

app-install:
	docker-compose run web make setup

app-db-prepare:
	docker-compose run web make db-prepare
