FROM    golang:1.11-alpine3.8 as build-stage
RUN     apk --no-cache --update add curl git dep
COPY    . /go/src/github.com/travelping/ping-exporter
WORKDIR /go/src/github.com/travelping/ping-exporter
RUN     dep ensure -v && \
        CGO_ENABLED=0 GOOS=linux go build -v -a -o /ping-exporter .

FROM    alpine:3.8
RUN     apk --no-cache add ca-certificates
WORKDIR /root/
COPY    --from=build-stage /ping-exporter .
CMD     ["./ping-exporter"]
