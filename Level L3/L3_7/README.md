# 📦 Warehouse Control

REST API сервис для управления складом товаров с ролевой аутентификацией по JWT.

---

## 🚀 Технологический стек

| Библиотека | Назначение |
|---|---|
| [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) | HTTP-фреймворк |
| [github.com/wb-go/wbf](https://github.com/wb-go/wbf) | Внутренний фреймворк WB (ginext, dbpg, zlog) |
| [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) | Генерация и валидация JWT-токенов |
| [github.com/google/uuid](https://github.com/google/uuid) | Генерация UUID для ID сущностей |
| [github.com/ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv) | Загрузка конфигурации из `.env` |
| [github.com/pressly/goose/v3](https://github.com/pressly/goose) | Миграции базы данных |
| [github.com/rs/zerolog](https://github.com/rs/zerolog) | Структурированное логирование |
| PostgreSQL (через `lib/pq`) | База данных |

---

## 📁 Структура проекта

```
L3_7/
├── cmd/
│   └── warehouse_control/
│       └── main.go              # Точка входа
├── internal/
│   ├── app/
│   │   └── app.go               # Инициализация и запуск сервера
│   ├── config/
│   │   └── config.go            # Конфигурация (читается из .env)
│   ├── entity/
│   │   └── entity.go            # Модели данных (Product, User, ProductLogs)
│   ├── handler/
│   │   └── handler.go           # HTTP-хэндлеры (контроллеры)
│   ├── middleware/
│   │   └── middleware.go        # JWT-аутентификация и CORS
│   ├── repository/
│   │   └── repository.go        # Работа с базой данных
│   ├── usecase/
│   │   └── usecase.go           # Бизнес-логика
│   └── document/
│       └── export_csv.go        # Экспорт истории в CSV
├── migrations/
│   ├── migrations.go            # Запуск миграций через goose
│   ├── 20251112143926_create_items_table.sql
│   ├── 20251113162738_user_tables.sql
│   ├── 20251113171659_product_logs.sql
│   └── 20251113171846_product_change_trigger.sql
├── web/
│   ├── index.html               # Фронтенд
│   ├── css/
│   └── js/
├── .env                         # Переменные окружения
├── Dockerfile
└── docker-compose.yml
```

---

## ⚙️ Конфигурация (`.env`)

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
TIMEZONE=Europe/Moscow

MAX_OPEN_CONNS=10
MAX_IDLE_CONNS=5
CONN_MAX_LIFETIME=30s

MAX_RETRIES=5
RETRY_DELAY=5s

JWT_SECRET=mister_popa

PACKAGE_WITH_MIGRATIONS=./migrations
```

---

## 🐳 Запуск через Docker Compose (рекомендуется)

```bash
docker-compose up --build
```

Сервис будет доступен по адресу: **http://localhost:8081**

---

## 🛠️ Запуск локально (без Docker)

1. Убедитесь, что PostgreSQL запущен и доступен.
2. Измените `DB_HOST=localhost` в файле `.env`.
3. Запустите сервис:

```bash
go run ./cmd/warehouse_control
```

---

## 🔐 Система ролей

Защищённые маршруты (`/api/v1/...`) требуют заголовок:

```
Authorization: Bearer <jwt_token>
```

| Роль | Права |
|---|---|
| `admin` | Полный доступ (GET, POST, PUT, DELETE) |
| `manager` | Просмотр и изменение (GET, PUT) |
| `viewer` | Только просмотр (GET) |

---

## 📡 API Endpoints

### 👤 Пользователи (публичные)

---

#### `POST /user/create` — Регистрация пользователя

**Тело запроса:**
```json
{
  "username": "john_doe",
  "role": "admin"
}
```

**Успешный ответ `200 OK`:**
```json
{
  "username": "john_doe"
}
```

**Ошибка `500`:**
```json
{
  "error": "user john_doe already exists"
}
```

---

#### `POST /user/login` — Вход и получение JWT-токена

**Тело запроса:**
```json
{
  "username": "john_doe",
  "role": "admin"
}
```

**Успешный ответ `200 OK`:**
```json
{
  "jwt_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибка `500`:**
```json
{
  "error": "user not found"
}
```

---

### 📦 Товары (требуют JWT)

> Все запросы ниже требуют заголовок `Authorization: Bearer <token>`

---

#### `POST /api/v1/items` — Создать товар

> 🔒 Только роль: `admin`

**Тело запроса:**
```json
{
  "name": "Ноутбук",
  "description": "Игровой ноутбук 16 дюймов",
  "quantity": 10
}
```

**Успешный ответ `200 OK`:**
```json
{
  "product_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

#### `GET /api/v1/items/:product_name` — Получить товар по названию

> 🔒 Роли: `admin`, `manager`, `viewer`

**URL параметр:** `product_name` — название товара

**Успешный ответ `200 OK`:**
```json
{
  "product": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Ноутбук",
    "description": "Игровой ноутбук 16 дюймов",
    "quantity": 10,
    "updated_at": "2025-11-13T12:00:00Z"
  }
}
```

---

#### `GET /api/v1/items` — Получить все товары

> 🔒 Роли: `admin`, `manager`, `viewer`

**Успешный ответ `200 OK`:**
```json
{
  "products": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Ноутбук",
      "description": "Игровой ноутбук 16 дюймов",
      "quantity": 10,
      "updated_at": "2025-11-13T12:00:00Z"
    }
  ]
}
```

---

#### `PUT /api/v1/items` — Обновить товар

> 🔒 Роли: `admin`, `manager`

**Тело запроса:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Ноутбук Pro",
  "description": "Обновлённая модель",
  "quantity": 5
}
```

**Успешный ответ `200 OK`:**
```json
{
  "product_updated": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Ноутбук Pro",
    "description": "Обновлённая модель",
    "quantity": 5,
    "updated_at": "2025-11-14T10:30:00Z"
  }
}
```

---

#### `DELETE /api/v1/items/:product_name` — Удалить товар

> 🔒 Только роль: `admin`

**URL параметр:** `product_name` — название товара

**Успешный ответ `200 OK`:**
```json
{
  "product_deleted": "Ноутбук Pro"
}
```

---

### 📋 История изменений

---

#### `GET /api/v1/product/logs/:product_id` — История изменений товара

> 🔒 Роли: `admin`, `manager`, `viewer`

**URL параметр:** `product_id` — UUID товара

**Успешный ответ `200 OK`:**
```json
{
  "logs": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440000",
      "old_name": "Ноутбук",
      "new_name": "Ноутбук Pro",
      "old_description": "Игровой ноутбук",
      "new_description": "Обновлённая модель",
      "old_quantity": 10,
      "new_quantity": 5,
      "changed_at": "2025-11-14T10:30:00Z"
    }
  ]
}
```

---

#### `GET /api/v1/product/save/:product_id` — Экспорт истории в CSV

> 🔒 Роли: `admin`, `manager`, `viewer`

**URL параметр:** `product_id` — UUID товара

**Ответ:** скачиваемый файл `<product_id>_<дата>.csv`

```
Content-Disposition: attachment; filename=550e8400_2025-11-14_10-30-05.csv
Content-Type: text/csv
```

---

## 🗄️ База данных

Миграции применяются **автоматически при старте** сервиса через [goose](https://github.com/pressly/goose).

| Таблица | Описание |
|---|---|
| `items` | Товары на складе |
| `users` | Пользователи системы |
| `product_logs` | История изменений товаров |

Триггер `product_change_trigger` автоматически записывает историю изменений в таблицу `product_logs` при каждом обновлении товара.
