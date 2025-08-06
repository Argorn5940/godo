package models

import (
	"fmt"
	"time"
)

// Task構造体（既存のコードに追加
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed   bool   `json:"completed"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTaskは新しいタスクを作成する
func NewTask(id int, title string) *Task {
	now := time.Now()
	return &Task{
		ID: id,
		Title: title,
		Completed: false,
		UpdatedAt: now,
	}
}

//Toggleはタスクの完了状態を切り替える
func (t *Task) Toggle() {
	t.Completed = !t.Completed
	t.UpdatedAt = time.Now()
}

//UpdateTitleはタスクのタイトルを更新する
func (t *Task) UpdateTitle(title string) {
	if len(title) > 30 {
		title = title[:30] //30文字制限
	}
	t.Title = title
	t.UpdatedAt = time.Now()
}

//TaskManagerはタスクの集合を管理する
type TaskManager struct {
	Tasks []*Task
	nextID int
}

//NewTaskManagerは新しいTaskManagerを作成する
func NewTaskManager(tasks []*Task) *TaskManager {
	nextID := 1
	if len(tasks) > 0 {
		//最大IDを見つけて+1する
		for _, task := range tasks {
			if task.ID >= nextID {
				nextID = task.ID + 1
			}
		}
	}
	return &TaskManager{
		Tasks: tasks,
		nextID: nextID,
	}
}

//GetTasksは全てのタスクを返す
func (tm *TaskManager) GetTasks() []*Task {
	return tm.Tasks
}

//AddTaskは新しいタスクを追加する
func (tm *TaskManager) AddTask(title string){
	if title == "" {
		return
	}

	task := NewTask(tm.nextID, title)
	tm.Tasks = append(tm.Tasks, task)
	tm.nextID++
}

//RemoveTaskは指定されたIDのタスクを削除する
func (tm *TaskManager) RemoveTask(id int)  error{
	for i, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks = append(tm.Tasks[:i], tm.Tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("タスクが見つかりませんでした: %d", id)
}

//GetTAskByIDは指定されたIDのタスクを取得する
func (tm *TaskManager) GetTaskByID(id int) *Task{
	for _, task := range tm.Tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

//GetStatsは完了・未完了の統計を返す
func (tm *TaskManager) GetStats() (completed, total int) {
	total = len(tm.Tasks)
	for _, task := range tm.Tasks {
		if task.Completed {
			completed++
		}
	}
	return completed, total
}