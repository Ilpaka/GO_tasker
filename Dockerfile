# Сборка статического бинарника
FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/tasker ./cmd/tasker

# Минимальный образ для запуска CLI
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/tasker /usr/local/bin/tasker

ENV GO_TASKER_DATA=/app/data/tasks.json

ENTRYPOINT ["tasker"]
CMD ["help"]
