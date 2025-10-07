# Клиент
FROM node:20 AS eikva_client
WORKDIR /app/eikva_client
COPY eikva-client/ .
RUN npm i --omit=dev
RUN npm run build

# Сервер
FROM golang:1.23.0 AS eikva_go
WORKDIR /app/eikva_go
COPY . .
COPY --from=eikva_client app/static ./static
RUN go mod download
RUN go build -o eikva_testcarft .

# Финал
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=eikva_go /app/eikva_go/eikva_testcarft .
COPY --from=eikva_go /app/eikva_go/static ./static
ENV  GIN_MODE=release
EXPOSE 3000

CMD ["./eikva_testcarft"]
