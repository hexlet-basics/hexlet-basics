web-exercises-load-php:
	bundle exec rake "exercises:load[php]"

web-exercises-load-python:
	bundle exec rake "exercises:load[python]"

web-exercises-load-ruby:
	bundle exec rake "exercises:load[ruby]"

web-exercises-load-racket:
	bundle exec rake "exercises:load[racket]"

web-exercises-load-javascript:
	bundle exec rake "exercises:load[javascript]"

web-exercises-load-java:
	bundle exec rake "exercises:load[java]"

web-exercises-load-html:
	bundle exec rake "exercises:load[html]"

web-exercises-load-css:
	bundle exec rake "exercises:load[css]"

web-exercises-load-elixir:
	bundle exec rake "exercises:load[elixir]"

web-exercises-load-go:
	bundle exec rake "exercises:load[go]"

caddy-docker-build-production:
	docker build --file services/caddy/Dockerfile.production --tag hexletbasics/services-caddy:$(VERSION) services/caddy

caddy-docker-push:
	docker push hexletbasics/services-caddy:$(VERSION)
