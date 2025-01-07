FROM alpine:latest
RUN apk update && apk add git go

USER root

COPY . /go-mysql
WORKDIR /go-mysql

RUN go build -o /go-mysqld github.com/cybergarage/go-mysql/examples/go-mysqld

ENTRYPOINT ["/go-mysqld"]
