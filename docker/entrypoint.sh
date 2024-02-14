#!/bin/sh

# Check if the secret file exists at /run/secrets/pixivfe_token
if [ -f /run/secrets/pixivfe_token ]; then
    export PIXIVFE_TOKEN=$(cat /run/secrets/pixivfe_token)
    echo "Info: PIXIVFE_TOKEN loaded from secret."
else
    echo "Info: PIXIVFE_TOKEN not loaded from secret. Loading the environment variable normally."
fi

# Execute the main application
exec /app/pixivfe "$@"
