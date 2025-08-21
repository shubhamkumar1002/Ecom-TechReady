FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o payment-service .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/payment-service /payment-service

COPY --from=build /app/.env .env

EXPOSE 8080

CMD ["/payment-service"]
