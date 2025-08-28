package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 新しいタスクを作成する関数
func NewTask(id int, title string) *Task {
	now := time.Now()
	return &Task{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TaskManager タスク管理を行う構造体
type TaskManager struct {
	tasks  []*Task
	nextID int
}

// NewTaskManager 新しいTaskManagerを作成する
func NewTaskManager(tasks []*Task) *TaskManager {
	nextID := 1
	if len(tasks) > 0 {
		// 既存のタスクから最大IDを取得
		for _, task := range tasks {
			if task.ID >= nextID {
				nextID = task.ID + 1
			}
		}
	}
	
	return &TaskManager{
		tasks:  tasks,
		nextID: nextID,
	}
}

// AddTask 新しいタスクを追加する
func (tm *TaskManager) AddTask(title string) {
	task := NewTask(tm.nextID, title)
	tm.tasks = append(tm.tasks, task)
	tm.nextID++
}

// GetTasks 全てのタスクを取得する
func (tm *TaskManager) GetTasks() []*Task {
	return tm.tasks
}

// DeleteTask 指定されたインデックスのタスクを削除する
func (tm *TaskManager) DeleteTask(index int) bool {
	if index < 0 || index >= len(tm.tasks) {
		return false
	}
	
	// スライスから要素を削除
	tm.tasks = append(tm.tasks[:index], tm.tasks[index+1:]...)
	return true
}

// ToggleTask 指定されたインデックスのタスクの完了状態を切り替える
func (tm *TaskManager) ToggleTask(index int) bool {
	if index < 0 || index >= len(tm.tasks) {
		return false
	}
	
	tm.tasks[index].Completed = !tm.tasks[index].Completed
	tm.tasks[index].UpdatedAt = time.Now()
	return true
}

// UpdateTask 指定されたインデックスのタスクのタイトルを更新する
func (tm *TaskManager) UpdateTask(index int, title string) bool {
	if index < 0 || index >= len(tm.tasks) {
		return false
	}
	
	tm.tasks[index].Title = title
	tm.tasks[index].UpdatedAt = time.Now()
	return true
}

// GetTaskByIndex 指定されたインデックスのタスクを取得する
func (tm *TaskManager) GetTaskByIndex(index int) *Task {
	if index < 0 || index >= len(tm.tasks) {
		return nil
	}
	return tm.tasks[index]
}

// GetStats 完了済みと未完了のタスク数を取得する
func (tm *TaskManager) GetStats() (completed, total int) {
	total = len(tm.tasks)
	for _, task := range tm.tasks {
		if task.Completed {
			completed++
		}
	}
	return completed, total
}