IMAGE_TAG ?= latest
REGISTRY ?= ghcr.io/gardener
IMAGE := $(REGISTRY)/gardener-extension-example:$(IMAGE_TAG)
GOARCH ?= $(shell go env GOARCH)
GARDENER_HACK_DIR = $(shell go list -m -f '{{.Dir}}' github.com/gardener/gardener)/hack

.PHONY: build install tidy test clean docker-image docker-push generate push-to-gardener-local

build:
	go build -tags netgo -o bin/gardener-extension-example ./cmd/gardener-extension-example

install:
	go install -tags netgo ./cmd/gardener-extension-example

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -rf bin

docker-image:
	docker build --tag $(IMAGE) .

docker-push:
	docker push $(IMAGE)

generate:
	bash hack/update-codegen.sh
	GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) IMAGE_TAG=$(IMAGE_TAG) go generate ./charts/...

push-to-gardener-local:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -tags netgo -o bin/gardener-extension-example ./cmd/gardener-extension-example
	docker build -f Dockerfile.dev -t $(IMAGE) .
	kind --name gardener-local load docker-image $(IMAGE)
