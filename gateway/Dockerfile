FROM golang:1.24-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gateway main.go router.go auth_middleware.go

EXPOSE 8080

CMD ["./gateway"]
