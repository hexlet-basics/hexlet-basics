FROM nginx:1.23.1

COPY --from=hexletbasics/services-app /opt/projects/hexlet-basics/public /var/www/hexlet-basics/shared/public

COPY services/web-nginx/files/nginx.hexlet.basics.conf.template /etc/nginx/conf.d/hexlet-basics.conf.template

# for security reason
RUN rm -rf /etc/nginx/conf.d/default.conf

COPY services/web-nginx/files/docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
