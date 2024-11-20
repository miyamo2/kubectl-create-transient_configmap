ARG GO_VERSION=1.23

FROM --platform=linux/amd64 golang:${GO_VERSION}-bookworm AS build

ADD . /tmp
WORKDIR /tmp

COPY . .

RUN env GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -trimpath -o ./bin/batch ./main.go \
    && chmod +x /tmp/bin/batch

FROM --platform=linux/amd64 scratch

COPY --from=build /tmp/bin/batch /go/bin/batch

ENTRYPOINT ["/go/bin/batch"]