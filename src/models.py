"""数据模型定义"""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional
from enum import Enum


class TaskStatus(Enum):
    """任务状态枚举"""
    PENDING = "pending"
    COMPLETED = "completed"


class TaskCategory(Enum):
    """任务分类枚举"""
    WORK = "work"
    STUDY = "study"
    LIFE = "life"
    OTHER = "other"


class Priority(Enum):
    """优先级枚举"""
    LOW = 1
    MEDIUM = 2
    HIGH = 3
    URGENT = 4


@dataclass
class Task:
    """待办事项数据模型"""

    title: str
    description: Optional[str] = None
    status: TaskStatus = TaskStatus.PENDING
    category: TaskCategory = TaskCategory.OTHER
    priority: Priority = Priority.MEDIUM
    created_at: datetime = field(default_factory=datetime.now)
    updated_at: datetime = field(default_factory=datetime.now)
    completed_at: Optional[datetime] = None
    id: Optional[int] = None

    def to_dict(self) -> dict:
        """转换为字典"""
        return {
            'id': self.id,
            'title': self.title,
            'description': self.description,
            'status': self.status.value,
            'category': self.category.value,
            'priority': self.priority.value,
            'created_at': self.created_at.isoformat(),
            'updated_at': self.updated_at.isoformat(),
            'completed_at': self.completed_at.isoformat() if self.completed_at else None,
        }

    @classmethod
    def from_dict(cls, data: dict) -> 'Task':
        """从字典创建任务对象"""
        return cls(
            id=data.get('id'),
            title=data['title'],
            description=data.get('description'),
            status=TaskStatus(data.get('status', 'pending')),
            category=TaskCategory(data.get('category', 'other')),
            priority=Priority(data.get('priority', 2)),
            created_at=datetime.fromisoformat(data['created_at']) if isinstance(data.get('created_at'), str) else data.get('created_at', datetime.now()),
            updated_at=datetime.fromisoformat(data['updated_at']) if isinstance(data.get('updated_at'), str) else data.get('updated_at', datetime.now()),
            completed_at=datetime.fromisoformat(data['completed_at']) if data.get('completed_at') else None,
        )

    def mark_completed(self):
        """标记为已完成"""
        self.status = TaskStatus.COMPLETED
        self.completed_at = datetime.now()
        self.updated_at = datetime.now()

    def mark_pending(self):
        """标记为未完成"""
        self.status = TaskStatus.PENDING
        self.completed_at = None
        self.updated_at = datetime.now()

    def __str__(self) -> str:
        """字符串表示"""
        status_icon = "✓" if self.status == TaskStatus.COMPLETED else "○"
        priority_str = "!" * self.priority.value
        return f"[{self.id}] {status_icon} {self.title} ({self.category.value}) {priority_str}"
