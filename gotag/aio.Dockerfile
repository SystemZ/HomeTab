FROM golang:latest
ADD . /build
RUN cd /build \
 && go get; exit 0
RUN cd /go/pkg/mod/github.com/discordapp/lilliput@v0.0.0-20191204003513-dd93dff726a5/deps/linux/lib \
 && cp libpng16.a libpng.a \
 && cp libpng16.la libpng.la \
 && cd /build \
 && go build -o /gotag

FROM alpine:latest
RUN apk add --no-cache bash tzdata openssl musl-dev ca-certificates
COPY --from=0 /gotag /
ADD tasktab /
ADD new /new
ADD /frontend/src /frontend
ENTRYPOINT ["/gotag"]
EXPOSE 4000