#!/usr/bin/env sh
set -eu

envsubst '${DO_SPACES_REGION} ${DO_SPACES_SITEMAP_BUCKET}' < /etc/nginx/conf.d/hexlet-basics.conf.template > /etc/nginx/conf.d/hexlet-basics.conf

exec "$@"
