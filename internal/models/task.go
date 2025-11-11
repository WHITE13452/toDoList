package models

import (
	"time"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
)

// TaskCategory 任务分类
type TaskCategory string

const (
	CategoryWork  TaskCategory = "work"
	CategoryStudy TaskCategory = "study"
	CategoryLife  TaskCategory = "life"
	CategoryOther TaskCategory = "other"
)

// Priority 优先级
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
	PriorityUrgent Priority = 4
)

// Task 待办事项
type Task struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Status      TaskStatus   `json:"status"`
	Category    TaskCategory `json:"category"`
	Priority    Priority     `json:"priority"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
}

// MarkCompleted 标记为已完成
func (t *Task) MarkCompleted() {
	t.Status = StatusCompleted
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// MarkPending 标记为未完成
func (t *Task) MarkPending() {
	t.Status = StatusPending
	t.CompletedAt = nil
	t.UpdatedAt = time.Now()
}

// NewTask 创建新任务
func NewTask(title, description string, category TaskCategory, priority Priority) *Task {
	now := time.Now()
	return &Task{
		Title:       title,
		Description: description,
		Status:      StatusPending,
		Category:    category,
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Statistics 统计信息
type Statistics struct {
	Total          int                    `json:"total"`
	Completed      int                    `json:"completed"`
	Pending        int                    `json:"pending"`
	CompletionRate float64                `json:"completion_rate"`
	ByCategory     map[TaskCategory]int   `json:"by_category"`
	ByPriority     map[Priority]int       `json:"by_priority"`
}
