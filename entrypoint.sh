#!/bin/sh
set -e

envsubst < /app/config/config.yaml > /app/config/config.yaml.tmp && mv /app/config/config.yaml.tmp /app/config/config.yaml

exec "$@"
