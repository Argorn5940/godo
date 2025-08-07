/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"godo/internal/models"
	"godo/internal/storage"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godo",
	Short: "超シンプルタスク管理アプリ",
	Run: func(cmd *cobra.Command, args []string) {
		// 保存機能のテスト
		testStorage()
	},
}
func testStorage() {
	fmt.Println("📄 Godo - タスク管理")
	fmt.Println("保存機能をテストしています...")

	// ストレージを初期化
	storage := storage.NewTaskStorage()
	fmt.Printf("保存場所: %s\n", storage.GetFilePath())

	// 既存のタスクを読み込み
	tasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		return
	}

	// TaskManagerを初期化
	manager := models.NewTaskManager(tasks)

	// サンプルタスクを追加（初回のみ）
	if len(tasks) == 0 {
		manager.AddTask("Go言語を学習する")
		manager.AddTask("bubbleteaを理解する")
		manager.AddTask("Cobraを覚える")
		fmt.Println("サンプルタスクを追加しました")
	}

	// タスク一覧を表示
	fmt.Println("\n現在のタスク:")
	for _, task := range manager.GetTasks() {
		status := "○"
		if task.Completed {
			status = "✓"
		}
		fmt.Printf("%s %s (ID: %d)\n", status, task.Title, task.ID)
	}

	// 統計を表示
	completed, total := manager.GetStats()
	fmt.Printf("\n統計: 完了 %d | 未完了 %d\n", completed, total-completed)

	// ファイルに保存
	if err := storage.SaveTasks(manager.GetTasks()); err != nil {
		fmt.Printf("保存エラー: %v\n", err)
		return
	}

	fmt.Println("✅ タスクをファイルに保存しました")
} 

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}



