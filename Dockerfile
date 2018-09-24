FROM    golang:1.11-alpine3.8 as build-stage
RUN     apk --no-cache --update add curl git dep
COPY    . /go/src/github.com/travelping/ping-exporter
WORKDIR /go/src/github.com/travelping/ping-exporter
RUN     dep ensure -v && \
        CGO_ENABLED=0 GOOS=linux go build -v -a -o /ping-exporter .

FROM    alpine:3.8
ARG     VCS_REF="master"
LABEL   \
        org.label-schema.name="travelping/ping-exporter" \
        org.label-schema.vendor="Travelping GmbH" \
        org.label-schema.description="Ping (multiple) targets and export round trip times via an HTTP endpoint suitable for Prometheus" \
        org.label-schema.url="https://github.com/travelping/ping-exporter" \
        org.label-schema.vcs-url="https://github.com/travelping/ping-exporter" \
        org.label-schema.vcs-ref="$VCS_REF" \
        org.label-schema.docker.cmd.help="docker exec -it $CONTAINER /ping-exporter --help"


RUN     apk --no-cache add ca-certificates
WORKDIR /root/
COPY    --from=build-stage /ping-exporter .
CMD     ["./ping-exporter"]
