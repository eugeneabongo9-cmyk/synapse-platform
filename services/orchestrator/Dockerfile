FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build deps if needed
RUN apk add --no-cache git ca-certificates

COPY go.mod ./
RUN go mod download || true

COPY . .

# Resolve all dependencies and build
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /orchestrator ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates curl
WORKDIR /
COPY --from=builder /orchestrator /orchestrator

EXPOSE 8080
USER nobody

ENTRYPOINT ["/orchestrator"]
