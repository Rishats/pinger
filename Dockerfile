# Этап 1: Сборка приложения
FROM golang:1.23-alpine as builder

# Устанавливаем зависимости и создаем рабочую директорию
WORKDIR /app

# Копируем все исходники в контейнер
COPY . .

# Устанавливаем зависимости (если они есть в go.mod и go.sum)
RUN go mod tidy

# Сборка бинарного файла
RUN go build -o pinger .

# Этап 2: Финальный образ для запуска
FROM alpine:latest

# Устанавливаем CA certificates для https-запросов
RUN apk --no-cache add ca-certificates

# Копируем скомпилированное приложение из этапа сборки
COPY --from=builder /app/pinger /usr/local/bin/pinger

# Порт, на котором будет работать ваш сервис
EXPOSE 2112

# Команда для запуска вашего пингера
CMD ["/usr/local/bin/pinger"]

