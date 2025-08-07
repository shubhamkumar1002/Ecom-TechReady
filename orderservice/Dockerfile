FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service .

FROM alpine:latest

WORKDIR /root

COPY --from=build /app/order-service /order-service

COPY --from=build /app/serviceaccountgcp.json /app/credentials.json

ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials.json

EXPOSE 8080
CMD ["/order-service"]
