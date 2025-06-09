# Usa una imagen oficial de Go
FROM golang:1.24-alpine

# Instala herramientas necesarias
RUN apk add --no-cache git

# Crea el directorio de trabajo
WORKDIR /app

# Copia los archivos go.mod y go.sum y descarga dependencias primero (cache eficiente)
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto del código
COPY . .

# Compila el binario
RUN go build -o user-service main.go

# Expone el puerto del microservicio
EXPOSE 8083

# Comando de inicio
CMD ["./user-service"]
