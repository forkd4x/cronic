#!/bin/sh
# cronic:
#   name: Example Shell Job
#   desc: Say hello every 10 seconds
#   cron: */10 * * * * *

echo "Hello, from $(basename "$0") using $SHELL"
sleep 5
echo "Bye, from $(basename "$0") using $SHELL"

