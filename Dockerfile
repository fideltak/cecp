# ====== Stage 1: Build the Go binary ======
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOARCH="amd64" GOOS="linux"  go build -trimpath --tags=kqueue --ldflags "-s -w" -o cecp ./cmd

# ====== Stage 2: Runtime container ======
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /app/cecp .
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown appuser:appgroup ./cecp
USER appuser
CMD ["./cecp"]