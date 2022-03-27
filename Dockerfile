FROM ruby:3.1

# NOTE https://github.com/webpack/webpack/issues/14532
ENV NODE_OPTIONS --openssl-legacy-provider
ENV NODE_VERSION 17.x

ENV DOCKER_CHANNEL stable
ENV DOCKER_VERSION 20.10.14

RUN curl -sL https://deb.nodesource.com/setup_${NODE_VERSION} | bash -

RUN curl -fsSL "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz" \
| tar -xzC /usr/local/bin --strip=1 docker/docker

RUN apt-get update && apt-get install -y \
  build-essential \
  bash-completion \
  libpq-dev \
  libsqlite3-dev \
  nodejs \
  && rm -rf /var/lib/apt/lists/*

RUN corepack enable
RUN yarn set version stable
# RUN npm install -g yarn

# ENV BUNDLE_PATH /root/hexlet-basics/vendor/bundle
ENV PROJECT_ROOT /app
WORKDIR ${PROJECT_ROOT}

ENV BUNDLE_APP_CONFIG ${PROJECT_ROOT}/.bundle/config
ENV GEM_HOME ${PROJECT_ROOT}/vendor/bundle
ENV BUNDLE_PATH ${GEM_HOME}

# RUN bundle config build.nokogiri --use-system-libraries
# BUNDLE_BUILD__NOKOGIRI: "--use-system-libraries"
