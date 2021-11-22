FROM golang:1.17.3-alpine

ADD . /go/src

WORKDIR /go/src

RUN go get -d -v ./...

RUN go install -v ./...
