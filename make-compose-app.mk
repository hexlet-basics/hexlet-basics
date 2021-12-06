app-bash:
	docker-compose run --rm web /bin/bash

app-language-load:
	docker-compose run --rm web make language-load L=${L}

app-lint:
	docker-compose run --rm --no-deps web make lint

app-test:
	docker-compose run --rm web make test

app-check: app-test app-lint

app-test-file:
	docker-compose run --rm web make test ${T}

app-rails-console:
	docker-compose run --rm web bin/rails c

app-setup-git-hooks:
	docker-compose run --rm web yarn run simple-git-hooks

app-languages-load:
	make app-language-load L='javascript'
	make app-language-load L='php'
	make app-language-load L='python'
	make app-language-load L='css'
	make app-language-load L='html'
	make app-language-load L='ruby'
	make app-language-load L='go'
	make app-language-load L='clojure'
	make app-language-load L='racket'
	make app-language-load L='java'
	make app-language-load L='elixir'
	make app-language-load L='clang'
	make app-language-load L='lua'
	make app-language-load L='csharp'
	make app-language-load L='typescript'
	make app-language-load L='prolog'
	make app-language-load L='haskell'
	make app-language-load L='cpp'
	make app-language-load L='bash'
	make app-language-load L='fortran'

app-setup: app-install app-db-prepare app-setup-git-hooks app-languages-load

app-install:
	docker-compose run --rm web make setup

app-db-prepare:
	docker-compose run --rm web make db-prepare
