FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /management-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /management-service .
COPY --from=builder /app/configs /app/configs
EXPOSE 8080
CMD ["./management-service"]