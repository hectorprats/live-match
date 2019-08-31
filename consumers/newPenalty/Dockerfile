FROM golang:latest AS builder

RUN mkdir /go/src/app
COPY /src /go/src/app
WORKDIR /go/src/app

RUN go get app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/newPenalty

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/newPenalty /go/bin/newPenalty
COPY src/config.json /go/bin/config.json
RUN chmod +x /go/bin/newPenalty

WORKDIR /go/bin/

ENTRYPOINT ["/go/bin/newPenalty"]