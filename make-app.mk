test:
	bin/rails test -d
	RAILS_LOCALE=ru bin/rails test -d

check: test lint vulnerability-check

vulnerability-check:
	bundle exec brakeman

ci-check: lint
	bin/rails db:create
	#ls -la public/packs-test
	DISABLE_SPRING=1 bin/rails test

language-load:
	bin/rails exercises:load[${L}]

fixtures-load:
	bin/rails db:fixtures:load

frontend: js-routes
	npx nodemon -L --watch webpack.config.js --exec npm run build:watch

setup:
	bin/setup
	yarn install

db-reset:
	bin/rails db:drop
	bin/rails db:create
	bin/rails db:migrate

start:
	rm -rf tmp/pids/server.pid
	bundle exec rails s -p 3000 -b '0.0.0.0'

start-production:
	bin/rails db:migrate
	bin/rails server -e production

clean:
	bin/rails db:drop

precompile-assets:
	bundle exec rails assets:precompile

lint: lint-eslint lint-rubocop lint-slim

lint-fix: lint-eslint-fix lint-rubocop-fix

lint-style:
	npx stylelint "app/assets/stylesheets/*.scss"

lint-style-fix:
	npx stylelint --fix "app/assets/stylesheets/*.scss"

lint-rubocop:
	bundle exec rubocop

lint-rubocop-fix:
	bundle exec rubocop -A

lint-eslint:
	npx eslint app/javascript --ext .js,.jsx

lint-eslint-fix:
	npx eslint app/javascript --ext .js,.jsx --fix

lint-slim:
	bundle exec slim-lint app/views/

js-routes:
	bundle exec rails js:routes

db-prepare:
	bin/rails db:drop || true
	bin/rails db:create || true
	bin/rails db:schema:load || true
	bin/rails db:migrate || true
	bin/rails db:fixtures:load || true

webpacker:
	bin/webpacker-dev-server

.PHONY: test
