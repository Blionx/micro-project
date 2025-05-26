FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o products-service main.go handlers.go

EXPOSE 8082

CMD ["./products-service"]
