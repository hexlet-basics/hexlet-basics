FROM ruby:3.1.3

# RUN bundle config --global frozen 1

ENV NODE_VERSION 19.x

RUN curl -sL https://deb.nodesource.com/setup_${NODE_VERSION} | bash -

RUN apt-get install -y nodejs

ENV DOCKER_CHANNEL stable
ENV DOCKER_VERSION 20.10.21

RUN curl -fsSL "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz" \
| tar -xzC /usr/local/bin --strip=1 docker/docker

RUN apt-get update && apt-get install -y \
  build-essential \
  bash-completion \
  libpq-dev \
  libsqlite3-dev \
  libvips42 \
  && rm -rf /var/lib/apt/lists/*

# ENV BUNDLE_PATH /root/hexlet-basics/vendor/bundle
ENV PROJECT_ROOT /opt/projects/hexlet-basics
RUN mkdir -p ${PROJECT_ROOT}

# NOTE: for initial make project-setup
RUN mkdir -p ${PROJECT_ROOT}/public

WORKDIR ${PROJECT_ROOT}

ENV BUNDLE_APP_CONFIG ${PROJECT_ROOT}/.bundle/config
ENV GEM_HOME ${PROJECT_ROOT}/vendor/bundle
ENV BUNDLE_PATH ${GEM_HOME}

# RUN bundle config build.nokogiri --use-system-libraries
# BUNDLE_BUILD__NOKOGIRI: "--use-system-libraries"
