# ---- Build stage ----
FROM golang:1.24.4 AS builder

WORKDIR /app

# Download dependencies first (build caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o incentive_calculator .

# ---- Run stage ----
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/incentive_calculator .
COPY --from=builder /app/resources ./resources
COPY --from=builder /app/cert ./cert

# Create log directory and allow OpenShift random UID
RUN mkdir -p /app/Logs && chmod -R 777 /app

# Run the app
# Default command (can be overridden at runtime)
CMD ["./incentive_calculator", "--env", "prod", "--port", "8090", "--base-config-path", "./resources"]
