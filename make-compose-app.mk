app-bash:
	docker-compose run --rm web /bin/bash

app-language-load:
	docker-compose run --rm web make language-load L=${L}

app-lint:
	docker-compose run --rm --no-deps web make lint

app-lint-fix:
	docker-compose run --rm --no-deps web make lint-fix

app-test:
	docker-compose run --rm web make test

app-check: app-test app-lint

app-test-file:
	docker-compose run --rm web make test-file T=${T}

app-rails-console:
	docker-compose run --rm web bin/rails c

app-setup-git-hooks:
	docker-compose run --rm web npx simple-git-hooks

app-lint-staged:
	docker-compose run -T --rm web npx lint-staged --relative

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
	make app-language-load L='kotlin'
	make app-language-load L='swift'
	make app-language-load L='rust'
	make app-language-load L='perl'
	make app-language-load L='ocaml'
	make app-language-load L='dart'
	make app-language-load L='crystal'
	make app-language-load L='powershell'
	make app-language-load L='dlang'


app-setup: app-install app-db-prepare app-setup-git-hooks app-languages-load

app-install:
	docker-compose run --rm web make setup

app-db-prepare:
	docker-compose run --rm web make db-prepare
