FROM golang:1.20 as build
WORKDIR /src

COPY cmd/ cmd
COPY internal/ internal
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o server cmd/main/main.go

FROM ubuntu:22.04 as production
WORKDIR /srv

RUN apt-get update && apt-get install docker.io -y

COPY --from=build /src/server server

CMD [ "./server" ]
