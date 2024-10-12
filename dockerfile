# Stage 1: Install dependencies
FROM golang:1.23.1-bookworm AS deps

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Stage 2: Build the application
FROM golang:1.23.1-bookworm AS builder

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY . .

RUN go build -ldflags="-w -s" -o main .

# Final stage: Run the application
FROM debian:bookworm-slim

WORKDIR /app

# Copy the built application
COPY --from=builder /app/main .

CMD ["./main"]
