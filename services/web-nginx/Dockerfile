ARG APP_IMAGE_WITH_TAG=hexletbasics/services-app
FROM ${APP_IMAGE_WITH_TAG} AS rails-app

FROM nginx:1.27.4

COPY --from=rails-app /rails/public /var/www/code-basics/shared/public

COPY services/web-nginx/files/nginx.code-basics.conf /etc/nginx/conf.d/code-basics.conf
COPY services/web-nginx/files/nginx.cors.conf /etc/nginx/snippets/cors.conf

COPY services/web-nginx/system/maintenance.html /system/maintenance.html
COPY services/web-nginx/system/maintenance.gif /system/maintenance.gif
COPY services/web-nginx/system/favicon.ico /system/favicon.ico

# for security reason
RUN rm -rf /etc/nginx/conf.d/default.conf

CMD ["nginx", "-g", "daemon off;"]
