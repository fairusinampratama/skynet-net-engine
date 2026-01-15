# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy Source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o skynet-net-engine-api cmd/server/main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS (if needed) and timezone
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/skynet-net-engine-api .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./skynet-net-engine-api"]
