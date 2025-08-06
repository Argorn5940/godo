package storage

import (
	"encoding/json"
	"fmt"
	"godo/internal/models"
	"os"
	"path/filepath"
)

const (
	//タスクファイル名
	TasksFileName = "tasks.json"
)

//TaskStorage はタスクの保存・読み込みを管理する
type TaskStorage struct {
	filePath string
}

//NewTaskStorageは新しいTaskStorageを作成する
func NewTaskStorage() *TaskStorage {
	//ホームディレクトリ配下にファイルを作成
	homeDir, err := os.UserHomeDir()
	if err != nil {
		//エラーの場合は現在のディレクトリを使用
		homeDir = "."
	}

	//ファイルパスを作成
	filePath := filepath.Join(homeDir, ".godo", TasksFileName)

	// .godoディレクトリが存在しない場合は作成
	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, 0755)

	return &TaskStorage{
		filePath: filePath,
	}
}

// LoadTasksはファイルからタスクを読み込む
func (ts *TaskStorage) LoadTasks() ([]*models.Task, error) {
	//ファイルが存在しない場合は空のスライスを返す
	if _, err := os.Stat(ts.filePath); os.IsNotExist(err) {
		return []*models.Task{}, nil
	}

	//ファイルを読み込む
	data, err := os.ReadFile(ts.filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}

	//JSONをパースする
	var tasks []*models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("JSONのパースに失敗しました: %w", err)
	}

	return tasks, nil

}

// SaveTasksはタスクはタスクをファイルに保存する
func (ts *TaskStorage) SaveTasks(tasks []*models.Task) error {
	//JSONに変換（見やすくインデント付き）
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONの作成に失敗しました: %w", err)
	}

	//ファイルに保存
	if err := os.WriteFile(ts.filePath, data, 0644); err != nil {
		return fmt.Errorf("ファイルへの保存に失敗しました: %w", err)
	}

	return nil
}

//GetFilePathは保存先のファイルパスを返す（デバック用）
func (ts *TaskStorage) GetFilePath() string {
	return ts.filePath
}