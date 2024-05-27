FROM golang:latest AS builder

WORKDIR /app

COPY service_b/go.mod service_b/go.sum ./
COPY shared/go.mod shared/go.sum ./shared/

RUN go mod download

COPY service_b/ .
COPY shared/ ./shared/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service_b ./service_b

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service_b .
CMD ["./service_b"]