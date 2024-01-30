FROM golang:1.21-alpine AS builder
LABEL authors="irudenko"

# Устанавливаем рабочую директорию
WORKDIR /app

COPY go.* /app/

RUN go mod download

# Копируем исходный код приложения
COPY . /app/

# Собираем бинарный файл приложения с отключением CGO
RUN CGO_ENABLED=0 go build -o notifier ./cmd/app/main.go

# Стадия запуска
FROM alpine

# Копируем бинарный файл из стадии сборки
COPY --from=builder /app/notifier /app/notifier

WORKDIR /app

# Запускаем приложение
ENTRYPOINT ["./notifier"]