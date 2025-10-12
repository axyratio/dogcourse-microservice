# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o course-service ./cmd/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/course-service .
COPY .env .

EXPOSE 8082
CMD ["./course-service"]
