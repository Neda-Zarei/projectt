FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /userplan-service ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /userplan-service .
# COPY --from=builder /app/configs /app/configs
EXPOSE 50051
CMD ["./userplan-service"]