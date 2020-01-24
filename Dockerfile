# docker run -it --rm --name=gobuild -v "$PWD":/usr/src/jblastor -w /usr/src/jblastor 2205a315f9c7

FROM golang:1.12.0-alpine3.9
MAINTAINER John F. Hogarty <hogihung@gmail.com>

RUN apk add git

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
