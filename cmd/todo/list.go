package main

import (
    "github.com/WHITE13452/toDoList/internal/cli"
    "github.com/WHITE13452/toDoList/internal/models"
    "github.com/spf13/cobra"
)

var (
    filterStatus   string
    filterCategory string
    sortBy         string  // 新增:排序字段
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "列出任务",
    Long:  "列出所有待办事项。可以使用 -s 和 -c 参数进行过滤,-o 参数进行排序。",
    Run: func(cmd *cobra.Command, args []string) {
        var status models.TaskStatus
        var category models.TaskCategory

        if filterStatus != "" {
            status = models.TaskStatus(filterStatus)
            if status != models.StatusPending && status != models.StatusCompleted {
                cli.PrintError("无效的状态,必须是 pending 或 completed")
                return
            }
        }

        if filterCategory != "" {
            category = models.TaskCategory(filterCategory)
            if category != models.CategoryWork && category != models.CategoryStudy &&
                category != models.CategoryLife && category != models.CategoryOther {
                cli.PrintError("无效的分类,必须是 work, study, life 或 other")
                return
            }
        }

        // 验证排序参数
        if sortBy != "" && sortBy != "priority" && sortBy != "created_at" && sortBy != "updated_at" {
            cli.PrintError("无效的排序字段,必须是 priority, created_at 或 updated_at")
            return
        }

        tasks, err := store.GetAllTasks(status, category, sortBy)
        if err != nil {
            cli.PrintError("获取任务列表失败: %v", err)
            return
        }

        cli.PrintTaskTable(tasks)
    },
}

func init() {
    rootCmd.AddCommand(listCmd)

    listCmd.Flags().StringVarP(&filterStatus, "status", "s", "", "按状态过滤 (pending/completed)")
    listCmd.Flags().StringVarP(&filterCategory, "category", "c", "", "按分类过滤 (work/study/life/other)")
    listCmd.Flags().StringVarP(&sortBy, "sort", "o", "", "排序方式 (priority/created_at/updated_at)")
}