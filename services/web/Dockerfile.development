FROM ruby:2.7.1

RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list

RUN apt-get update -qq \
  && apt-get install -y --no-install-recommends \
  build-essential \
  libpq-dev \
  libsqlite3-dev  \
  nodejs \
  yarn

ENV DOCKER_CHANNEL edge
ENV DOCKER_VERSION 19.03.8
RUN curl -fsSL "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz" \
  | tar -xzC /usr/local/bin --strip=1 docker/docker

ENV BUNDLE_PATH /app/vendor/bundle

WORKDIR /app
