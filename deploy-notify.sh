#!/bin/bash

# curl -X POST https://sentry.io/api/0/organizations/hexlet-sq/releases/ \
#      -H "Authorization: Bearer ${SENTRY_API_TOKEN}" \
#      -H 'Content-Type: application/json' \
#      -d "{\"environment\":\"${RAILS_ENV}\", \"version\":\"${HEXLET_BASICS_RELEASE_VERSION}\", \"projects\":[\"hexlet-basics\"]}"

# curl -v -X POST https://sentry.io/api/0/organizations/hexlet-sq/releases/$HEXLET_BASICS_RELEASE_VERSION/deploys/ \
#      -H "Authorization: Bearer ${SENTRY_API_TOKEN}" \
#      -H 'Content-Type: application/json' \
#      -d "{\"environment\":\"${RAILS_ENV}\"}"

# NOTE: релиз создаётся ранее при сборке на github. Потому создаём только деплой к релизу
curl -v -X POST https://sentry.hexlet.io/api/0/organizations/hexlet/releases/$HEXLET_BASICS_RELEASE_VERSION/deploys/ \
     -H "Authorization: Bearer ${SENTRY_API_TOKEN}" \
     -H 'Content-Type: application/json' \
     -d "{\"environment\":\"${RAILS_ENV}\"}"
