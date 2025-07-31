# cronic:
#   name: Example Dockerfile Job
#   desc: Say hello every 40 seconds
#   cron: */40 * * * * *
#   cmd: docker build -f $f -t ${f%.*} . && docker run --rm ${f%.*}
FROM alpine:latest
CMD echo Hello, from container $HOSTNAME && \
    sleep 20 && \
    echo Bye, from container $HOSTNAME
