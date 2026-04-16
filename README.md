# GO_tasker

Публичный репозиторий: https://github.com/Ilpaka/GO_tasker

Минимальный **CLI-трекер задач** на Go: добавление, просмотр, отметка «выполнено», удаление. Данные хранятся в **JSON-файле** между запусками.

## Цель проекта

Учебный мини-проект: показать базовую структуру Go-приложения (`cmd` / `internal`), работу с файловым хранилищем и оформление репозитория (Git, ветки, PR, GitHub Projects, Docker).

## Функционал

- добавление задачи (`add`);
- просмотр всех задач (`list`);
- отметка выполненной (`done`);
- удаление (`remove`);
- сохранение в `data/tasks.json` (или путь из `GO_TASKER_DATA`).

## Структура проекта

```
GO_tasker/
├── cmd/tasker/          точка входа (main)
├── internal/
│   ├── app/             команды CLI и вывод
│   ├── storage/         чтение/запись JSON
│   └── task/            тип «задача»
├── data/                каталог файла данных (в репозитории — только .gitkeep)
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── go.mod
├── .gitignore
└── README.md
```

## Запуск локально

Требуется Go **1.22+**.

```bash
cd GO_tasker
go run ./cmd/tasker help
```

Сборка бинарника:

```bash
go build -o tasker ./cmd/tasker
./tasker list
```

## Примеры использования

```bash
# справка
go run ./cmd/tasker help

# добавить задачу (описание можно опустить)
go run ./cmd/tasker add "Сдать отчёт" "Проверить раздел про Docker"

# список
go run ./cmd/tasker list

# отметить выполненной
go run ./cmd/tasker done 1

# удалить
go run ./cmd/tasker remove 1
```

Файл по умолчанию: `data/tasks.json`. Другой путь:

```bash
export GO_TASKER_DATA=/tmp/мои_задачи.json
go run ./cmd/tasker list
```

## Работа с Git

Репозиторий оформлен как учебный: есть **осмысленные коммиты** на русском, разделение работы по **веткам** и слияние в `main`.

### Ветки

- **`main`** — стабильная версия после слияний.
- **`feature/cli-commands`** — интерфейс командной строки и логика команд.
- **`feature/docker-support`** — контейнеризация (Dockerfile, Compose, `.dockerignore`).

Сначала в `main` появились модуль и хранилище, затем в feature-ветке — CLI, после merge — Docker-ветка, затем документация в `main`.

### Pull Request

Пример оформления (см. раздел «Текст Pull Request» в отчёте к заданию): слияние **`feature/docker-support` → `main`**, в PR описаны файлы Docker и проверка запуска.

### GitHub Project

Для проекта заведена доска с задачами (минимум 6): планирование этапов от каркаса репозитория до Docker. Поля: **Priority**, **Start Date**, **End Date**, **Original Estimate**; статусы: **Proposed**, **Active**, **Resolved**, **Completed**. Через веб-интерфейс GitHub: **Projects → New project** (Table), привязка к репозиторию **Ilpaka/GO_tasker**, импорт или ручное добавление пунктов из [списка задач](https://github.com/Ilpaka/GO_tasker/issues). Для автоматизации в CLI позже выполните `gh auth refresh -s project,read:project`.

---

## Docker и Docker Compose

Контейнеризация нужна, чтобы **одинаково** запускать утилиту на любой машине с Docker, не ставя Go, и чтобы данные жили в примонтированной папке `./data`.

### Сборка образа

```bash
docker build -t go-tasker:local .
```

Ожидаемо в конце: `Successfully tagged go-tasker:local`. Проверка:

```bash
docker image ls | grep go-tasker
```

### Запуск контейнера

Интерактивно, с привязкой данных к каталогу проекта:

```bash
docker run --rm -v "$(pwd)/data:/app/data" go-tasker:local list
docker run --rm -v "$(pwd)/data:/app/data" go-tasker:local add "Задача из контейнера"
```

Проверка: команды завершаются без ошибки, `list` показывает задачи, файл появляется на хосте в `./data/tasks.json`.

Остановка: для `--rm` контейнер удаляется после выхода; долгий процесс — `docker stop <id>`.

### Docker Compose

Запуск (актуальный синтаксис V2):

```bash
docker compose up --build
```

По умолчанию сервис выводит справку (`help`). Рабочие команды удобно запускать так:

```bash
docker compose run --rm tasker list
docker compose run --rm tasker add "Покупки" "Молоко, яйца"
docker compose run --rm tasker done 1
```

Остановка сервиса, если поднимали `up` в фоне или в отдельном терминале:

```bash
docker compose down
```

Удаление контейнеров и сети compose; тома при bind-mount `./data` данные на диске остаются.

---

## Лицензия

Учебный проект, без лицензии.
