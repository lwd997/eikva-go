FROM golang:1.23.0

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main .

EXPOSE 3000
EXPOSE 3001

RUN mv .env.example .env

CMD ["./main"]
