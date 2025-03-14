#!/bin/sh

envsubst < /app/config.template.json > /app/config.json

# Copy initial files to the volume if they don't exist
if [ -d "/app/storage" ]; then
  cp -n -r /app/storage/* /storage/
fi

exec ./main