package tools

import (
	"encoding/json"
	"fmt"

	"github.com/WHITE13452/toDoList/internal/models"
	"github.com/WHITE13452/toDoList/internal/storage"
	"github.com/sashabaranov/go-openai"
)

// TodoTools AI Agent 可调用的工具集
type TodoTools struct {
	storage *storage.Storage
}

// New 创建工具实例
func New(storage *storage.Storage) *TodoTools {
	return &TodoTools{storage: storage}
}

// GetToolDefinitions 获取工具定义（OpenAI Function Calling 格式）
func (t *TodoTools) GetToolDefinitions() []openai.Tool {
	return []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_all_tasks",
				Description: "获取所有待办事项列表。可以根据状态（pending/completed）或分类（work/study/life/other）进行过滤。",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"status": {
							"type": "string",
							"enum": ["pending", "completed"],
							"description": "任务状态过滤：pending(待办) 或 completed(已完成)"
						},
						"category": {
							"type": "string",
							"enum": ["work", "study", "life", "other"],
							"description": "任务分类过滤：work(工作)、study(学习)、life(生活)、other(其他)"
						}
					}
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "add_task",
				Description: "添加一个新的待办事项。需要提供标题，描述、分类和优先级为可选参数。",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"title": {
							"type": "string",
							"description": "任务标题（必填）"
						},
						"description": {
							"type": "string",
							"description": "任务描述（可选）"
						},
						"category": {
							"type": "string",
							"enum": ["work", "study", "life", "other"],
							"description": "任务分类，默认为 other"
						},
						"priority": {
							"type": "integer",
							"enum": [1, 2, 3, 4],
							"description": "优先级：1(低)、2(中)、3(高)、4(紧急)，默认为 2"
						}
					},
					"required": ["title"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "update_task_status",
				Description: "更新任务的完成状态。可以标记任务为已完成或未完成。",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"task_id": {
							"type": "integer",
							"description": "要更新的任务 ID"
						},
						"status": {
							"type": "string",
							"enum": ["pending", "completed"],
							"description": "新的任务状态：pending(未完成) 或 completed(已完成)"
						}
					},
					"required": ["task_id", "status"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "delete_task",
				Description: "删除指定的待办事项。",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"task_id": {
							"type": "integer",
							"description": "要删除的任务 ID"
						}
					},
					"required": ["task_id"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "search_tasks",
				Description: "在待办事项中搜索包含指定关键词的任务（标题或描述）。",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"keyword": {
							"type": "string",
							"description": "搜索关键词"
						}
					},
					"required": ["keyword"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_statistics",
				Description: "获取待办事项的统计信息，包括总数、完成数、待办数、完成率、分类统计和优先级分布。",
				Parameters:  json.RawMessage(`{"type": "object", "properties": {}}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_task_detail",
				Description: "获取指定任务的详细信息。",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"task_id": {
							"type": "integer",
							"description": "任务 ID"
						}
					},
					"required": ["task_id"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "batch_complete_tasks",
				Description: "批量标记多个任务为已完成。",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"task_ids": {
							"type": "array",
							"items": {"type": "integer"},
							"description": "要标记为完成的任务 ID 列表"
						}
					},
					"required": ["task_ids"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "batch_delete_tasks",
				Description: "批量删除多个任务。",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"task_ids": {
							"type": "array",
							"items": {"type": "integer"},
							"description": "要删除的任务 ID 列表"
						}
					},
					"required": ["task_ids"]
				}`),
			},
		},
	}
}

// ExecuteTool 执行工具调用
func (t *TodoTools) ExecuteTool(name, arguments string) (string, error) {
	switch name {
	case "get_all_tasks":
		return t.getAllTasks(arguments)
	case "add_task":
		return t.addTask(arguments)
	case "update_task_status":
		return t.updateTaskStatus(arguments)
	case "delete_task":
		return t.deleteTask(arguments)
	case "search_tasks":
		return t.searchTasks(arguments)
	case "get_statistics":
		return t.getStatistics()
	case "get_task_detail":
		return t.getTaskDetail(arguments)
	case "batch_complete_tasks":
		return t.batchCompleteTasks(arguments)
	case "batch_delete_tasks":
		return t.batchDeleteTasks(arguments)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func (t *TodoTools) getAllTasks(arguments string) (string, error) {
	var args struct {
		Status   models.TaskStatus   `json:"status"`
		Category models.TaskCategory `json:"category"`
	}

	if arguments != "" && arguments != "{}" {
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return "", fmt.Errorf("failed to parse arguments: %w", err)
		}
	}

	tasks, err := t.storage.GetAllTasks(args.Status, args.Category)
	if err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"success": true,
		"count":   len(tasks),
		"tasks":   tasks,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t *TodoTools) addTask(arguments string) (string, error) {
	var args struct {
		Title       string               `json:"title"`
		Description string               `json:"description"`
		Category    models.TaskCategory  `json:"category"`
		Priority    models.Priority      `json:"priority"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	// 设置默认值
	if args.Category == "" {
		args.Category = models.CategoryOther
	}
	if args.Priority == 0 {
		args.Priority = models.PriorityMedium
	}

	task := models.NewTask(args.Title, args.Description, args.Category, args.Priority)
	if err := t.storage.AddTask(task); err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("任务已添加，ID: %d", task.ID),
		"task":    task,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t *TodoTools) updateTaskStatus(arguments string) (string, error) {
	var args struct {
		TaskID int64             `json:"task_id"`
		Status models.TaskStatus `json:"status"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	task, err := t.storage.GetTask(args.TaskID)
	if err != nil {
		return "", err
	}
	if task == nil {
		result := map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("任务 %d 不存在", args.TaskID),
		}
		data, _ := json.Marshal(result)
		return string(data), nil
	}

	if args.Status == models.StatusCompleted {
		task.MarkCompleted()
	} else {
		task.MarkPending()
	}

	if err := t.storage.UpdateTask(task); err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("任务 %d 已标记为 %s", args.TaskID, args.Status),
		"task":    task,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t *TodoTools) deleteTask(arguments string) (string, error) {
	var args struct {
		TaskID int64 `json:"task_id"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	if err := t.storage.DeleteTask(args.TaskID); err != nil {
		result := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		data, _ := json.Marshal(result)
		return string(data), nil
	}

	result := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("任务 %d 已删除", args.TaskID),
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t *TodoTools) searchTasks(arguments string) (string, error) {
	var args struct {
		Keyword string `json:"keyword"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	tasks, err := t.storage.SearchTasks(args.Keyword)
	if err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"success": true,
		"count":   len(tasks),
		"keyword": args.Keyword,
		"tasks":   tasks,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t *TodoTools) getStatistics() (string, error) {
	stats, err := t.storage.GetStatistics()
	if err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"success":    true,
		"statistics": stats,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t *TodoTools) getTaskDetail(arguments string) (string, error) {
	var args struct {
		TaskID int64 `json:"task_id"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	task, err := t.storage.GetTask(args.TaskID)
	if err != nil {
		return "", err
	}

	if task == nil {
		result := map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("任务 %d 不存在", args.TaskID),
		}
		data, _ := json.Marshal(result)
		return string(data), nil
	}

	result := map[string]interface{}{
		"success": true,
		"task":    task,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t *TodoTools) batchCompleteTasks(arguments string) (string, error) {
	var args struct {
		TaskIDs []int64 `json:"task_ids"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	successCount := 0
	failedIDs := []int64{}

	for _, id := range args.TaskIDs {
		task, err := t.storage.GetTask(id)
		if err != nil || task == nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		task.MarkCompleted()
		if err := t.storage.UpdateTask(task); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successCount++
	}

	result := map[string]interface{}{
		"success":       true,
		"message":       fmt.Sprintf("成功标记 %d 个任务为已完成", successCount),
		"success_count": successCount,
		"failed_ids":    failedIDs,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t *TodoTools) batchDeleteTasks(arguments string) (string, error) {
	var args struct {
		TaskIDs []int64 `json:"task_ids"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	successCount := 0
	failedIDs := []int64{}

	for _, id := range args.TaskIDs {
		if err := t.storage.DeleteTask(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		successCount++
	}

	result := map[string]interface{}{
		"success":       true,
		"message":       fmt.Sprintf("成功删除 %d 个任务", successCount),
		"success_count": successCount,
		"failed_ids":    failedIDs,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
