"""AI Agent 可调用的工具定义"""

from typing import List, Optional, Dict, Any
from .storage import TaskStorage
from .models import Task, TaskStatus, TaskCategory, Priority


class TodoTools:
    """TodoList 工具集合，供 AI Agent 调用"""

    def __init__(self, storage: TaskStorage):
        self.storage = storage

    def get_tools_definition(self) -> List[Dict[str, Any]]:
        """获取工具定义，用于 Claude API

        Returns:
            工具定义列表
        """
        return [
            {
                "name": "get_all_tasks",
                "description": "获取所有待办事项列表。可以根据状态（pending/completed）或分类（work/study/life/other）进行过滤。",
                "input_schema": {
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
                }
            },
            {
                "name": "add_task",
                "description": "添加一个新的待办事项。需要提供标题，描述、分类和优先级为可选参数。",
                "input_schema": {
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
                }
            },
            {
                "name": "update_task_status",
                "description": "更新任务的完成状态。可以标记任务为已完成或未完成。",
                "input_schema": {
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
                }
            },
            {
                "name": "delete_task",
                "description": "删除指定的待办事项。",
                "input_schema": {
                    "type": "object",
                    "properties": {
                        "task_id": {
                            "type": "integer",
                            "description": "要删除的任务 ID"
                        }
                    },
                    "required": ["task_id"]
                }
            },
            {
                "name": "search_tasks",
                "description": "在待办事项中搜索包含指定关键词的任务（标题或描述）。",
                "input_schema": {
                    "type": "object",
                    "properties": {
                        "keyword": {
                            "type": "string",
                            "description": "搜索关键词"
                        }
                    },
                    "required": ["keyword"]
                }
            },
            {
                "name": "get_statistics",
                "description": "获取待办事项的统计信息，包括总数、完成数、待办数、完成率、分类统计和优先级分布。",
                "input_schema": {
                    "type": "object",
                    "properties": {}
                }
            },
            {
                "name": "get_task_detail",
                "description": "获取指定任务的详细信息。",
                "input_schema": {
                    "type": "object",
                    "properties": {
                        "task_id": {
                            "type": "integer",
                            "description": "任务 ID"
                        }
                    },
                    "required": ["task_id"]
                }
            },
            {
                "name": "batch_complete_tasks",
                "description": "批量标记多个任务为已完成。",
                "input_schema": {
                    "type": "object",
                    "properties": {
                        "task_ids": {
                            "type": "array",
                            "items": {"type": "integer"},
                            "description": "要标记为完成的任务 ID 列表"
                        }
                    },
                    "required": ["task_ids"]
                }
            },
            {
                "name": "batch_delete_tasks",
                "description": "批量删除多个任务。",
                "input_schema": {
                    "type": "object",
                    "properties": {
                        "task_ids": {
                            "type": "array",
                            "items": {"type": "integer"},
                            "description": "要删除的任务 ID 列表"
                        }
                    },
                    "required": ["task_ids"]
                }
            }
        ]

    def execute_tool(self, tool_name: str, tool_input: Dict[str, Any]) -> Dict[str, Any]:
        """执行工具调用

        Args:
            tool_name: 工具名称
            tool_input: 工具输入参数

        Returns:
            执行结果
        """
        try:
            if tool_name == "get_all_tasks":
                return self._get_all_tasks(
                    status=tool_input.get("status"),
                    category=tool_input.get("category")
                )

            elif tool_name == "add_task":
                return self._add_task(
                    title=tool_input["title"],
                    description=tool_input.get("description"),
                    category=tool_input.get("category", "other"),
                    priority=tool_input.get("priority", 2)
                )

            elif tool_name == "update_task_status":
                return self._update_task_status(
                    task_id=tool_input["task_id"],
                    status=tool_input["status"]
                )

            elif tool_name == "delete_task":
                return self._delete_task(task_id=tool_input["task_id"])

            elif tool_name == "search_tasks":
                return self._search_tasks(keyword=tool_input["keyword"])

            elif tool_name == "get_statistics":
                return self._get_statistics()

            elif tool_name == "get_task_detail":
                return self._get_task_detail(task_id=tool_input["task_id"])

            elif tool_name == "batch_complete_tasks":
                return self._batch_complete_tasks(task_ids=tool_input["task_ids"])

            elif tool_name == "batch_delete_tasks":
                return self._batch_delete_tasks(task_ids=tool_input["task_ids"])

            else:
                return {"success": False, "error": f"未知工具: {tool_name}"}

        except Exception as e:
            return {"success": False, "error": str(e)}

    def _get_all_tasks(self, status: Optional[str] = None,
                      category: Optional[str] = None) -> Dict[str, Any]:
        """获取所有任务"""
        status_filter = TaskStatus(status) if status else None
        category_filter = TaskCategory(category) if category else None

        tasks = self.storage.get_all_tasks(
            status=status_filter,
            category=category_filter
        )

        return {
            "success": True,
            "count": len(tasks),
            "tasks": [task.to_dict() for task in tasks]
        }

    def _add_task(self, title: str, description: Optional[str] = None,
                 category: str = "other", priority: int = 2) -> Dict[str, Any]:
        """添加任务"""
        task = Task(
            title=title,
            description=description,
            category=TaskCategory(category),
            priority=Priority(priority)
        )
        task = self.storage.add_task(task)

        return {
            "success": True,
            "message": f"任务已添加，ID: {task.id}",
            "task": task.to_dict()
        }

    def _update_task_status(self, task_id: int, status: str) -> Dict[str, Any]:
        """更新任务状态"""
        task = self.storage.get_task(task_id)
        if not task:
            return {"success": False, "error": f"任务 {task_id} 不存在"}

        if status == "completed":
            task.mark_completed()
        else:
            task.mark_pending()

        self.storage.update_task(task)

        return {
            "success": True,
            "message": f"任务 {task_id} 已标记为{status}",
            "task": task.to_dict()
        }

    def _delete_task(self, task_id: int) -> Dict[str, Any]:
        """删除任务"""
        task = self.storage.get_task(task_id)
        if not task:
            return {"success": False, "error": f"任务 {task_id} 不存在"}

        if self.storage.delete_task(task_id):
            return {
                "success": True,
                "message": f"任务 {task_id} 已删除"
            }
        else:
            return {"success": False, "error": "删除失败"}

    def _search_tasks(self, keyword: str) -> Dict[str, Any]:
        """搜索任务"""
        tasks = self.storage.search_tasks(keyword)

        return {
            "success": True,
            "count": len(tasks),
            "keyword": keyword,
            "tasks": [task.to_dict() for task in tasks]
        }

    def _get_statistics(self) -> Dict[str, Any]:
        """获取统计信息"""
        stats = self.storage.get_statistics()

        return {
            "success": True,
            "statistics": stats
        }

    def _get_task_detail(self, task_id: int) -> Dict[str, Any]:
        """获取任务详情"""
        task = self.storage.get_task(task_id)
        if not task:
            return {"success": False, "error": f"任务 {task_id} 不存在"}

        return {
            "success": True,
            "task": task.to_dict()
        }

    def _batch_complete_tasks(self, task_ids: List[int]) -> Dict[str, Any]:
        """批量完成任务"""
        success_count = 0
        failed_ids = []

        for task_id in task_ids:
            task = self.storage.get_task(task_id)
            if task:
                task.mark_completed()
                if self.storage.update_task(task):
                    success_count += 1
                else:
                    failed_ids.append(task_id)
            else:
                failed_ids.append(task_id)

        return {
            "success": True,
            "message": f"成功标记 {success_count} 个任务为已完成",
            "success_count": success_count,
            "failed_ids": failed_ids
        }

    def _batch_delete_tasks(self, task_ids: List[int]) -> Dict[str, Any]:
        """批量删除任务"""
        success_count = 0
        failed_ids = []

        for task_id in task_ids:
            if self.storage.delete_task(task_id):
                success_count += 1
            else:
                failed_ids.append(task_id)

        return {
            "success": True,
            "message": f"成功删除 {success_count} 个任务",
            "success_count": success_count,
            "failed_ids": failed_ids
        }
