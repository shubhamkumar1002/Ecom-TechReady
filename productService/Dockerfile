FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o product-service .

FROM alpine:latest

WORKDIR /root

COPY --from=build /app/product-service /product-service

COPY --from=build /app/.env .env

EXPOSE 8080
CMD ["/product-service"]
