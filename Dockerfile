# syntax=docker/dockerfile:1

FROM golang:1.19.2-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY src/ ./src

RUN go build -o /webserver ./src

CMD ["/webserver"]