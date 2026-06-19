# Subscriptions Service

REST API сервис для управления онлайн-подписками пользователей.

Проект реализован как тестовое задание для позиции **Junior Golang Developer (Effective Mobile)**.

---

# Возможности

Сервис поддерживает:

* Создание подписки
* Получение подписки по UUID
* Получение списка всех подписок
* Обновление подписки
* Удаление подписки
* Подсчёт суммарной стоимости подписок за выбранный период
* Фильтрацию по:

  * `user_id`
  * `service_name`

---

# Технологический стек

* Go 1.25
* [Gin](https://gin-gonic.com/?utm_source=chatgpt.com)
* [PostgreSQL](https://www.postgresql.org/?utm_source=chatgpt.com)
* [lib/pq](https://github.com/lib/pq?utm_source=chatgpt.com)
* Docker / Docker Compose
* [Swagger (swaggo)](https://github.com/swaggo/swag?utm_source=chatgpt.com)
* [Zap Logger](https://pkg.go.dev/go.uber.org/zap?utm_source=chatgpt.com)

---

# Архитектура проекта

```text
.
├── cmd/app                 # Entry point
├── docs                    # Swagger documentation
└── internal
    ├── config             # Config (.env + yaml)
    ├── database           # DB connection
    ├── dto                # Request / Response DTO
    ├── errs               # Custom errors
    ├── handler            # HTTP handlers
    ├── logger             # Logger
    ├── migrations         # SQL migrations
    ├── model              # Database models
    ├── repository         # Data access layer
    ├── response           # API responses
    ├── router             # Routes
    └── service            # Business logic
```

---

# Конфигурация

Создать файл `.env`:

```env
DB_HOST=--------
DB_PORT=----
DB_USER=-----
DB_PASSWORD=--
```

Дополнительные настройки находятся в:

```text
internal/config/config.yml
```

---

# Запуск локально

### Установить зависимости

```bash
go mod tidy
```

### Запустить приложение

```bash
go run ./cmd/app
```

---

# Запуск через Docker

```bash
docker-compose up --build
```

---

# Миграции

Миграции находятся в:

```text
internal/migrations
```

Применяются **автоматически при запуске приложения**.

---

# Swagger Documentation

После запуска приложения Swagger доступен по адресу:

```text
http://localhost:8080/swagger/index.html
```

---

# API Endpoints

## CRUD

### Создать подписку

```http
POST /api/subscriptions
```

### Получить список подписок

```http
GET /api/subscriptions
```

### Получить подписку по ID

```http
GET /api/subscriptions/{id}
```

### Обновить подписку

```http
PUT /api/subscriptions/{id}
```

### Удалить подписку

```http
DELETE /api/subscriptions/{id}
```

---

## Подсчёт общей стоимости подписок

### Без фильтрации

```http
GET /api/subscriptions/total?from=07-2025&to=07-2026
```

### С фильтрацией

```http
GET /api/subscriptions/total?from=07-2025&to=07-2026&user_id=<uuid>&service_name=Yandex Plus
```

---

# Тестирование

Запуск unit tests:

```bash
go test ./...
```

---

# Логирование

Сервис использует structured logging через Zap logger.

Логи включают:

* HTTP requests
* Ошибки
* Database events
* Business events

---

# Автор:

**Shamsiddin Toshzoda**
