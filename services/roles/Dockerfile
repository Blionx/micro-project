FROM golang:1.24-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go mod tidy
RUN go build -o roles-service

EXPOSE 8083

CMD ["./roles-service"]
