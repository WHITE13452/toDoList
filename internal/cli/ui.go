package cli

import (
	"fmt"
	"strings"

	"github.com/WHITE13452/toDoList/internal/models"
	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgCyan)
	dimColor     = color.New(color.Faint)
)

// PrintSuccess 打印成功消息
func PrintSuccess(format string, args ...interface{}) {
	successColor.Printf("✓ "+format+"\n", args...)
}

// PrintError 打印错误消息
func PrintError(format string, args ...interface{}) {
	errorColor.Printf("✗ "+format+"\n", args...)
}

// PrintInfo 打印信息
func PrintInfo(format string, args ...interface{}) {
	infoColor.Printf(format+"\n", args...)
}

// PrintTask 打印单个任务
func PrintTask(task *models.Task, detailed bool) {
	statusIcon := "○"
	if task.Status == models.StatusCompleted {
		statusIcon = "✓"
	}

	priorityText := getPriorityText(task.Priority)

	if detailed {
		fmt.Println(strings.Repeat("─", 60))
		fmt.Printf("ID: %d\n", task.ID)
		fmt.Printf("标题: %s\n", task.Title)
		fmt.Printf("状态: %s %s\n", statusIcon, task.Status)
		fmt.Printf("分类: %s\n", task.Category)
		fmt.Printf("优先级: %s\n", priorityText)
		fmt.Printf("创建时间: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("更新时间: %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))
		if task.CompletedAt != nil {
			fmt.Printf("完成时间: %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
		}
		if task.Description != "" {
			fmt.Printf("\n描述:\n%s\n", task.Description)
		}
		fmt.Println(strings.Repeat("─", 60))
	} else {
		if task.Status == models.StatusCompleted {
			successColor.Printf("[%d] %s %s (%s, %s)\n",
				task.ID, statusIcon, task.Title, task.Category, priorityText)
		} else {
			fmt.Printf("[%d] %s %s (%s, %s)\n",
				task.ID, statusIcon, task.Title, task.Category, priorityText)
		}
	}
}

// PrintTaskTable 以表格形式打印任务列表
func PrintTaskTable(tasks []*models.Task) {
	if len(tasks) == 0 {
		dimColor.Println("暂无任务")
		return
	}

	// 打印表头
	fmt.Println(strings.Repeat("═", 80))
	fmt.Printf("%-6s %-6s %-32s %-10s %-8s %-16s\n",
		"ID", "状态", "标题", "分类", "优先级", "创建时间")
	fmt.Println(strings.Repeat("─", 80))

	// 打印任务
	for _, task := range tasks {
		statusIcon := "○"
		if task.Status == models.StatusCompleted {
			statusIcon = "✓"
		}

		priorityStr := strings.Repeat("!", int(task.Priority))

		// 截断长标题
		title := task.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}

		if task.Status == models.StatusCompleted {
			successColor.Printf("%-6d %-6s %-32s %-10s %-8s %-16s\n",
				task.ID, statusIcon, title, task.Category, priorityStr,
				task.CreatedAt.Format("2006-01-02 15:04"))
		} else {
			fmt.Printf("%-6d %-6s %-32s %-10s %-8s %-16s\n",
				task.ID, statusIcon, title, task.Category, priorityStr,
				task.CreatedAt.Format("2006-01-02 15:04"))
		}
	}

	fmt.Println(strings.Repeat("═", 80))
	dimColor.Printf("总计: %d 个任务\n", len(tasks))
}

// PrintStatistics 打印统计信息
func PrintStatistics(stats *models.Statistics) {
	fmt.Println(strings.Repeat("═", 60))
	infoColor.Println("                    📊 统计信息")
	fmt.Println(strings.Repeat("═", 60))

	fmt.Printf("📋 总任务数: %d\n", stats.Total)
	successColor.Printf("✓ 已完成: %d\n", stats.Completed)
	fmt.Printf("○ 待完成: %d\n", stats.Pending)
	fmt.Printf("📈 完成率: %.1f%%\n", stats.CompletionRate)

	if len(stats.ByCategory) > 0 {
		fmt.Println("\n📁 按分类统计:")
		for cat, count := range stats.ByCategory {
			fmt.Printf("  • %s: %d\n", cat, count)
		}
	}

	if len(stats.ByPriority) > 0 {
		fmt.Println("\n⚡ 待办任务优先级分布:")
		priorityNames := map[models.Priority]string{
			models.PriorityLow:    "低",
			models.PriorityMedium: "中",
			models.PriorityHigh:   "高",
			models.PriorityUrgent: "紧急",
		}
		for priority := models.PriorityUrgent; priority >= models.PriorityLow; priority-- {
			if count, ok := stats.ByPriority[priority]; ok {
				fmt.Printf("  • %s: %d\n", priorityNames[priority], count)
			}
		}
	}

	fmt.Println(strings.Repeat("═", 60))
}

// PrintAgentWelcome 打印 Agent 欢迎信息
func PrintAgentWelcome() {
	fmt.Println(strings.Repeat("═", 60))
	infoColor.Println("            🤖 TodoList AI Agent")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println()
	fmt.Println("我是你的智能待办助手，可以帮你管理任务。")
	fmt.Println()
	fmt.Println("你可以问我：")
	fmt.Println("• 'list' 或 '显示所有任务'")
	fmt.Println("• '统计' 或 '总结一下'")
	fmt.Println("• '添加任务：写周报'")
	fmt.Println("• '完成任务 1'")
	fmt.Println("• '搜索包含会议的任务'")
	fmt.Println("• 或者用自然语言描述你想做什么")
	fmt.Println()
	dimColor.Println("输入 'exit' 或 'quit' 退出。")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println()
}

// PrintAgentThinking 打印 Agent 思考中
func PrintAgentThinking(toolName string) {
	dimColor.Printf("🔧 调用工具: %s...\n", toolName)
}

// PrintAgentResponse 打印 Agent 响应
func PrintAgentResponse(response string) {
	fmt.Println(strings.Repeat("─", 60))
	infoColor.Println("Agent:")
	fmt.Println(response)
	fmt.Println(strings.Repeat("─", 60))
}

func getPriorityText(priority models.Priority) string {
	switch priority {
	case models.PriorityLow:
		return "低"
	case models.PriorityMedium:
		return "中"
	case models.PriorityHigh:
		return "高"
	case models.PriorityUrgent:
		return "紧急"
	default:
		return fmt.Sprintf("%d", priority)
	}
}
