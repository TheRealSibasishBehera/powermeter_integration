FROM golang:1.18 as builder

USER root

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o main ./cmd/

ENTRYPOINT ["/app/main"]
