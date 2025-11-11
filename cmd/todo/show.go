package main

import (
	"strconv"

	"github.com/WHITE13452/toDoList/internal/cli"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [task_id]",
	Short: "显示任务详情",
	Long:  "显示指定任务的详细信息。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			cli.PrintError("无效的任务 ID")
			return
		}

		task, err := store.GetTask(taskID)
		if err != nil {
			cli.PrintError("获取任务失败: %v", err)
			return
		}
		if task == nil {
			cli.PrintError("任务 %d 不存在", taskID)
			return
		}

		cli.PrintTask(task, true)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
