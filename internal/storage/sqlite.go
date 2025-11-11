package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/WHITE13452/toDoList/internal/models"
)

// Storage SQLite 存储实现
type Storage struct {
	db *sql.DB
}

// New 创建新的存储实例
func New(dbPath string) (*Storage, error) {
	// 如果没有指定路径，使用默认路径
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dbPath = filepath.Join(home, ".todolist.db")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &Storage{db: db}
	if err := storage.initDatabase(); err != nil {
		return nil, err
	}

	return storage, nil
}

// initDatabase 初始化数据库表
func (s *Storage) initDatabase() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		category TEXT NOT NULL DEFAULT 'other',
		priority INTEGER NOT NULL DEFAULT 2,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		completed_at DATETIME
	);
	CREATE INDEX IF NOT EXISTS idx_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_category ON tasks(category);
	CREATE INDEX IF NOT EXISTS idx_priority ON tasks(priority);
	`

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	return nil
}

// AddTask 添加任务
func (s *Storage) AddTask(task *models.Task) error {
	query := `
	INSERT INTO tasks (title, description, status, category, priority, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(query,
		task.Title, task.Description, task.Status, task.Category,
		task.Priority, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.ID = id
	return nil
}

// GetTask 获取单个任务
func (s *Storage) GetTask(id int64) (*models.Task, error) {
	query := `
	SELECT id, title, description, status, category, priority,
	       created_at, updated_at, completed_at
	FROM tasks WHERE id = ?
	`

	var task models.Task
	var completedAt sql.NullTime

	err := s.db.QueryRow(query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status,
		&task.Category, &task.Priority, &task.CreatedAt,
		&task.UpdatedAt, &completedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, nil
}

// GetAllTasks 获取所有任务
func (s *Storage) GetAllTasks(status models.TaskStatus, category models.TaskCategory) ([]*models.Task, error) {
	query := `
	SELECT id, title, description, status, category, priority,
	       created_at, updated_at, completed_at
	FROM tasks WHERE 1=1
	`
	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY priority DESC, created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		var completedAt sql.NullTime

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status,
			&task.Category, &task.Priority, &task.CreatedAt,
			&task.UpdatedAt, &completedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// UpdateTask 更新任务
func (s *Storage) UpdateTask(task *models.Task) error {
	query := `
	UPDATE tasks
	SET title = ?, description = ?, status = ?, category = ?,
	    priority = ?, updated_at = ?, completed_at = ?
	WHERE id = ?
	`

	task.UpdatedAt = time.Now()

	var completedAt interface{}
	if task.CompletedAt != nil {
		completedAt = task.CompletedAt
	}

	result, err := s.db.Exec(query,
		task.Title, task.Description, task.Status, task.Category,
		task.Priority, task.UpdatedAt, completedAt, task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// DeleteTask 删除任务
func (s *Storage) DeleteTask(id int64) error {
	query := "DELETE FROM tasks WHERE id = ?"

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// SearchTasks 搜索任务
func (s *Storage) SearchTasks(keyword string) ([]*models.Task, error) {
	query := `
	SELECT id, title, description, status, category, priority,
	       created_at, updated_at, completed_at
	FROM tasks
	WHERE title LIKE ? OR description LIKE ?
	ORDER BY priority DESC, created_at DESC
	`

	pattern := "%" + keyword + "%"
	rows, err := s.db.Query(query, pattern, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		var completedAt sql.NullTime

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status,
			&task.Category, &task.Priority, &task.CreatedAt,
			&task.UpdatedAt, &completedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// GetStatistics 获取统计信息
func (s *Storage) GetStatistics() (*models.Statistics, error) {
	stats := &models.Statistics{
		ByCategory: make(map[models.TaskCategory]int),
		ByPriority: make(map[models.Priority]int),
	}

	// 总数和完成数
	err := s.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&stats.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	err = s.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'completed'").Scan(&stats.Completed)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed count: %w", err)
	}

	stats.Pending = stats.Total - stats.Completed

	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Completed) / float64(stats.Total) * 100
	}

	// 按分类统计
	rows, err := s.db.Query("SELECT category, COUNT(*) FROM tasks GROUP BY category")
	if err != nil {
		return nil, fmt.Errorf("failed to get category stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category models.TaskCategory
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, err
		}
		stats.ByCategory[category] = count
	}

	// 按优先级统计（仅待办）
	rows, err = s.db.Query("SELECT priority, COUNT(*) FROM tasks WHERE status = 'pending' GROUP BY priority")
	if err != nil {
		return nil, fmt.Errorf("failed to get priority stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var priority models.Priority
		var count int
		if err := rows.Scan(&priority, &count); err != nil {
			return nil, err
		}
		stats.ByPriority[priority] = count
	}

	return stats, nil
}

// Close 关闭数据库连接
func (s *Storage) Close() error {
	return s.db.Close()
}
