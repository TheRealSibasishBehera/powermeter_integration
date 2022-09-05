FROM golang:1.18 as builder

USER root

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN apt-get update && apt-get install -y iputils-ping
COPY . .

RUN go build -o main ./cmd/
EXPOSE 8881

ENTRYPOINT ["/app/main"]
