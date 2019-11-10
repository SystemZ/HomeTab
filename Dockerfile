FROM alpine:latest
RUN apk add --no-cache bash tzdata
ADD tasktab /
ADD templates /templates
ENTRYPOINT ["/tasktab"]
EXPOSE 3000