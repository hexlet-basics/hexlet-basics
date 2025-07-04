# BUNDLE_BUILD__NOKOGIRI: "--use-system-libraries"
# syntax=docker/dockerfile:1
# check=error=true

# This Dockerfile is designed for production, not development. Use with Kamal or build'n'run by hand:
# docker build -t store .
# docker run -d -p 80:80 -e RAILS_MASTER_KEY=<value from config/master.key> --name store store

# For a containerized dev environment, see Dev Containers: https://guides.rubyonrails.org/getting_started_with_devcontainer.html

# Make sure RUBY_VERSION matches the Ruby version in .ruby-version
ARG RUBY_VERSION=3.4.4
FROM docker.io/library/ruby:$RUBY_VERSION-slim AS base

# Rails app lives here
WORKDIR /rails

# Install base packages
RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y curl libjemalloc2 libvips sqlite3 libpq-dev libyaml-dev make && \
    rm -rf /var/lib/apt/lists /var/cache/apt/archives

ENV LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2

# Set production environment
ENV RAILS_ENV="production" \
    BUNDLE_DEPLOYMENT="1" \
    BUNDLE_PATH="/usr/local/bundle" \
    BUNDLE_WITHOUT="development"

ARG DOCKER_CHANNEL=stable
ARG DOCKER_VERSION=27.3.1
# NOTE: Жестко указан gid группы docker на нодах. Надо подумать как переделать лучше
RUN groupadd --gid 998 docker
RUN curl -fsSL "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz" \
  | tar -xzC /usr/local/bin --strip=1 docker/docker

# Throw-away build stage to reduce size of final image
FROM base AS build

RUN curl -sL https://deb.nodesource.com/setup_24.x | sh -

# Install packages needed to build gems
RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y build-essential git pkg-config nodejs && \
    rm -rf /var/lib/apt/lists /var/cache/apt/archives

# Install application gems
COPY .ruby-version Gemfile Gemfile.lock ./
RUN bundle install && \
    rm -rf ~/.bundle/ "${BUNDLE_PATH}"/ruby/*/cache "${BUNDLE_PATH}"/ruby/*/bundler/gems/*/.git && \
    bundle exec bootsnap precompile --gemfile

COPY package.json package-lock.json ./
RUN npm ci

# Copy application code
COPY . .

# Precompile bootsnap code for faster boot times
RUN bundle exec bootsnap precompile app/ lib/

ARG RELEASE_VERSION="unknown"
ARG APP_HOST
ENV VITE_RELEASE_VERSION=$RELEASE_VERSION
ENV VITE_APP_HOST=$APP_HOST
ENV VITE_NODE_ENV=production

# Precompiling assets for production without requiring secret RAILS_MASTER_KEY
RUN --mount=type=secret,id=sentry-org,env=VITE_SENTRY_ORG \
    --mount=type=secret,id=sentry-url,env=VITE_SENTRY_URL \
    --mount=type=secret,id=sentry-project,env=VITE_SENTRY_PROJECT \
    --mount=type=secret,id=sentry-dsn,env=VITE_SENTRY_DSN \
    --mount=type=secret,id=sentry-auth-token,env=VITE_SENTRY_AUTH_TOKEN \
    # --mount=type=secret,id=carrotquest-api-key,env=VITE_CARROTQUEST_API_KEY \
    --mount=type=secret,id=yandex-metrika-counter-id,env=VITE_YANDEX_METRIKA_COUNTER_ID \
    --mount=type=secret,id=n8n_base_url,env=N8N_BASE_URL \
    SECRET_KEY_BASE_DUMMY=1 make build-assets

# NOTE: .git is included for Sentry, then removed to keep the image clean.
RUN rm -rf .git

# Final stage for app image
FROM base

# Copy built artifacts: gems, application
COPY --from=build "${BUNDLE_PATH}" "${BUNDLE_PATH}"
COPY --from=build /rails /rails
COPY --from=build /usr/bin/node /usr/bin/node

# Run and own only the runtime files as a non-root user for security
RUN groupadd --system --gid 1000 rails && \
    useradd rails --uid 1000 --gid 1000 --create-home --shell /bin/bash && \
    chown -R rails:rails db log storage tmp && \
    usermod -aG docker rails

# NOTE: используя дефолт 1000:1000 не обновляются групповые права, в результате docker недоступен
USER rails

# Entrypoint prepares the database.
ENTRYPOINT ["/rails/bin/docker-entrypoint"]

# Start server via Thruster by default, this can be overwritten at runtime
EXPOSE 80

CMD ["./bin/thrust", "./bin/rails", "server"]
