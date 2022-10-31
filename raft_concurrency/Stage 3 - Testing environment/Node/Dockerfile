# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

WORKDIR /app
COPY . .

# RUN apk add build-base
RUN go mod download

RUN go build -o main
CMD ["./main"]

EXPOSE 5555

