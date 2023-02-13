# syntax=docker/dockerfile:1
FROM golang:1.16-alpine
WORKDIR /app

COPY . .

RUN go mod download
EXPOSE 8080

RUN go build -o main
CMD ["./main"]