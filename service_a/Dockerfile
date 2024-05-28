FROM golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY service_a/go.mod service_a/go.sum ./
COPY shared/go.mod shared/go.sum ../shared/

# Download all dependencies
RUN go mod download

# Copy the source code
COPY service_a/ ./service_a
COPY shared/ ./shared

# Build the application
WORKDIR /app/service_a
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/service_a/service_a

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service_a/service_a .
CMD ["./service_a"]