package models

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	before := time.Now().Add(-time.Second)
	task := NewTask(1, "test")
	if task.ID != 1 {
		t.Fatalf("expected ID 1, got %d", task.ID)
	}
	if task.Title != "test" {
		t.Fatalf("expected Title 'test', got %q", task.Title)
	}
	if task.Completed {
		t.Fatalf("expected Completed false, got true")
	}
	if task.CreatedAt.Before(before) || task.UpdatedAt.Before(before) {
		t.Fatalf("timestamps not set correctly")
	}
}

func TestTaskManager_AddGetStats(t *testing.T) {
	m := NewTaskManager([]*Task{})
	m.AddTask("a")
	m.AddTask("b")
	if completed, total := m.GetStats(); completed != 0 || total != 2 {
		t.Fatalf("expected (0,2), got (%d,%d)", completed, total)
	}
	// toggle first
	if !m.ToggleTask(0) {
		t.Fatalf("toggle should succeed")
	}
	if completed, total := m.GetStats(); completed != 1 || total != 2 {
		t.Fatalf("expected (1,2), got (%d,%d)", completed, total)
	}
}

func TestTaskManager_DeleteUpdateBounds(t *testing.T) {
	m := NewTaskManager([]*Task{})
	m.AddTask("a")
	m.AddTask("b")
	if !m.UpdateTask(1, "bb") {
		t.Fatalf("update should succeed")
	}
	if m.GetTaskByIndex(1).Title != "bb" {
		t.Fatalf("title not updated")
	}
	if !m.DeleteTask(0) {
		t.Fatalf("delete should succeed")
	}
	if m.DeleteTask(10) {
		t.Fatalf("out-of-range delete should fail")
	}
	if m.ToggleTask(-1) {
		t.Fatalf("negative index toggle should fail")
	}
}

func TestNewTaskManagerWithExistingTasksSetsNextID(t *testing.T) {
	// existing tasks with max ID 5
	existing := []*Task{
		{ID: 2, Title: "x"},
		{ID: 5, Title: "y"},
	}
	m := NewTaskManager(existing)
	m.AddTask("z")
	lastTask := m.GetTaskByIndex(len(m.GetTasks()) - 1)
	if lastTask.ID != 6 {
		// nextID should continue from max(existing)+1
		t.Fatalf("expected new task ID 6, got %d", lastTask.ID)
	}
}


