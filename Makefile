IMAGE_TAG ?= latest
REGISTRY ?= ghcr.io/gardener
IMAGE := $(REGISTRY)/gardener-extension-example:$(IMAGE_TAG)
ADMISSION_IMAGE := $(REGISTRY)/gardener-extension-admission-example:$(IMAGE_TAG)
GOARCH ?= $(shell go env GOARCH)
GARDENER_HACK_DIR = $(shell go list -m -f '{{.Dir}}' github.com/gardener/gardener)/hack

.PHONY: build install tidy test clean docker-image docker-image-extension docker-image-admission docker-push generate push-to-gardener-local

build:
	go build -tags netgo -o bin/gardener-extension-example ./cmd/gardener-extension-example
	go build -tags netgo -o bin/gardener-extension-admission-example ./cmd/gardener-extension-admission-example

install:
	go install -tags netgo ./cmd/...

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -rf bin

docker-image: docker-image-extension docker-image-admission

docker-image-extension:
	docker build --target gardener-extension-example --tag $(IMAGE) .

docker-image-admission:
	docker build --target gardener-extension-admission-example --tag $(ADMISSION_IMAGE) .

docker-push:
	docker push $(IMAGE)
	docker push $(ADMISSION_IMAGE)

generate:
	bash hack/update-codegen.sh
	GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) IMAGE_TAG=$(IMAGE_TAG) go generate ./charts/...

push-to-gardener-local:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -tags netgo -o bin/gardener-extension-example ./cmd/gardener-extension-example
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -tags netgo -o bin/gardener-extension-admission-example ./cmd/gardener-extension-admission-example
	docker build -f Dockerfile.dev --target gardener-extension-example -t $(IMAGE) .
	docker build -f Dockerfile.dev --target gardener-extension-admission-example -t $(ADMISSION_IMAGE) .
	kind --name gardener-local load docker-image $(IMAGE) $(ADMISSION_IMAGE)
