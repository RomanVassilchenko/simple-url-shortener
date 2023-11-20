FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o simple-url-shortener cmd/simple-url-shortener/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/. .

EXPOSE 8083

ENV CONFIG_PATH="./config/prod.yaml" \
    POSTGRES_USER=${POSTGRES_USER} \
    POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
    HTTP_SERVER_PASSWORD=${HTTP_SERVER_PASSWORD} \
    POSTGRES_URL="postgres"

# Command to run the application
CMD ["./simple-url-shortener"]
