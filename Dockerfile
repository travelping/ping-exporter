FROM golang:1.10-alpine as builder
RUN apk --no-cache add curl git && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY . /go/src/github.com/openvnf/cgw-exporter/
WORKDIR /go/src/github.com/openvnf/cgw-exporter/
RUN dep ensure && CGO_ENABLED=0 GOOS=linux go build -a -o cgw-exporter .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/openvnf/cgw-exporter/cgw-exporter .
RUN mkdir -p /etc/defaults/cgw-exporter/ && touch /etc/defaults/cgw-exporter/cgw-exporter.yaml
CMD ["./cgw-exporter"]
