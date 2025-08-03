package storage

import (
	"encoding/json"
	"fmt"
	"godo/internal/models"
	"os"
	"path/filepath"
	"time"
)

const (
	dataFileName = "tasks.json"
	dataDir = ".godo"
)

// Storage はタスクデーターの保存・読込を管理する構造体
type Storage struct {
	filePath string
}

// New は新しいStirageインスタンスを作成する
func New() (*Storage, error) {
	//ホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
	}

	//データディレクトリのパスを作成
	dataPath := filepath.Join(homeDir, dataDir)

	//データディレクトリが存在しない場合は作成
	if  err := os.MkdirAll(dataPath, 0755); err != nil {
		return nil, fmt.Errorf("データディレクトリの作成に失敗: %w", err)
	}

	//ファイルパスを指定
	filePath := filepath.Join(dataPath, dataFileName)

	return &Storage{
		filePath: filePath,
	}, nil
}

//LoadTasks は保存されているタスクを読み込む
func (s *Storage) LoadTasks() ([]models.Task, error) {
	//ファイルが存在しない場合は空のスライスを返す
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []models.Task{}, nil
	}

	//ファイルを読み込む
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	//JSONをパース
	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("JSONのパースに失敗: %w", err)
	}

	return tasks, nil
}

// SabeTasksはタスクをファイルに保存する
func (s *Storage) SaveTasks(tasks []models.Task) error {
	//JSONに変換（読みやすい形式で）
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONの変換に失敗: %w", err)
	}
	
	//ファイルに書き込み
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("ファイルへの書き込みに失敗: %w", err)
	}

	return nil
}

// GetFilePathはデータファイルのパスを返す(デバック用)
func (s *Storage) GetFilePath() string {
	return s.filePath
}

//BackupTasksは現在のタスクデーターをバックアップする
func (s *Storage) BackupTasks() error {
	//元のファイルが存在しない場合は何もしない
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return nil
	}

	//バックアップファイル名を作成（タイムスタンプ付き）
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup_%s", s.filePath, timestamp)

	//ファイルをコピー
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("バックアップ対象ファイルの読み込みに失敗: %w", err)
	}
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("バックアップファイルの作成に失敗: %w", err)
	}

	return nil
}