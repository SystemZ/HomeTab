FROM alpine:latest
RUN apk add --no-cache bash tzdata openssl musl-dev ca-certificates
ADD tasktab /
ADD new /new
ADD templates /templates
ENTRYPOINT ["/tasktab"]
