# Forum Service

Микросервисный форум с чатом, реализованный на Go.

## Архитектура

Проект состоит из двух микросервисов:
- `forum_service` - основной сервис форума
- `auth_service` - сервис аутентификации

## Технологии

- Go
- PostgreSQL
- gRPC
- WebSocket
- Gin
- Swagger
- golang-migrate
- zap (логирование)

## Требования

- Go 1.21+
- PostgreSQL 14+
- Docker (опционально)

## Установка и запуск

1. Клонировать репозиторий:
```bash
git clone https://github.com/your-username/miqxzzforum.git
cd miqxzzforum
```

2. Установить зависимости:
```bash
go mod download
```

3. Настроить переменные окружения:
```bash
cp .env.example .env
# Отредактировать .env файл
```

4. Запустить миграции:
```bash
cd forum-backend/forum_service
migrate -path migrations -database "postgresql://user:password@localhost:5432/forum?sslmode=disable" up
```

5. Запустить сервисы:
```bash
# Запуск forum_service
cd forum-backend/forum_service
go run cmd/main.go

# Запуск auth_service
cd forum-backend/auth_service
go run cmd/main.go
```

## API Документация

Swagger UI доступен по адресу:
- Forum Service: http://localhost:8081/swagger/index.html
- Auth Service: http://localhost:8082/swagger/index.html

## Тесты

Запуск тестов:
```bash
go test ./... -v
```

## Лицензия

MIT 