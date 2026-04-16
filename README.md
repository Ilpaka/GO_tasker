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
├── Makefile             fmt / lint / check
├── .golangci.yml        настройки golangci-lint
├── .pre-commit-config.yaml
├── step_notes.txt       короткие заметки по доработке (для себя)
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

## Форматирование, линтер и pre-commit

Отдельного конфига для форматтера в Go не нужно: стиль задаёт встроенный **gofmt**, а запускать удобно командой **`go fmt ./...`** (она как раз прогоняет gofmt по пакетам проекта).

**Линтер:** [golangci-lint](https://golangci-lint.run/) с коротким файлом `.golangci.yml`: включены только `errcheck`, `govet`, `staticcheck`, `unused`, `ineffassign` и проверка **`gofmt`** (чтобы забытое форматирование не проскочило в репозиторий). Раздувать список линтеров не стоит — для маленького CLI хватает этого набора.

**Pre-commit:** один раз ставится сам фреймворк (`pip install pre-commit` или как у вас принято в Python-окружении), затем в корне репозитория:

```bash
cd GO_tasker
pre-commit install
```

Перед каждым `git commit` запустятся хуки из `.pre-commit-config.yaml`: обрезка хвостовых пробелов и перевод строки в конце файла, лёгкая проверка YAML, лимит на огромные файлы, **`go fmt`** по Go-файлам и **`golangci-lint`**. Если что-то падает, коммит не создаётся — правите вывод, снова `git add`, повторяете коммит.

Проверить всё вручную без коммита:

```bash
make fmt          # go fmt ./...
make lint         # golangci-lint run ./...
make check        # fmt и lint подряд
pre-commit run --all-files   # то же по смыслу, что увидит хук, по всему дереву
```

Поставить `golangci-lint`, если его ещё нет: с [страницы релизов](https://github.com/golangci/golangci-lint/releases) под свою ОС или через `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` (бинарник окажется в `$(go env GOPATH)/bin`, этот путь должен быть в `PATH`).

## Работа с Git

Репозиторий оформлен как учебный: есть **осмысленные коммиты** на русском, разделение работы по **веткам** и слияние в `main`. Достаточно обычного **`git`** и браузера на GitHub — **GitHub CLI (`gh`) не нужен**.

### Клонирование и отправка изменений (только git)

Если репозиторий уже есть на GitHub:

```bash
git clone https://github.com/Ilpaka/GO_tasker.git
cd GO_tasker
```

Свой форк или другой URL подставьте вместо ссылки выше.

Работа в новой ветке и push:

```bash
git checkout main
git pull
git checkout -b feature/моя-ветка
# правки, затем:
git add .
git commit -m "Краткое осмысленное сообщение на русском"
git push -u origin feature/моя-ветка
```

Слияние локально (без PR), если нужно у себя:

```bash
git checkout main
git merge --no-ff feature/моя-ветка -m "Слияние feature/моя-ветка"
git push origin main
```

### Ветки в этом проекте

- **`main`** — стабильная версия после слияний.
- **`feature/cli-commands`** — интерфейс командной строки и логика команд.
- **`feature/docker-support`** — контейнеризация (Dockerfile, Compose, `.dockerignore`).

Сначала в `main` появились модуль и хранилище, затем в feature-ветке — CLI, после merge — Docker-ветка, затем документация в `main`.

### Pull Request (через сайт GitHub, без gh)

После `git push` откройте репозиторий в браузере — GitHub предложит **Compare & pull request** для новой ветки. Либо: **Pull requests → New pull request**, база **`main`**, сравнение с вашей веткой, заголовок и описание на русском, затем **Create pull request** и **Merge pull request**.

Пример по смыслу задания: слияние **`feature/docker-support` → `main`**, в описании PR — какие файлы Docker добавлены и как проверяли сборку.

### GitHub Project (через сайт)

Доска с задачами (минимум 6 этапов): **Projects → New project** (шаблон **Table**), при необходимости привяжите к репозиторию. Поля в настройках проекта: **Priority**, **Start Date**, **End Date**, **Original Estimate**; статусы колонок или поля статуса: **Proposed**, **Active**, **Resolved**, **Completed**.

Готовые пункты можно завести вручную или добавить существующие [issues](https://github.com/Ilpaka/GO_tasker/issues) в проект кнопкой **Add item** на доске.

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
