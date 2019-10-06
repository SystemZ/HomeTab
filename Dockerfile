FROM alpine:latest
ADD tasktab /
RUN apk add --no-cache \
    libc6-compat bash
ENTRYPOINT /tasktab
EXPOSE 3000