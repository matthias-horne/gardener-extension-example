FROM golang:1.27 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -tags netgo -o /gardener-extension-example ./cmd/gardener-extension-example

FROM alpine:3.24
COPY --from=builder /gardener-extension-example /gardener-extension-example
USER 65532:65532
ENTRYPOINT ["/gardener-extension-example"]
