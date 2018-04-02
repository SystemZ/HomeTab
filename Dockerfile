FROM alpine:latest
ADD tasktab /
RUN chmod +x /tasktab
ENTRYPOINT /tasktab
EXPOSE 3000