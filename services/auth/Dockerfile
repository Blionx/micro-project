FROM golang:1.24-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service main.go jwt.go handlers.go

EXPOSE 8081

CMD ["./auth-service"]
