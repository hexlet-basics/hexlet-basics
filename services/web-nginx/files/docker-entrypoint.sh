#!/usr/bin/env sh
set -eu

envsubst '${SITEMAPS_S3_HOST}' < /etc/nginx/conf.d/hexlet-basics.conf.template > /etc/nginx/conf.d/hexlet-basics.conf

exec "$@"
