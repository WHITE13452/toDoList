package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/WHITE13452/toDoList/internal/cli"
	"github.com/spf13/cobra"
)

var skipConfirm bool

var deleteCmd = &cobra.Command{
	Use:   "delete [task_id]",
	Short: "删除任务",
	Long:  "删除指定的待办事项。默认需要确认，使用 -y 参数跳过确认。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			cli.PrintError("无效的任务 ID")
			return
		}

		// 获取任务信息
		task, err := store.GetTask(taskID)
		if err != nil {
			cli.PrintError("获取任务失败: %v", err)
			return
		}
		if task == nil {
			cli.PrintError("任务 %d 不存在", taskID)
			return
		}

		// 确认删除
		if !skipConfirm {
			cli.PrintTask(task, false)
			fmt.Printf("\n确定要删除任务 %d 吗？(y/N): ", taskID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("已取消")
				return
			}
		}

		if err := store.DeleteTask(taskID); err != nil {
			cli.PrintError("删除任务失败: %v", err)
			return
		}

		cli.PrintSuccess("任务 %d 已删除", taskID)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "跳过确认")
}
