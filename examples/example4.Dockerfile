# TODO: Set filename as ENV inside docker container
# cronic:
#   name: Example Dockerfile Job
#   desc: Say hello every 14 seconds
#   cron: */14 * * * * *
FROM alpine:latest
CMD echo Hello, from container $HOSTNAME && \
    sleep 7 && \
    echo Bye, from container $HOSTNAME
