FROM golang:1.22.3

WORKDIR /app

COPY . .
COPY .env /app/.env

RUN go mod tidy
RUN go build -o /kubernetes-api ./main.go

EXPOSE 8080

CMD ["/kubernetes-api"]