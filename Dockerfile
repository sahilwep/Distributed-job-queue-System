FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api
RUN go build -o worker ./cmd/worker

EXPOSE 8080

CMD ["./api"]