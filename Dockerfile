FROM ruby:3.0.3

RUN apt-get update \
      && apt-get install -y build-essential libpq-dev libsqlite3-dev

RUN curl -sL https://deb.nodesource.com/setup_17.x | bash -

# NOTE https://github.com/webpack/webpack/issues/14532
ENV NODE_OPTIONS --openssl-legacy-provider

ENV DOCKER_CHANNEL stable
ENV DOCKER_VERSION 20.10.11
RUN curl -fsSL "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz" \
  | tar -xzC /usr/local/bin --strip=1 docker/docker

RUN apt-get update && apt-get install -y nodejs
RUN npm install -g npm-check-updates
RUN npm install -g yarn

# RUN yarn set version berry
# RUN yarn config set --home enableTelemetry 0

ENV BUNDLE_PATH /root/hexlet-basics/vendor/bundle
WORKDIR /root/hexlet-basics
