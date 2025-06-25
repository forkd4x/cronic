#!/bin/sh
# cronic:
#   name: Example Shell Job
#   desc: Say hello every 7 seconds
#   cron: */7 * * * * *

echo "Hello, from $(basename "$0") using $SHELL"
sleep 3
echo "Hello, again, from $(basename "$0") using $SHELL"

