package main

import (
	"fmt"
	"os"

	"github.com/WHITE13452/toDoList/internal/storage"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	dbPath  string
	store   *storage.Storage
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "📋 TodoList - 智能待办事项管理工具",
	Long: `TodoList 是一个功能强大的命令行待办事项管理工具，集成了 AI Agent 智能助手。

支持传统 CLI 命令和 AI Agent 交互两种模式。`,
	Version: "1.0.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 加载 .env 文件
		_ = godotenv.Load()

		// 初始化存储
		var err error
		store, err = storage.New(dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize storage: %v\n", err)
			os.Exit(1)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// 关闭存储
		if store != nil {
			store.Close()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "", "数据库文件路径 (默认: ~/.todolist.db)")
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
