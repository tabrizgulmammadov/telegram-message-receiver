FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

# Create storage directory
RUN mkdir -p /messages/voices /messages/texts
VOLUME ["/messages"]

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["./main"]