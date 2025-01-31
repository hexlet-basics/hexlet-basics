#!/bin/bash
curl -X POST https://sentry.io/api/0/organizations/hexlet-sq/releases/ \
     -H "Authorization: Bearer ${SENTRY_API_TOKEN}" \
     -H 'Content-Type: application/json' \
     -d "{\"environment\":\"${RAILS_ENV}\", \"version\":\"${HEXLET_BASICS_RELEASE_VERSION}\", \"projects\":[\"hexlet-basics\"]}"

curl -v -X POST https://sentry.io/api/0/organizations/hexlet-sq/releases/$HEXLET_BASICS_RELEASE_VERSION/deploys/ \
     -H "Authorization: Bearer ${SENTRY_API_TOKEN}" \
     -H 'Content-Type: application/json' \
     -d "{\"environment\":\"${RAILS_ENV}\"}"
