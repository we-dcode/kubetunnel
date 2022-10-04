
FROM golang:1.18-bullseye as build-env

COPY ./ /go/src/github.com/dcode/kubetunnel/

WORKDIR /go/src/github.com/dcode/kubetunnel

# install all dependencies
RUN go get ./...

RUN go build -o /kubetunnel cmd/server/kubetunnel-server.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build-env /kubetunnel /kubetunnel

EXPOSE 8080

USER nonroot:nonroot

# Put back once we have an application
CMD ["/kubetunnel"]

