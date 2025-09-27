package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"godo/internal/models"
)

// helper to write arbitrary tasks JSON to a path
func writeTasksJSON(t *testing.T, path string, tasks []*models.Task) {
	t.Helper()
	data, err := json.Marshal(tasks)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func TestNewTaskStorage_usesHomeDirAndCreatesDir(t *testing.T) {
	// On Windows, os.UserHomeDir prefers USERPROFILE over HOME.
	// Set both to ensure NewTaskStorage resolves to our temp home.
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)
	ts := NewTaskStorage()
	if ts.GetFilePath() == "" {
		t.Fatalf("expected non-empty file path")
	}
	// directory should exist
	dir := filepath.Dir(ts.GetFilePath())
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("expected dir to exist: %v", err)
	}
}

func TestLoadTasks_whenFileMissingReturnsEmptySlice(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)
	ts := NewTaskStorage()
	// Ensure file does not exist
	if _, err := os.Stat(ts.GetFilePath()); !os.IsNotExist(err) {
		t.Fatalf("expected file to not exist initially")
	}
	tasks, err := ts.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks error: %v", err)
	}
	if tasks == nil || len(tasks) != 0 {
		t.Fatalf("expected empty slice, got %#v", tasks)
	}
}

func TestSaveThenLoadTasks_roundTrip(t *testing.T) {
	// isolate HOME so storage writes to temp
	temp := t.TempDir()
	t.Setenv("HOME", temp)
	t.Setenv("USERPROFILE", temp)
	ts := NewTaskStorage()

	tasks := []*models.Task{
		models.NewTask(1, "alpha"),
		models.NewTask(2, "beta"),
	}
	if err := ts.SaveTasks(tasks); err != nil {
		t.Fatalf("SaveTasks error: %v", err)
	}

	loaded, err := ts.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks error: %v", err)
	}
	if len(loaded) != len(tasks) {
		t.Fatalf("length mismatch: got %d want %d", len(loaded), len(tasks))
	}
	for i := range tasks {
		if loaded[i].ID != tasks[i].ID || loaded[i].Title != tasks[i].Title || loaded[i].Completed != tasks[i].Completed {
			t.Fatalf("task %d mismatch: got %+v want %+v", i, loaded[i], tasks[i])
		}
	}

	// Ensure file exists at expected path
	if _, err := os.Stat(ts.GetFilePath()); err != nil {
		t.Fatalf("expected file at %s: %v", ts.GetFilePath(), err)
	}
}


