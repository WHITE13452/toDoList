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
	Long:  "删除指定的待办事项。可以使用任务ID或关键词搜索。默认需要确认,使用 -y 参数跳过确认。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]

		// 尝试将输入解析为任务ID
		taskID, err := strconv.ParseInt(input, 10, 64)
		if err == nil {
			// 按 ID 删除
            deleteByID(taskID)
            return
		}

		// 按关键词搜索并删除
        deleteByKeyword(input)
	},
}

// deleteByID 根据任务ID删除任务
func deleteByID(taskID int64) {
	task, err := store.GetTask(taskID)
	if err != nil {
		cli.PrintError("获取任务失败: %v", err)
		return
	}
	if task == nil {
		cli.PrintError("任务 %d 不存在", taskID)
		return
	}

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
}


// deleteByKeyword 按关键词搜索并删除任务
func deleteByKeyword(keyword string) {
	// 搜索任务
    tasks, err := store.SearchTasks(keyword)
    if err != nil {
        cli.PrintError("搜索失败: %v", err)
        return
    }

    if len(tasks) == 0 {
        fmt.Printf("未找到包含 '%s' 的任务\n", keyword)
        return
    }

    // 显示搜索结果
    fmt.Printf("\n找到 %d 个匹配的任务:\n\n", len(tasks))
    for i, task := range tasks {
        fmt.Printf("[%d] ", i+1)
        cli.PrintTask(task, false)
    }
	 
	// 提示用户选择
    fmt.Printf("\n请选择要删除的任务 (输入序号 1-%d, 输入 0 取消): ", len(tasks))
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(input)

    // 解析选择
    choice, err := strconv.Atoi(input)
    if err != nil || choice < 0 || choice > len(tasks) {
        cli.PrintError("无效的选择")
        return
    }

    if choice == 0 {
        fmt.Println("已取消")
        return
    }

    // 选中的任务
    selectedTask := tasks[choice-1]

    // 二次确认
    if !skipConfirm {
        fmt.Println("\n您选择删除的任务:")
        cli.PrintTask(selectedTask, true)
        fmt.Printf("\n确定要删除这个任务吗？(y/N): ")
        response, _ := reader.ReadString('\n')
        response = strings.TrimSpace(strings.ToLower(response))
        if response != "y" && response != "yes" {
            fmt.Println("已取消")
            return
        }
    }

    // 执行删除
    if err := store.DeleteTask(selectedTask.ID); err != nil {
        cli.PrintError("删除任务失败: %v", err)
        return
    }

    cli.PrintSuccess("任务 %d 已删除", selectedTask.ID)
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "跳过确认")
}
