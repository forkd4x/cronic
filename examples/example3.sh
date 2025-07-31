#!/bin/sh
# cronic:
#   name: Example Shell Job
#   desc: Say hello every 30 seconds
#   cron: */30 * * * * *

echo "Hello, from $(basename "$0") using $SHELL"
sleep 15
echo "Bye, from $(basename "$0") using $SHELL"

