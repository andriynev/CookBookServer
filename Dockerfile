FROM golang:1.12-alpine as base

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl make gcc musl-dev openssh-client inotify-tools libc6-compat && \
    git config --global http.sslVerify false

ENV CGO_ENABLED 0

WORKDIR /app

COPY ./bin/migrate /app/bin/migrate
COPY ./bin/swag /app/bin/swag

FROM base as dev

# Compile Delve
# RUN go get -u github.com/derekparker/delve/cmd/dlv
RUN go get -u github.com/githubnemo/CompileDaemon

COPY ./Makefile /app/Makefile
COPY ./go.mod /app/go.mod
RUN make dependencies
RUN go get github.com/ugorji/go/codec@none

COPY . .

VOLUME [ "/data/media" ]
ENTRYPOINT ["/bin/bash", "/app/entrypoint-DEV.sh"]

EXPOSE 80 40000

FROM base as prod

COPY ./Makefile /app/Makefile
COPY ./entrypoint.sh /app/entrypoint.sh
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

COPY ./src/api /app/src/api

RUN make dependencies

# Build go binary
RUN make build-app

VOLUME [ "/data/media" ]
ENTRYPOINT ["/bin/bash", "/app/entrypoint.sh"]
EXPOSE 80
