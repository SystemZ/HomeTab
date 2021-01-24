FROM alpine:latest
RUN apk add --no-cache bash tzdata openssl musl-dev ca-certificates
ADD hometab /
ADD new /new
ADD templates /templates
ENTRYPOINT ["/hometab"]
