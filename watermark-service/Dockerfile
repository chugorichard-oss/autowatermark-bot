FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o watermark-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/watermark-service .
EXPOSE 5000
CMD ["./watermark-service"]
