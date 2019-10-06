FROM alpine:latest
RUN apk add --no-cache bash
ADD tasktab /
ADD templates /templates
ENTRYPOINT ["/tasktab"]
EXPOSE 3000