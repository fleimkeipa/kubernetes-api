FROM golang:1.22.3

WORKDIR /app

COPY . .

# Utilize cache and avoid re-downloading if no changes are detected
RUN go mod download
RUN go mod verify

RUN go build -o /kubernetes-api ./main.go

EXPOSE 8080

CMD ["/kubernetes-api"]
