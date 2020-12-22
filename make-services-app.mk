app-bash:
	docker-compose run app /bin/bash

app-language-load:
	docker-compose run app make language-load L=${L}

app-test: app-db-prepare
	docker-compose run app make test
	docker-compose run app make lint

app-rails-console:
	docker-compose run app bin/rails c

app-setup: env-prepare app-install app-db-prepare
	make app-language-load L='javascript'
	make app-language-load L='php'
	make app-language-load L='python'
	make app-language-load L='ruby'
	make app-language-load L='css'

app-install:
	docker-compose run app make setup

app-db-prepare:
	docker-compose run app make db-prepare
