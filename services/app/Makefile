test:
	RAILS_ENV=test bin/rails test

fixtures-load:
	bin/rails db:fixtures:load

setup:
	cp -n .env.example .env || true
	bin/setup
	yarn install

db-reset:
	bin/rails db:drop
	bin/rails db:create
	bin/rails db:migrate

start:
	bundle exec rails s -p 3000 -b '0.0.0.0'

start-production:
	bin/rails db:migrate
	bin/rails server -e production

clean:
	bin/rails db:drop

precompile-assets:
	bundle exec rails assets:precompile

lint: lint-eslint lint-rubocop

linter-fix:
	bundle exec rubocop --auto-correct

lint-rubocop:
	bundle exec rubocop

lint-eslint:
	npx eslint app/javascript

js-routes:
	bundle exec rails js_routes:generate

db-prepare:
	bin/rails db:drop || true
	bin/rails db:create || true
	bin/rails db:schema:load || true
	bin/rails db:migrate || true
	bin/rails db:fixtures:load || true

.PHONY: test
