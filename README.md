# Сервис количественных исследований

Готовый MVP сервиса опросов в стиле Чистой Архитектуры для учебного проекта VK.

## Что реализовано

- Личный кабинет внутренних пользователей:
  - создание, редактирование, просмотр, удаление опросов;
  - получение результатов;
  - выгрузка результатов в `.xls`.
- Публичное прохождение опросов по уникальной ссылке:
  - получение опроса по токену;
  - старт сессии;
  - сохранение незавершенного прогресса;
  - отправка финальных ответов.
- Внешний API для интеграций:
  - получение результатов опроса по токену интеграции.
- Демо фронтенд Mini App:
  - страница прохождения опроса с сохранением прогресса и отправкой ответов.
- Безопасность в MVP:
  - непредсказуемые публичные токены (`crypto/rand`);
  - CSRF-проверка для личного кабинета;
  - экранирование текстовых ответов при показе/выгрузке;
  - базовый внутренний контроль доступа через заголовок;
  - защитные HTTP-заголовки (CSP, `X-Frame-Options`, `nosniff`).
- Логирование во всех слоях через `pkg/logger`.

## Запуск

```bash
go run ./cmd/api
```

Демо-интерфейс мини-приложения:

- `http://localhost:8080/miniapp`

Переменные окружения:

- `APP_PORT` (по умолчанию `8080`)
- `LOG_LEVEL` (`debug|info|warn|error`)
- `CSRF_TOKEN` (по умолчанию `dev-csrf-token`)
- `EXTERNAL_API_TOKEN` (по умолчанию `dev-external-token`)
- `PUBLIC_BASE_URL` (в проде укажите публичный https URL приложения, например `https://your-app.example.com`)
- `VK_APP_ID` (числовой идентификатор приложения, например `54495216`)
- `VK_SERVICE_TOKEN`
- `VK_SECURE_KEY`

## Примеры API

### 1. Создать опрос (кабинет)

```bash
curl -X POST http://localhost:8080/api/cabinet/surveys \
  -H 'Content-Type: application/json' \
  -H 'X-Internal-User: analyst' \
  -H 'X-CSRF-Token: dev-csrf-token' \
  -d '{
    "title":"Опрос по новому интерфейсу",
    "description":"Проверяем гипотезы UX",
    "questions":[
      {
        "id":"q1",
        "title":"Нравится ли интерфейс?",
        "type":"single_choice",
        "options":[{"id":"yes","text":"Да"},{"id":"no","text":"Нет"}]
      },
      {
        "id":"q2",
        "title":"Комментарий",
        "type":"free_text",
        "options":[]
      }
    ]
  }'
```

### 2. Получить публичный опрос

```bash
curl http://localhost:8080/api/public/surveys/<public_token>
```

### 3. Запустить сессию

```bash
curl -X POST http://localhost:8080/api/public/sessions \
  -H 'Content-Type: application/json' \
  -d '{"public_token":"<public_token>"}'
```

### 4. Сохранить прогресс

```bash
curl -X PUT http://localhost:8080/api/public/sessions/progress \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id":"<session_id>",
    "answers":{"q1":["yes"],"q2":["Пока думаю"]}
  }'
```

### 5. Завершить опрос

```bash
curl -X POST http://localhost:8080/api/public/sessions/submit \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id":"<session_id>",
    "answers":{"q1":["yes"],"q2":["<script>alert(1)</script>"]}
  }'
```

### 6. Получить результаты (кабинет)

```bash
curl http://localhost:8080/api/cabinet/surveys/<survey_id>/results \
  -H 'X-Internal-User: analyst'
```

### 7. Экспортировать xls

```bash
curl http://localhost:8080/api/cabinet/surveys/<survey_id>/export \
  -H 'X-Internal-User: analyst' \
  -o survey_results.xls
```

### 8. Получить результаты через внешний API

```bash
curl http://localhost:8080/api/external/surveys/<survey_id>/results \
  -H 'X-External-Token: dev-external-token'
```

## Тесты

```bash
go test ./...
```

## Деплой (Docker)

Сборка:

```bash
docker build -t research-flow:latest .
```

Запуск:

```bash
docker run --rm -p 8080:8080 \
  -e APP_PORT=8080 \
  -e LOG_LEVEL=info \
  -e CSRF_TOKEN=your-csrf-token \
  -e EXTERNAL_API_TOKEN=your-external-token \
  -e PUBLIC_BASE_URL=https://your-public-domain \
  -e VK_APP_ID=54495216 \
  -e VK_SERVICE_TOKEN=... \
  -e VK_SECURE_KEY=... \
  research-flow:latest
```

Проверка:

```bash
curl -i https://your-public-domain/health
```

URL для настройки в VK:

- `https://your-public-domain/miniapp`

## Важно для продакшна

В учебном MVP использованы in-memory адаптеры и простой контроль доступа заголовками.
Для финальной сдачи под требования VK нужно добавить:

- интеграцию auth.vk.team (OAuth2/SAML2);
- хранение секретов в Vault;
- storage-адаптеры PostgreSQL/Redis;
- интеграцию Security Gate и устранение находок;
- деплой Mini App + демозапись.
