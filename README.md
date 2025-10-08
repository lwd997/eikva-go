# Приложение Тесткрафт (Быстрый старт)

## Переменные среды

Список переменных среды и их описание можно посмотреть в `.env.example`

## Сборка и запуск

### Докер

```bash
# Собираем образ
docker build -t testcraft:latest .

# Запускаем образ
 docker run  \
    -p 3000:3000 \
    -e NO_SSL_VERIFY=1 \
    -e JWT_SECRET=foo \
    -e OPEN_AI_API_KEY=foo \
    -e OPEN_AI_BASE_URL=http://foo.bar \
    -e OPEN_AI_COMPLETIONS_PATHNAME=/bar/foo \
    -e LLM_TOKEN_TRESHOLD=1000 \
    --name testcraft \
   testcraft:latest
```

### Напрямую

```bash
# Сборка клиента (нужно иметь nodejs)
cd eikva-client
npm i --omit=dev
npm run build


# Сборка сервера (нужно иметь golang)
go mod downlod
RUN go build -o eikva_testcarft .

# Запуск
./eikva_testcarft
```


# Eikva / AIQA (ТЕСТКРАФТ) — платформа управления и генерации тест‑кейсов

Единый репозиторий (монорепо) для:

- **eikva-go** — бэкенд на Go (Gin + SQLite + JWT + WebSocket) с интеграцией с OpenAI для автогенерации тест‑кейсов.
- **eikva-client** — веб‑клиент на React + TypeScript + Vite для создания, хранения и редактирования тест‑кейсов и групп.

> Проект сочетает ручное управление тест‑артефактами и автоматическую генерацию с помощью LLM.

---

## Возможности

### Клиент (AIQA)
- Регистрация и вход пользователя (JWT).
- Управление группами тест‑кейсов: создание, переименование, удаление.
- Управление тест‑кейсами: добавление, редактирование (название, предусловие, описание, постусловие).
- Просмотр и редактирование шагов тест‑кейса.
- Удобные инструменты редактирования и синхронизации с сервером (debounce ~0.5s).

### Бэкенд (Eikva)
- CRUD для групп, тест‑кейсов и их шагов.
- Генерация тест‑кейсов с помощью LLM (OpenAI **qwen3:latest**).
- WebSocket‑уведомления о завершении генерации.
- Аутентификация через JWT (access + refresh, HS512).
- SQLite в режиме WAL через `sqlx`.

---

## Технологии

- **Frontend:** React, TypeScript, Vite, ESLint, Biome.
- **Backend:** Go 1.23, Gin, Gorilla WebSocket, sqlx (SQLite), JWT (HS512), OpenAI API.

---

## Структура репозитория

```
./
├── eikva-client/              # Клиент (React + TS + Vite)
│   ├── public/                # Статика (иконки, шрифты)
│   └── src/                   # Исходники
│       ├── components/        # Переиспользуемые компоненты
│       ├── hooks/             # Кастомные React‑хуки
│       ├── http/              # HTTP‑клиент для API
│       ├── models/            # Интерфейсы TypeScript
│       ├── pages/             # Страницы (Login, Main, Group, ...)
│       ├── styles/            # Стили
│       └── main.tsx           # Точка входа
│   ├── vite.config.ts
│   ├── tsconfig*.json
│   ├── eslint.config.js
│   ├── biome.json
│   └── package.json
│
└── eikva-go/                  # Бэкенд (Go + Gin + SQLite + OpenAI)
    ├── ai/                    # Генерация тест‑кейсов (ai.go, prompts)
    ├── controllers/           # Обработчики HTTP
    ├── database/              # Миграции и доступ к БД
    ├── models/                # Модели/структуры
    ├── routes/                # Инициализация маршрутов
    ├── middlewares/           # JWT, recovery и пр.
    ├── session/               # Логика JWT‑сессий
    ├── ws/                    # WebSocket сервис
    ├── env_vars/              # Переменные окружения
    ├── requests/              # HTTP‑клиент к OpenAI
    ├── tools/                 # Утилиты и валидация
    ├── sample.json            # Пример ответа от AI
    └── main.go                # Точка входа
```

---

## Архитектура (упрощённо)

```
[AIQA (React, 5173)]  ->  [Eikva API (Go, 3000)]  ->  [SQLite (WAL)]
                                |
                                └-> [OpenAI API / qwen3:latest]

WebSocket уведомления на ws://localhost:3000/ws
Тестовая страница WebSocket: http://localhost:3000/static/index.html
```

> В dev‑режиме фронтенд проксирует API на `localhost:3000`. При запуске бэкенда также поднимается фейк‑сервер OpenAI на `:3001` (для локальной разработки).

---

## Требования

- Node.js ≥ 16 (рекомендуется LTS)
- npm или yarn
- Go ≥ 1.23
- SQLite3
- Доступ к модели LLM по OpenAI API (или локальный совместимый сервер)

---

## Быстрый старт (локально)

### 1) Клонирование

```bash
git clone <URL‑репозитория>
cd <папка‑проекта>
```

### 2) Настройка бэкенда

Создайте файл `.env` в каталоге `eikva-go` со значениями:

```dotenv
JWT_SECRET=<секрет для JWT>
OPEN_AI_API_KEY=<ваш OpenAI API ключ>
OPEN_AI_BASE_URL=https://api.openai.com
OPEN_AI_COMPLETIONS_PATHNAME=/v1/chat/completions
LLM_TOKEN_TRESHOLD=<int>
```

> **Важно:** ранние варианты документации могли указывать путь `/v1/completions`. В текущей версии используется чат‑эндпоинт `/v1/chat/completions`.

Запустите бэкенд:

```bash
cd eikva-go
go mod download
go run main.go
```

- HTTP API поднимется на `http://localhost:3000`
- Фейк OpenAI (для dev) — на `http://localhost:3001`
- База SQLite будет создана автоматически (режим WAL)

### 3) Запуск фронтенда

```bash
cd ../eikva-client
npm install
npm run dev
```

- Откроется `http://localhost:5173` (HMR включён)
- Прокси настроен на API `localhost:3000`

### 4) Предпросмотр/сборка клиента

```bash
npm run build     # сборка production
npm run preview   # локальный предпросмотр собранной версии
```

---

## Использование

1. **Регистрация и вход**
   - Перейдите на страницу входа, выполните «Регистрация», затем войдите под созданными учётными данными.
   - Токены (access/refresh) сохраняются в `localStorage`.

2. **Группы тест‑кейсов**
   - В боковой панели создайте новую группу, при наведении доступны: переименование (R), удаление (D), сохранение (S).

3. **Тест‑кейсы и шаги**
   - Внутри группы: добавляйте тест‑кейсы, редактируйте поля (название, предусловие, описание, постусловие), просматривайте и редактируйте шаги.
   - Изменения синхронизируются с сервером с задержкой ~0.5 сек (debounce).

4. **Автогенерация тест‑кейсов (AI)**
   - На бэкенде доступна операция генерации для группы: см. `POST /test-cases/start-generation`.
   - О прогрессе/завершении генерации можно получать уведомления через WebSocket.

5. **WebSocket уведомления**
   - Подключение: `ws://localhost:3000/ws` (если не переопределено в коде).
   - Тестовая страница: `http://localhost:3000/static/index.html`.

---

## REST API (кратко)

### Аутентификация
- `POST /auth/register` — регистрация пользователя
- `POST /auth/login` — вход, ответ: `{ access_token, refresh_token }`
- `POST /auth/update-tokens` — обновление токенов
- (защищённые)
  - `POST /auth/logout`
  - `GET  /auth/whoami`

### Группы тест‑кейсов
- `GET  /groups/get` — список групп
- `POST /groups/add` — создать `{ name }`
- `POST /groups/delete` — удалить `{ uuid }`
- `POST /groups/rename` — переименовать `{ uuid, name }`
- `GET  /groups/get-test-cases/:groupUUID` — тест‑кейсы группы

### Тест‑кейсы
- `POST /test-cases/add` — создать пустой `{ test_case_group }`
- `POST /test-cases/start-generation` — сгенерировать `{ test_case_group, amount, user_input }`
- `POST /test-cases/update` — обновить поля тест‑кейса
- `POST /test-cases/delete` — удалить `{ uuid }`
- `GET  /test-cases/get-steps/:testCaseUUID` — получить шаги тест‑кейса

### Шаги тест‑кейса
- `POST /steps/add` — создать пустой шаг `{ test_case }`
- `POST /steps/update` — обновить `{ uuid, description, data, expected_result }`
- `POST /steps/delete` — удалить `{ uuid }`
- `POST /steps/swap` — поменять местами `{ first, second }`

---

## Примеры запросов (cURL)

**Регистрация и логин**
```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"demo","password":"demo123"}'

curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"demo","password":"demo123"}'
# => { "access_token": "...", "refresh_token": "..." }
```

**Создание группы**
```bash
curl -X POST http://localhost:3000/groups/add \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Моя первая группа"}'
```

**Старт генерации тест‑кейсов**
```bash
curl -X POST http://localhost:3000/test-cases/start-generation \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
        "test_case_group":"<GROUP_UUID>",
        "amount": 5,
        "user_input": "Сократить требования и сгенерировать базовые позитивные и негативные сценарии"
      }'
```

---

## Развёртывание

1. Соберите клиент: `npm run build`; отдавайте содержимое `dist/` через любой статический веб‑сервер / CDN.
2. Соберите и запустите бэкенд (пример):
   ```bash
   cd eikva-go
   go build -o eikva
   ./eikva
   ```
3. Настройте переменные окружения (`.env`) и обратите внимание на секреты (JWT, OpenAI ключи).
4. За проксирование/SSL может отвечать внешний веб‑сервер (nginx, Caddy и т.д.).

---

## Контакты / поддержка

- Вопросы по коду, архитектуре, ML-части и интеграции: Раков Леонид (ИЛ Новосибирск), Гущин Сергей (ИЛ Новосибирск), Юферев Виталий (ИЛ Новосибирск)

