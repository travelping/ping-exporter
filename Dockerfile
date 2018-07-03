FROM golang:1.10-alpine as builder
RUN apk --no-cache add curl git && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY . /go/src/github.com/travelping/ping-exporter
WORKDIR /go/src/github.com/travelping/ping-exporter
RUN dep ensure && CGO_ENABLED=0 GOOS=linux go build -a -o ping-exporter .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/travelping/ping-exporter/ping-exporter .
RUN mkdir -p /etc/ping-exporter/ && touch /etc/ping-exporter/ping-exporter.yaml
CMD ["./ping-exporter"]
