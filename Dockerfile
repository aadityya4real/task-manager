# ---------- Builder ----------
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o task-manager .

# ---------- Production ----------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/task-manager .
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080

CMD ["./task-manager"]