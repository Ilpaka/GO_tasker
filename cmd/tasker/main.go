package main

import (
	"os"
	"path/filepath"

	"go_tasker/internal/app"
	"go_tasker/internal/storage"
)

func dataPath() string {
	if p := os.Getenv("GO_TASKER_DATA"); p != "" {
		return p
	}
	return filepath.Join("data", "tasks.json")
}

func main() {
	store := storage.NewJSONStore(dataPath())
	application := app.New(store)

	var argv []string
	if len(os.Args) > 1 {
		argv = os.Args[1:]
	}
	code := application.Run(argv)
	os.Exit(code)
}
