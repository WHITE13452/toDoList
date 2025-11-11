package main

import (
	"fmt"

	"github.com/WHITE13452/toDoList/internal/cli"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "搜索任务",
	Long:  "在待办事项中搜索包含指定关键词的任务（标题或描述）。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]

		tasks, err := store.SearchTasks(keyword)
		if err != nil {
			cli.PrintError("搜索失败: %v", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Printf("未找到包含 '%s' 的任务\n", keyword)
			return
		}

		fmt.Printf("\n找到 %d 个匹配的任务:\n\n", len(tasks))
		for _, task := range tasks {
			cli.PrintTask(task, false)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
