#!/bin/sh
# cronic:
#   name: Example Shell Job
#   desc: Say hello every 14 seconds
#   cron: */14 * * * * *

echo "Hello, from $(basename "$0") using $SHELL"
sleep 7
echo "Hello, again, from $(basename "$0") using $SHELL"

