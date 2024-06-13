FROM golang AS builder
LABEL REPO="www.github.com/marvasgit/KubeStateWatch"

RUN apt-get update && \
    dpkg --add-architecture arm64 &&\
    apt-get install -y --no-install-recommends build-essential && \
    apt-get clean && \
    mkdir -p "$GOPATH/src/github.com/marvasgit/KubeStateWatch"

ADD . "$GOPATH/src/github.com/marvasgit/KubeStateWatch"

RUN cd "$GOPATH/src/github.com/marvasgit/KubeStateWatch" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=$(dpkg --print-architecture) go build -a --installsuffix cgo --ldflags="-s" -o /statemonitor

FROM cgr.dev/chainguard/static:latest-glibc

COPY --from=builder /statemonitor /bin/statemonitor

ENTRYPOINT ["/bin/statemonitor"]
