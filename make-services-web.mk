web-exercises-load-php:
	bundle exec rake "exercises:load[${L}]"

# caddy-docker-build-production:
# 	docker build --file services/caddy/Dockerfile.production --tag hexletbasics/services-caddy:$(VERSION) services/caddy

# caddy-docker-push:
# 	docker push hexletbasics/services-caddy:$(VERSION)
