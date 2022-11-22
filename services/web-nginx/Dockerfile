FROM ruby:3.1.2 as rails

ENV PROJECT_ROOT /app
WORKDIR ${PROJECT_ROOT}

ENV NODE_VERSION 19.x

RUN curl -sL https://deb.nodesource.com/setup_${NODE_VERSION} | bash -

RUN apt-get install -y nodejs

COPY package.json .
COPY package-lock.json .

RUN npm install

COPY Gemfile .
COPY Gemfile.lock .

RUN bundle install --jobs $(nproc)

COPY . .

RUN APP_HOST=example.test bin/rails js:routes:typescript
RUN npm run build
RUN APP_HOST=example.test NODE_ENV=production RAILS_ENV=production SECRET_KEY_BASE=jopa bin/rails assets:precompile

FROM nginx:1.23.1

COPY --from=rails /app/public /var/www/hexlet-basics/shared/public

COPY services/web-nginx/files/nginx.hexlet.basics.conf.template /etc/nginx/conf.d/hexlet-basics.conf.template

# for security reason
RUN rm -rf /etc/nginx/conf.d/default.conf

COPY services/web-nginx/files/docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
