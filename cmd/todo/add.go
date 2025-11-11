package main

import (
	"github.com/WHITE13452/toDoList/internal/cli"
	"github.com/WHITE13452/toDoList/internal/models"
	"github.com/spf13/cobra"
)

var (
	taskDescription string
	taskCategory    string
	taskPriority    int
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "添加新任务",
	Long:  "添加一个新的待办事项。标题为必填参数，描述、分类和优先级为可选。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]

		// 验证分类
		category := models.TaskCategory(taskCategory)
		if category != models.CategoryWork && category != models.CategoryStudy &&
			category != models.CategoryLife && category != models.CategoryOther {
			cli.PrintError("无效的分类，必须是 work, study, life 或 other")
			return
		}

		// 验证优先级
		priority := models.Priority(taskPriority)
		if priority < models.PriorityLow || priority > models.PriorityUrgent {
			cli.PrintError("无效的优先级，必须是 1-4")
			return
		}

		// 创建任务
		task := models.NewTask(title, taskDescription, category, priority)

		// 保存任务
		if err := store.AddTask(task); err != nil {
			cli.PrintError("添加任务失败: %v", err)
			return
		}

		cli.PrintSuccess("任务已添加 (ID: %d)", task.ID)
		cli.PrintTask(task, false)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&taskDescription, "description", "d", "", "任务描述")
	addCmd.Flags().StringVarP(&taskCategory, "category", "c", "other", "任务分类 (work/study/life/other)")
	addCmd.Flags().IntVarP(&taskPriority, "priority", "p", 2, "优先级 (1:低 2:中 3:高 4:紧急)")
}
