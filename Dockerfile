FROM golang:1.25.2-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]

