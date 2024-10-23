FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /rest-api ./api-rest/app.go

EXPOSE 3002

ENTRYPOINT ["/rest-api"]
