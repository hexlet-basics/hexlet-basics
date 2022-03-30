#!/bin/bash
# https://sentry.io/settings/hexlet-sq/projects/hexlet-basics/release-tracking/
curl https://sentry.io/api/hooks/release/builtin/6298540/5214185412a3b448f1afde83bed4e5965fa199755a50cd5a074dab276bf63d1c/ \
  -X POST \
  -H 'Content-Type: application/json' \
  -d '{"version": "abcdefg"}'
