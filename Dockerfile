# TODO Switch to ruby 3.1 Rails 7.0
FROM ruby:3.0.3

# NOTE https://github.com/webpack/webpack/issues/14532
ENV NODE_OPTIONS --openssl-legacy-provider
ENV NODE_VERSION 17.x

ENV DOCKER_CHANNEL stable
ENV DOCKER_VERSION 20.10.11

RUN curl -sL https://deb.nodesource.com/setup_${NODE_VERSION} | bash -

RUN curl -fsSL "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz" \
| tar -xzC /usr/local/bin --strip=1 docker/docker

RUN apt-get update && apt-get install -y     \
  build-essential \
  libpq-dev \
  libsqlite3-dev \
  nodejs \
  && rm -rf /var/lib/apt/lists/*

RUN npm install -g \
  npm-check-updates \
  yarn

ENV BUNDLE_PATH /root/hexlet-basics/vendor/bundle
WORKDIR /root/hexlet-basics
