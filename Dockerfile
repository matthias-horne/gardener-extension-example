FROM golang:1.27 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -tags netgo -o /gardener-extension-example ./cmd/gardener-extension-example
RUN CGO_ENABLED=0 go build -tags netgo -o /gardener-extension-admission-example ./cmd/gardener-extension-admission-example

FROM alpine:3.24 AS base
USER 65532:65532

FROM base AS gardener-extension-admission-example
COPY --from=builder /gardener-extension-admission-example /gardener-extension-admission-example
ENTRYPOINT ["/gardener-extension-admission-example"]

FROM base AS gardener-extension-example
COPY --from=builder /gardener-extension-example /gardener-extension-example
ENTRYPOINT ["/gardener-extension-example"]
