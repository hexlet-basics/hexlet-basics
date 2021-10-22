app-bash:
	docker-compose run web /bin/bash

app-language-load:
	docker-compose run web make language-load L=${L}

app-lint:
	docker-compose run web make lint

app-test:
	docker-compose run web make test

app-check: app-test app-lint

app-test-file:
	docker-compose run web make test ${T}

app-rails-console:
	docker-compose run web bin/rails c

app-setup-git-hooks:
	docker-compose run web yarn run simple-git-hooks

app-languages-load:
	make app-language-load L='javascript'
	make app-language-load L='php'
	# make app-language-load L='python'
	make app-language-load L='css'
	# make app-language-load L='html'
	# make app-language-load L='ruby'
	make app-language-load L='go'
	make app-langugage-load L='clojure'


app-setup: app-install app-db-prepare app-setup-git-hooks app-languages-load

app-install:
	docker-compose run web make setup

app-db-prepare:
	docker-compose run web make db-prepare
