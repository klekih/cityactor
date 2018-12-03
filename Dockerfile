FROM golang

ADD . /go/src/github.com/klekih/cityactor

RUN go install github.com/klekih/cityactor

ENTRYPOINT /go/bin/cityactor
