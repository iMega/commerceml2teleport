ARG GOLANG_VERSION=1.10
ARG GRPCURL_COMMIT

# Build application
FROM golang:1.10-alpine as builder

WORKDIR /go/src/github.com/imega/commerceml2teleport

COPY . .

RUN go build -v -o rel/app

# Build release container

FROM golang:${GOLANG_VERSION}-alpine
# HTTP сервер
EXPOSE 8080
# GRPC сервер
EXPOSE 9000

WORKDIR /
RUN apk add --update git && \
    git clone https://github.com/fullstorydev/grpcurl.git $GOPATH/src/github.com/fullstorydev/grpcurl && \
    cd $GOPATH/src/github.com/fullstorydev/grpcurl && \
    git checkout ${GRPCURL_COMMIT} -b temp_branch && \
    go get -v github.com/fullstorydev/grpcurl/cmd/grpcurl && \
    apk del -v git && \
    rm -rf $GOPATH/src/*

HEALTHCHECK --interval=10s --timeout=2s \
    CMD grpcurl -plaintext -d {} 127.0.0.1:9000 grpc.health.v1.Health.Check | grep -w "SERVING"

COPY --from=builder /go/src/github.com/imega/commerceml2teleport/rel/app .
CMD ["/app"]
