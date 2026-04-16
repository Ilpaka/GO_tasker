package app

import (
	"fmt"
	"os"
	"strings"

	"go_tasker/internal/storage"
	"go_tasker/internal/task"
)

// App связывает хранилище и вывод в консоль.
type App struct {
	store *storage.JSONStore
	out   *os.File
}

func New(store *storage.JSONStore) *App {
	return &App{store: store, out: os.Stdout}
}

func (a *App) println(msg string) {
	_, _ = fmt.Fprintln(a.out, msg)
}

func (a *App) Usage() {
	a.println(`GO_tasker — простой трекер задач (CLI).

Использование:
  tasker add <заголовок> [описание]   добавить задачу
  tasker list                       показать все задачи
  tasker done <id>                  отметить задачу выполненной
  tasker remove <id>                удалить задачу
  tasker help                       эта справка

Файл данных по умолчанию: data/tasks.json
Переопределение: переменная окружения GO_TASKER_DATA`)
}

func (a *App) Add(args []string) int {
	if len(args) == 0 {
		a.println("ошибка: укажите заголовок задачи, например: tasker add \"Купить хлеб\"")
		return 1
	}
	title := strings.TrimSpace(args[0])
	desc := ""
	if len(args) > 1 {
		desc = strings.TrimSpace(strings.Join(args[1:], " "))
	}
	if title == "" {
		a.println("ошибка: заголовок не может быть пустым")
		return 1
	}
	t, err := a.store.Add(title, desc)
	if err != nil {
		a.println("ошибка сохранения: " + err.Error())
		return 1
	}
	a.println(fmt.Sprintf("добавлена задача #%d: %s", t.ID, t.Title))
	return 0
}

func (a *App) List() int {
	tasks, err := a.store.List()
	if err != nil {
		a.println("ошибка чтения: " + err.Error())
		return 1
	}
	if len(tasks) == 0 {
		a.println("задач пока нет — добавьте первую командой tasker add")
		return 0
	}
	for _, t := range tasks {
		a.println(formatTask(t))
	}
	return 0
}

func formatTask(t task.Task) string {
	status := "активна"
	if t.Done {
		status = "выполнена"
	}
	line := fmt.Sprintf("#%d [%s] %s", t.ID, status, t.Title)
	if strings.TrimSpace(t.Description) != "" {
		line += "\n    " + t.Description
	}
	return line
}

func (a *App) Done(idStr string) int {
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id <= 0 {
		a.println("ошибка: укажите положительный числовой id, например: tasker done 2")
		return 1
	}
	if err := a.store.SetDone(id); err != nil {
		if err == storage.ErrNotFound {
			a.println(fmt.Sprintf("задача #%d не найдена", id))
		} else {
			a.println("ошибка: " + err.Error())
		}
		return 1
	}
	a.println(fmt.Sprintf("задача #%d отмечена как выполненная", id))
	return 0
}

func (a *App) Remove(idStr string) int {
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id <= 0 {
		a.println("ошибка: укажите положительный числовой id, например: tasker remove 2")
		return 1
	}
	if err := a.store.Remove(id); err != nil {
		if err == storage.ErrNotFound {
			a.println(fmt.Sprintf("задача #%d не найдена", id))
		} else {
			a.println("ошибка: " + err.Error())
		}
		return 1
	}
	a.println(fmt.Sprintf("задача #%d удалена", id))
	return 0
}

func (a *App) Run(argv []string) int {
	if len(argv) < 1 {
		a.Usage()
		return 1
	}
	cmd := argv[0]
	args := argv[1:]
	switch cmd {
	case "help", "-h", "--help":
		a.Usage()
		return 0
	case "add":
		return a.Add(args)
	case "list":
		return a.List()
	case "done":
		if len(args) < 1 {
			a.println("ошибка: укажите id задачи")
			return 1
		}
		return a.Done(args[0])
	case "remove":
		if len(args) < 1 {
			a.println("ошибка: укажите id задачи")
			return 1
		}
		return a.Remove(args[0])
	default:
		a.println("неизвестная команда: " + cmd)
		a.Usage()
		return 1
	}
}
