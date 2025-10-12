# Stage 1: build
FROM golang:1.24.6-alpine AS builder
WORKDIR /app

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build verbose จะเห็น error แบบเต็ม ๆ ถ้ามี
RUN go build -v -o auth-service ./cmd/main.go

# Stage 2: run
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/auth-service .
EXPOSE 8081
CMD ["./auth-service"]
