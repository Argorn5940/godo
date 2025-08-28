/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"godo/internal/ui"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godo",
	Short: "超シンプルタスク管理アプリ",
	Long: `Godo - Go言語で作った超シンプルなタスク管理アプリ
	
bubbleteaを使ったターミナルUIで、直感的にタスク管理ができます。

操作方法:
  Enter     - タスクの完了/未完了を切り替え
  n         - 新しいタスクを追加
  e         - 選択したタスクを編集
  d         - 選択したタスクを削除
  ↑/↓ or j/k - タスクの選択を移動
  q         - アプリケーションを終了`,
	Run: func(cmd *cobra.Command, args []string) {
		// TUIアプリケーションを開始
		if err := ui.RunApp(); err != nil {
			fmt.Printf("アプリケーション実行エラー: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}