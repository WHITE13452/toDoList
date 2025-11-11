package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/WHITE13452/toDoList/internal/agent"
	"github.com/WHITE13452/toDoList/internal/cli"
	"github.com/WHITE13452/toDoList/internal/tools"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "启动 AI Agent 交互模式",
	Long: `启动 AI Agent 交互模式。

在这个模式下，你可以用自然语言与 AI 助手对话来管理任务。

示例：
• "显示所有未完成的任务"
• "帮我添加一个任务：准备项目演示"
• "完成任务 3"
• "有哪些工作相关的任务？"
• "给我一个总结"

环境变量：
• QWEN_API_KEY - Qwen API Key（必需）
• QWEN_API_BASE - API Base URL（可选，默认: https://dashscope.aliyuncs.com/compatible-mode/v1）
• QWEN_MODEL - 模型名称（可选，默认: qwen-plus）`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取配置
		apiKey := os.Getenv("QWEN_API_KEY")
		if apiKey == "" {
			cli.PrintError("未找到 QWEN_API_KEY 环境变量")
			fmt.Println("\n请确保设置了 QWEN_API_KEY 环境变量。")
			fmt.Println("你可以创建一个 .env 文件并添加：")
			fmt.Println("QWEN_API_KEY=your_api_key_here")
			return
		}

		baseURL := os.Getenv("QWEN_API_BASE")
		if baseURL == "" {
			baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
		}

		model := os.Getenv("QWEN_MODEL")
		if model == "" {
			model = "qwen-plus"
		}

		// 创建 Agent
		todoTools := tools.New(store)
		agentInstance := agent.New(agent.Config{
			APIKey:  apiKey,
			BaseURL: baseURL,
			Model:   model,
		}, todoTools)

		// 打印欢迎信息
		cli.PrintAgentWelcome()

		// 交互循环
		reader := bufio.NewReader(os.Stdin)
		ctx := context.Background()

		for {
			// 获取用户输入
			fmt.Print("\n你: ")
			userInput, err := reader.ReadString('\n')
			if err != nil {
				cli.PrintError("读取输入失败: %v", err)
				break
			}

			userInput = strings.TrimSpace(userInput)
			if userInput == "" {
				continue
			}

			// 检查退出命令
			if userInput == "exit" || userInput == "quit" || userInput == "退出" || userInput == "q" {
				cli.PrintInfo("再见！")
				break
			}

			// 处理快捷命令
			switch userInput {
			case "list", "ls", "列表", "显示":
				userInput = "显示所有待办任务"
			case "stats", "statistics", "统计":
				userInput = "显示统计信息和总结"
			case "help", "h", "帮助":
				fmt.Println("\n可用命令：")
				fmt.Println("• list/ls - 显示所有任务")
				fmt.Println("• stats - 显示统计信息")
				fmt.Println("• help - 显示此帮助")
				fmt.Println("• exit - 退出")
				fmt.Println("\n或者直接用自然语言描述你想做什么，例如：")
				fmt.Println("• '帮我添加一个任务：准备项目演示'")
				fmt.Println("• '完成任务 3'")
				fmt.Println("• '有哪些工作相关的未完成任务？'")
				continue
			case "clear", "cls", "清屏":
				agentInstance.ClearHistory()
				cli.PrintInfo("对话历史已清空")
				continue
			}

			// 与 Agent 对话
			fmt.Println()
			response, err := agentInstance.Chat(ctx, userInput)
			if err != nil {
				cli.PrintError("Agent 错误: %v", err)
				continue
			}

			// 显示响应
			cli.PrintAgentResponse(response)
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
