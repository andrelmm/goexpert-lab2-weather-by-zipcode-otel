FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY service_a/ .

COPY ot ./ot

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]