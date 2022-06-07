# FROM golang:1.18-alpine AS build-env

# ENV GOPATH=/ \
#     GOOS=linux \
#     GOARCH=amd64 \
#     CGO_ENABLED=0

# RUN apk update && \
#     apk add git && \
#     apk add make

# WORKDIR /root/
# ADD . .

# RUN make    


FROM golang:1.18-alpine
MAINTAINER s-kostyaev

LABEL name=tribonacci
LABEL version=0.0.1
LABEL architecrture=amd64
LABEL source="ssh://git@github.com:s-kostyaev/tribonacci.git"

RUN apk add build-base
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN mkdir /app
COPY ./bin/tribonacci-web /app/tribonacci-web
COPY ./cmd/tribonacci-web/docs /app/docs
WORKDIR /app/

EXPOSE 8080
EXPOSE 2345
ENTRYPOINT ["dlv", "--headless", "exec", "/app/tribonacci-web", "-l", ":2345", "--log", "--log-output", "debugger,rpc,dap"]
