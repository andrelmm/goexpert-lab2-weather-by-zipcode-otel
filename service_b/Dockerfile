FROM golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY service_b/go.mod service_b/go.sum ./
COPY shared/go.mod shared/go.sum ../shared/

# Download all dependencies
RUN go mod download

# Copy the source code
COPY service_b/ ./service_b
COPY shared/ ./shared

# Build the application
WORKDIR /app/service_b
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/service_b/service_b

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service_b/service_b .
CMD ["./service_b"]