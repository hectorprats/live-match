FROM golang:latest

RUN mkdir /go/src/app
COPY /src /go/src/app
WORKDIR /go/src/app

RUN go get app
RUN go install

CMD ["/go/bin/app"]