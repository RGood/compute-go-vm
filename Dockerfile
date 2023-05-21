FROM golang:alpine as build
WORKDIR /src

COPY cmd/ cmd
COPY internal/ internal
COPY go.mod go.mod
COPY go.sum go.sum

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/main/main.go

FROM alpine:3.18 as production
WORKDIR /srv

RUN apk update && apk add --no-cache docker-cli
RUN apk add --no-cache bash

COPY --from=build /src/server server

CMD [ "/srv/server" ]
