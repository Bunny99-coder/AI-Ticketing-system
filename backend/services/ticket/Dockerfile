# Stage 1: Build
FROM golang:1.25-alpine AS builder
LABEL build-time="2025-11-01T10:28:33Z"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ticket-service cmd/ticket-service/main.go

# Stage 2: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ticket-service .
CMD ["./ticket-service"]