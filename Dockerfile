# ====== Stage 1: Build the Go binary ======
FROM golang:1.23 AS builder
WORKDIR /app
# Copy go.mod and go.sum
# COPY go.mod go.sum ./
# RUN go mod download
# Copy the source code
COPY . .
# Build the Go binary
RUN GO111MODULE=on CGO_ENABLED=0 GOARCH="amd64" GOOS="linux"  go build -trimpath --tags=kqueue --ldflags "-s -w" -o cecp ./cmd

# ====== Stage 2: Minimal runtime container ======
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /app/cecp .
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown appuser:appgroup ./cecp
USER appuser
CMD ["./cecp"]