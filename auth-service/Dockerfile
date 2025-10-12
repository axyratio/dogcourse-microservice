# Stage 1: build
FROM golang:1.24.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth-service ./cmd/main.go

# Stage 2: run
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/auth-service .
COPY .env .
EXPOSE 8081
CMD ["./auth-service"]
