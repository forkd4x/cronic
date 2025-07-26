# cronic:
#   name: Example Dockerfile Job
#   desc: Say hello every 14 seconds
#   cron: */14 * * * * *
#   cmd: docker build -f $f -t ${f%.*} . && docker run --rm ${f%.*}
FROM alpine:latest
CMD echo Hello, from container $HOSTNAME && \
    sleep 7 && \
    echo Bye, from container $HOSTNAME
