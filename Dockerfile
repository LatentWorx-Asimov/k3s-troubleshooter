FROM golang:1.23-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the worker binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /agent-worker ./cmd/worker

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /agent-worker /agent-worker
USER 65532:65532

ENTRYPOINT ["/agent-worker"]
