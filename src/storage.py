"""数据持久化层 - 使用 SQLite"""

import sqlite3
from pathlib import Path
from typing import List, Optional
from datetime import datetime
import os

from .models import Task, TaskStatus, TaskCategory, Priority


class TaskStorage:
    """任务存储管理器"""

    def __init__(self, db_path: Optional[str] = None):
        """初始化存储

        Args:
            db_path: 数据库文件路径，默认为用户家目录下的 .todolist.db
        """
        if db_path is None:
            home = Path.home()
            db_path = str(home / ".todolist.db")

        self.db_path = db_path
        self._init_database()

    def _init_database(self):
        """初始化数据库表结构"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute("""
            CREATE TABLE IF NOT EXISTS tasks (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                title TEXT NOT NULL,
                description TEXT,
                status TEXT NOT NULL DEFAULT 'pending',
                category TEXT NOT NULL DEFAULT 'other',
                priority INTEGER NOT NULL DEFAULT 2,
                created_at TEXT NOT NULL,
                updated_at TEXT NOT NULL,
                completed_at TEXT
            )
        """)

        conn.commit()
        conn.close()

    def add_task(self, task: Task) -> Task:
        """添加任务

        Args:
            task: 任务对象

        Returns:
            添加后的任务对象（包含 ID）
        """
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute("""
            INSERT INTO tasks (title, description, status, category, priority,
                             created_at, updated_at, completed_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            task.title,
            task.description,
            task.status.value,
            task.category.value,
            task.priority.value,
            task.created_at.isoformat(),
            task.updated_at.isoformat(),
            task.completed_at.isoformat() if task.completed_at else None
        ))

        task.id = cursor.lastrowid
        conn.commit()
        conn.close()

        return task

    def get_task(self, task_id: int) -> Optional[Task]:
        """获取单个任务

        Args:
            task_id: 任务 ID

        Returns:
            任务对象或 None
        """
        conn = sqlite3.connect(self.db_path)
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()

        cursor.execute("SELECT * FROM tasks WHERE id = ?", (task_id,))
        row = cursor.fetchone()
        conn.close()

        if row:
            return self._row_to_task(row)
        return None

    def get_all_tasks(self, status: Optional[TaskStatus] = None,
                     category: Optional[TaskCategory] = None) -> List[Task]:
        """获取所有任务

        Args:
            status: 可选的状态过滤
            category: 可选的分类过滤

        Returns:
            任务列表
        """
        conn = sqlite3.connect(self.db_path)
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()

        query = "SELECT * FROM tasks WHERE 1=1"
        params = []

        if status:
            query += " AND status = ?"
            params.append(status.value)

        if category:
            query += " AND category = ?"
            params.append(category.value)

        query += " ORDER BY priority DESC, created_at DESC"

        cursor.execute(query, params)
        rows = cursor.fetchall()
        conn.close()

        return [self._row_to_task(row) for row in rows]

    def update_task(self, task: Task) -> bool:
        """更新任务

        Args:
            task: 任务对象

        Returns:
            是否更新成功
        """
        if not task.id:
            return False

        task.updated_at = datetime.now()

        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute("""
            UPDATE tasks
            SET title = ?, description = ?, status = ?, category = ?,
                priority = ?, updated_at = ?, completed_at = ?
            WHERE id = ?
        """, (
            task.title,
            task.description,
            task.status.value,
            task.category.value,
            task.priority.value,
            task.updated_at.isoformat(),
            task.completed_at.isoformat() if task.completed_at else None,
            task.id
        ))

        success = cursor.rowcount > 0
        conn.commit()
        conn.close()

        return success

    def delete_task(self, task_id: int) -> bool:
        """删除任务

        Args:
            task_id: 任务 ID

        Returns:
            是否删除成功
        """
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute("DELETE FROM tasks WHERE id = ?", (task_id,))
        success = cursor.rowcount > 0

        conn.commit()
        conn.close()

        return success

    def search_tasks(self, keyword: str) -> List[Task]:
        """搜索任务

        Args:
            keyword: 搜索关键词

        Returns:
            匹配的任务列表
        """
        conn = sqlite3.connect(self.db_path)
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()

        cursor.execute("""
            SELECT * FROM tasks
            WHERE title LIKE ? OR description LIKE ?
            ORDER BY priority DESC, created_at DESC
        """, (f"%{keyword}%", f"%{keyword}%"))

        rows = cursor.fetchall()
        conn.close()

        return [self._row_to_task(row) for row in rows]

    def get_statistics(self) -> dict:
        """获取统计信息

        Returns:
            包含统计信息的字典
        """
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        # 总任务数
        cursor.execute("SELECT COUNT(*) FROM tasks")
        total = cursor.fetchone()[0]

        # 已完成任务数
        cursor.execute("SELECT COUNT(*) FROM tasks WHERE status = 'completed'")
        completed = cursor.fetchone()[0]

        # 待办任务数
        pending = total - completed

        # 按分类统计
        cursor.execute("""
            SELECT category, COUNT(*) as count
            FROM tasks
            GROUP BY category
        """)
        by_category = {row[0]: row[1] for row in cursor.fetchall()}

        # 按优先级统计（仅待办）
        cursor.execute("""
            SELECT priority, COUNT(*) as count
            FROM tasks
            WHERE status = 'pending'
            GROUP BY priority
        """)
        by_priority = {row[0]: row[1] for row in cursor.fetchall()}

        conn.close()

        return {
            'total': total,
            'completed': completed,
            'pending': pending,
            'completion_rate': (completed / total * 100) if total > 0 else 0,
            'by_category': by_category,
            'by_priority': by_priority
        }

    def _row_to_task(self, row: sqlite3.Row) -> Task:
        """将数据库行转换为任务对象"""
        return Task(
            id=row['id'],
            title=row['title'],
            description=row['description'],
            status=TaskStatus(row['status']),
            category=TaskCategory(row['category']),
            priority=Priority(row['priority']),
            created_at=datetime.fromisoformat(row['created_at']),
            updated_at=datetime.fromisoformat(row['updated_at']),
            completed_at=datetime.fromisoformat(row['completed_at']) if row['completed_at'] else None
        )
