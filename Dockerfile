FROM ruby:3.2.1

# RUN bundle config --global frozen 1

ENV NODE_VERSION 19.x

RUN curl -sL https://deb.nodesource.com/setup_${NODE_VERSION} | bash -

RUN apt-get install -y nodejs git

ENV DOCKER_CHANNEL stable
ENV DOCKER_VERSION 20.10.21

RUN curl -fsSL "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz" \
| tar -xzC /usr/local/bin --strip=1 docker/docker

RUN curl -sS https://www.postgresql.org/media/keys/ACCC4CF8.asc | gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
RUN echo \
  "deb [signed-by=/etc/apt/trusted.gpg.d/postgresql.gpg] https://apt.postgresql.org/pub/repos/apt \
  `. /etc/os-release && echo "$VERSION_CODENAME"`-pgdg main" | \
  tee /etc/apt/sources.list.d/pgdg.list

RUN apt-get update && apt-get install -y \
  build-essential \
  bash-completion \
  postgresql-client-16 \
  libsqlite3-dev \
  libvips42 \
  && rm -rf /var/lib/apt/lists/*

# NOTE: for sorbet
WORKDIR /var/tmp

RUN wget -O watchman.zip https://github.com/facebook/watchman/releases/download/v2022.12.26.00/watchman-v2022.12.26.00-linux.zip && unzip watchman.zip

WORKDIR /var/tmp/watchman-v2022.12.26.00-linux

RUN mkdir -p /usr/local/{bin,lib} /usr/local/var/run/watchman
RUN cp bin/* /usr/local/bin
RUN cp lib/* /usr/local/lib
RUN chmod 755 /usr/local/bin/watchman
RUN chmod 2777 /usr/local/var/run/watchman

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
