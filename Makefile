test:
	bin/rails test

setup:
	bin/setup
	bin/rails db:fixtures:load

fixtures-load:
	bin/rails db:fixtures:load

clean:
	bin/rails db:drop

db-reset:
	bin/rails db:drop
	bin/rails db:create
	bin/rails db:migrate
	bin/rails db:fixtures:load

start:
	bundle exec heroku local

precompile-assets:
	bundle exec rails assets:precompile


lint: lint-eslint lint-rubocop

linter-fix:
	bundle exec rubocop --auto-correct

lint-rubocop:
	bundle exec rubocop

lint-eslint:
	npx eslint app/javascript

deploy:
	git push heroku master

heroku-console:
	heroku run rails console

heroku-logs:
	heroku logs --tail

ci-test:
	make setup
	make precompile-assets
	make lint
	make test

.PHONY: test
