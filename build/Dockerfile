# this failed
#FROM gcr.io/distroless/base:latest

# this probably needs CGO=0
#FROM alpine:latest
#RUN apk add --no-cache bash tzdata openssl musl-dev ca-certificates

FROM ubuntu:20.04

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates tzdata\
 && rm -rf /var/lib/{apt,dpkg,cache,log}/ \
 && apt-get clean autoclean
ADD hometab /
RUN chmod +x /hometab
ADD frontend /new
ADD templates /templates
ENTRYPOINT ["/hometab"]
