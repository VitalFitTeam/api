FROM golang:1.25.0-alpine AS development


RUN apk add --no-cache git make
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go install github.com/air-verse/air@latest

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY .air.toml go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080