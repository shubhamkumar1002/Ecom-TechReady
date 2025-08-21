FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service .

FROM alpine:latest

WORKDIR /root

COPY --from=build /app/order-service /order-service

COPY --from=build /app/.env .env

EXPOSE 8080
CMD ["/order-service"]
