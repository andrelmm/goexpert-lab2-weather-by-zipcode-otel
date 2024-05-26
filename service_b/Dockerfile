FROM golang:latest as builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]