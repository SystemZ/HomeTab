FROM alpine:latest
ADD tasktab /
RUN apk add --no-cache bash
ENTRYPOINT ["/tasktab"]
EXPOSE 3000