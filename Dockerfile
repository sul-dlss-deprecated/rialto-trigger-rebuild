FROM golang:alpine as BASE
WORKDIR /go/src/github.com/sul-dlss/rialto-trigger-rebuild
COPY . .
RUN apk update && \
    apk add --no-cache --virtual .build-deps git && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    apk del .build-deps
RUN CGO_ENABLED=0 GOOS=linux go build -o trigger-rebuild -ldflags "-s" -a -installsuffix cgo main.go

FROM scratch
COPY --from=BASE /go/src/github.com/sul-dlss/rialto-trigger-rebuild/trigger-rebuild .
CMD ["./trigger-rebuild"]
