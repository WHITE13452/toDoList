package agent

import (
	"context"
	"fmt"

	"github.com/WHITE13452/toDoList/internal/tools"
	"github.com/sashabaranov/go-openai"
)

// Agent AI 智能助手
type Agent struct {
	client   *openai.Client
	tools    *tools.TodoTools
	model    string
	messages []openai.ChatCompletionMessage
}

// Config Agent 配置
type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}

// New 创建新的 Agent 实例
func New(config Config, tools *tools.TodoTools) *Agent {
	clientConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	if config.Model == "" {
		config.Model = "qwen-plus"
	}

	return &Agent{
		client: client,
		tools:  tools,
		model:  config.Model,
		messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `你是一个智能待办事项管理助手。你可以帮助用户管理他们的任务列表。

你的能力包括：
1. 查看和总结待办事项
2. 添加新任务
3. 标记任务完成或未完成
4. 删除任务
5. 搜索特定任务
6. 提供统计信息和分析
7. 批量操作任务

使用技巧：
- 当用户询问任务情况时，先调用 get_all_tasks 或 get_statistics 获取信息
- 对于模糊的任务描述，可以使用 search_tasks 查找
- 批量操作时使用 batch_complete_tasks 或 batch_delete_tasks
- 提供建议时要考虑任务的优先级和分类
- 用清晰、友好的中文与用户交流

重要：
- 在执行删除等重要操作前，最好确认用户的意图
- 提供统计和总结时，用简洁明了的方式呈现
- 如果任务很多，可以先总结再列出重点`,
			},
		},
	}
}

// Chat 与 Agent 对话
func (a *Agent) Chat(ctx context.Context, userMessage string) (string, error) {
	// 添加用户消息
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	// 最多循环 10 次来处理工具调用
	for i := 0; i < 10; i++ {
		// 创建聊天完成请求
		req := openai.ChatCompletionRequest{
			Model:    a.model,
			Messages: a.messages,
			Tools:    a.tools.GetToolDefinitions(),
		}

		resp, err := a.client.CreateChatCompletion(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to create chat completion: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("no response from API")
		}

		choice := resp.Choices[0]
		message := choice.Message

		// 添加助手消息到历史
		a.messages = append(a.messages, message)

		// 如果没有工具调用，返回文本响应
		if len(message.ToolCalls) == 0 {
			return message.Content, nil
		}

		// 处理工具调用
		for _, toolCall := range message.ToolCalls {
			functionName := toolCall.Function.Name
			arguments := toolCall.Function.Arguments

			// 执行工具
			result, err := a.tools.ExecuteTool(functionName, arguments)
			if err != nil {
				result = fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error())
			}

			// 添加工具结果到消息历史
			a.messages = append(a.messages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    result,
				ToolCallID: toolCall.ID,
			})
		}

		// 继续循环以获取最终响应
	}

	return "", fmt.Errorf("too many tool calls, conversation limit reached")
}

// ClearHistory 清空对话历史
func (a *Agent) ClearHistory() {
	// 保留系统提示，清除其他消息
	systemMessage := a.messages[0]
	a.messages = []openai.ChatCompletionMessage{systemMessage}
}

// GetHistory 获取对话历史
func (a *Agent) GetHistory() []openai.ChatCompletionMessage {
	return a.messages
}
