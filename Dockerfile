FROM ruby:2.6

ENV PROJECT_ROOT /app

RUN mkdir $PROJECT_ROOT
WORKDIR $PROJECT_ROOT

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list

RUN apt-get update -qq \
  && apt-get install -y --no-install-recommends \
  build-essential \
  libpq-dev \
  libsqlite3-dev  \
  nodejs \
  yarn

COPY Gemfile Gemfile
COPY Gemfile.lock Gemfile.lock

RUN gem install bundler
RUN bundle config path vendor/bundle
RUN bundle install --jobs 4 --retry 3

ADD . $PROJECT_ROOT
