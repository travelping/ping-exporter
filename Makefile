ping-exporter: force vendor
	go build -v

# fetch dependencies as defined via Gopkg.{toml,lock}
vendor:
	dep ensure -v

build-docker-image:
	docker build -f Dockerfile -t travelping/ping-exporter:latest .

.PHONY: force
