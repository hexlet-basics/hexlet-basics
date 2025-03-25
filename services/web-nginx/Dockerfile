ARG APP_IMAGE_WITH_TAG=hexletbasics/services-app
FROM ${APP_IMAGE_WITH_TAG} AS rails-app

FROM nginx:1.27.4

COPY --from=rails-app /rails/public /var/www/code-basics/shared/public

COPY services/web-nginx/files/nginx.hexlet.basics.conf /etc/nginx/conf.d/code-basics.conf

# for security reason
RUN rm -rf /etc/nginx/conf.d/default.conf

CMD ["nginx", "-g", "daemon off;"]
