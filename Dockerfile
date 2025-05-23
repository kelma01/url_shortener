FROM golang:1.24 AS builder

WORKDIR /url_shortener

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "./cmd/server/main.go"]

EXPOSE 8080