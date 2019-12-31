FROM alpine:latest
RUN apk add --no-cache bash tzdata
ADD tasktab /
ADD new /new
ADD templates /templates
ENTRYPOINT ["/tasktab"]
EXPOSE 3000