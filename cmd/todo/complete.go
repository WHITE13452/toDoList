package main

import (
	"strconv"

	"github.com/WHITE13452/toDoList/internal/cli"
	"github.com/spf13/cobra"
)

var uncomplete bool

var completeCmd = &cobra.Command{
	Use:   "complete [task_id]",
	Short: "标记任务完成/未完成",
	Long:  "标记任务为已完成或未完成（使用 -u 参数）。",
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

		if uncomplete {
			task.MarkPending()
		} else {
			task.MarkCompleted()
		}

		if err := store.UpdateTask(task); err != nil {
			cli.PrintError("更新任务失败: %v", err)
			return
		}

		if uncomplete {
			cli.PrintSuccess("任务 %d 已标记为未完成", taskID)
		} else {
			cli.PrintSuccess("任务 %d 已完成", taskID)
		}

		cli.PrintTask(task, false)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	completeCmd.Flags().BoolVarP(&uncomplete, "uncomplete", "u", false, "标记为未完成")
}
