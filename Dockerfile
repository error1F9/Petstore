FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/main .

COPY ./migrations ./migrations

COPY .env .

RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["/app/main"]