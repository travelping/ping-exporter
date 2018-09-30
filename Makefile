ping-exporter: vendor
	go build -v -o $@

# fetch dependencies as defined via Gopkg.{toml,lock}
vendor:
	dep ensure -v

build-docker-image:
	docker build -f Dockerfile -t travelping/ping-exporter:latest .

.PHONY: ping-exporter
