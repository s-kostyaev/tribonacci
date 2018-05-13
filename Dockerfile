FROM golang:1.10-alpine AS build-env

ENV GOPATH=/ \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

RUN apk update && \
    apk add git && \
    apk add make

WORKDIR /src/github.com/s-kostyaev/tribonacci
ADD . .

RUN make    


FROM alpine:3.7
MAINTAINER s-kostyaev

LABEL name=tribonacci
LABEL version=0.0.1
LABEL architecrture=amd64
LABEL source="ssh://git@github.com:s-kostyaev/tribonacci.git"

RUN mkdir /app
COPY --from=build-env /src/github.com/s-kostyaev/tribonacci/bin/tribonacci-web /app/tribonacci-web
COPY ./cmd/tribonacci-web/docs /app/docs
WORKDIR /app/

EXPOSE 8080
ENTRYPOINT ["/app/tribonacci-web"]
