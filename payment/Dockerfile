FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o payment-service .

FROM alpine:latest

WORKDIR /app

COPY ./docs ./docs

COPY --from=build /app/payment-service /payment-service
# --from=build /app/.env .env

COPY --from=build /app/serviceaccountgcp.json /app/credentials.json

ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials.json

EXPOSE 8080

CMD ["/payment-service"]
