FROM golang:latest AS builder

WORKDIR /app

COPY service_a/go.mod service_a/go.sum ./
COPY shared/go.mod shared/go.sum ./shared/

RUN go mod download

COPY service_a/ .
COPY shared/ ./shared/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service_a ./service_a

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service_a .
CMD ["./service_a"]