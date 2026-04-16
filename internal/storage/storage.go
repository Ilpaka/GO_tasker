package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go_tasker/internal/task"
)

var ErrNotFound = errors.New("задача не найдена")

// JSONStore хранит задачи в одном JSON-файле.
type JSONStore struct {
	path string
	mu   sync.Mutex
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{path: path}
}

func (s *JSONStore) filePath() string {
	return s.path
}

func (s *JSONStore) load() ([]task.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath())
	if err != nil {
		if os.IsNotExist(err) {
			return []task.Task{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []task.Task{}, nil
	}
	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("разбор JSON: %w", err)
	}
	return tasks, nil
}

func (s *JSONStore) save(tasks []task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.filePath()), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.filePath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.filePath())
}

// List возвращает копию всех задач.
func (s *JSONStore) List() ([]task.Task, error) {
	return s.load()
}

// Add добавляет задачу и возвращает её с назначенным ID.
func (s *JSONStore) Add(title, description string) (task.Task, error) {
	tasks, err := s.load()
	if err != nil {
		return task.Task{}, err
	}
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newTask := task.Task{
		ID:          maxID + 1,
		Title:       title,
		Description: description,
		Done:        false,
	}
	tasks = append(tasks, newTask)
	if err := s.save(tasks); err != nil {
		return task.Task{}, err
	}
	return newTask, nil
}

// SetDone помечает задачу выполненной.
func (s *JSONStore) SetDone(id int) error {
	tasks, err := s.load()
	if err != nil {
		return err
	}
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			return s.save(tasks)
		}
	}
	return ErrNotFound
}

// Remove удаляет задачу по ID.
func (s *JSONStore) Remove(id int) error {
	tasks, err := s.load()
	if err != nil {
		return err
	}
	var out []task.Task
	found := false
	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		out = append(out, t)
	}
	if !found {
		return ErrNotFound
	}
	return s.save(out)
}
