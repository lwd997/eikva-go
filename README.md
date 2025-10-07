# Приложение Тесткрафт

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
