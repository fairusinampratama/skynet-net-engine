# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Dependency Caching
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
# We intentionally ignore the 'web' folder via .dockerignore, so it won't be copied
RUN CGO_ENABLED=0 GOOS=linux go build -o /net-engine cmd/server/main.go

# Runtime Stage
FROM gcr.io/distroless/static-debian12

WORKDIR /
COPY --from=builder /net-engine /net-engine

# Expose API Port
EXPOSE 8080

ENTRYPOINT ["/net-engine"]
