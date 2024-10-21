# Usar una imagen base de Go
FROM golang:1.23 AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum y descargar las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar la aplicación
RUN go build -o /rest-api ./api-rest/app.go

# Exponer el puerto que usará la aplicación
EXPOSE 3002

# Comando para ejecutar la aplicación
ENTRYPOINT ["/rest-api"]
