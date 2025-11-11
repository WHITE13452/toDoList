"""CLI 命令实现"""

import click
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich.text import Text
from typing import Optional

from .models import Task, TaskStatus, TaskCategory, Priority
from .storage import TaskStorage

console = Console()


class CLI:
    """CLI 命令处理器"""

    def __init__(self, storage: TaskStorage):
        self.storage = storage

    def add_task(self, title: str, description: Optional[str] = None,
                 category: str = "other", priority: int = 2):
        """添加任务"""
        try:
            task = Task(
                title=title,
                description=description,
                category=TaskCategory(category),
                priority=Priority(priority)
            )
            task = self.storage.add_task(task)
            console.print(f"✓ 任务已添加 (ID: {task.id})", style="bold green")
            self._display_task(task)
        except Exception as e:
            console.print(f"✗ 添加失败: {str(e)}", style="bold red")

    def list_tasks(self, status: Optional[str] = None,
                  category: Optional[str] = None,
                  show_all: bool = False):
        """列出任务"""
        try:
            status_filter = TaskStatus(status) if status else None
            category_filter = TaskCategory(category) if category else None

            tasks = self.storage.get_all_tasks(
                status=status_filter,
                category=category_filter
            )

            if not tasks:
                console.print("暂无任务", style="yellow")
                return

            # 创建表格
            table = Table(title="📋 任务列表", show_header=True, header_style="bold magenta")
            table.add_column("ID", style="cyan", width=6)
            table.add_column("状态", width=6)
            table.add_column("标题", style="white", min_width=20)
            table.add_column("分类", width=8)
            table.add_column("优先级", width=8)
            table.add_column("创建时间", width=16)

            for task in tasks:
                status_icon = "✓" if task.status == TaskStatus.COMPLETED else "○"
                status_style = "green" if task.status == TaskStatus.COMPLETED else "yellow"
                priority_str = "!" * task.priority.value

                # 截断长标题
                title = task.title if len(task.title) <= 30 else task.title[:27] + "..."

                table.add_row(
                    str(task.id),
                    Text(status_icon, style=status_style),
                    title,
                    task.category.value,
                    priority_str,
                    task.created_at.strftime("%Y-%m-%d %H:%M")
                )

            console.print(table)
            console.print(f"\n总计: {len(tasks)} 个任务", style="dim")

        except Exception as e:
            console.print(f"✗ 列出任务失败: {str(e)}", style="bold red")

    def complete_task(self, task_id: int, uncomplete: bool = False):
        """标记任务完成/未完成"""
        try:
            task = self.storage.get_task(task_id)
            if not task:
                console.print(f"✗ 任务 {task_id} 不存在", style="bold red")
                return

            if uncomplete:
                task.mark_pending()
                self.storage.update_task(task)
                console.print(f"✓ 任务 {task_id} 已标记为未完成", style="bold green")
            else:
                task.mark_completed()
                self.storage.update_task(task)
                console.print(f"✓ 任务 {task_id} 已完成", style="bold green")

            self._display_task(task)

        except Exception as e:
            console.print(f"✗ 操作失败: {str(e)}", style="bold red")

    def delete_task(self, task_id: int):
        """删除任务"""
        try:
            task = self.storage.get_task(task_id)
            if not task:
                console.print(f"✗ 任务 {task_id} 不存在", style="bold red")
                return

            if self.storage.delete_task(task_id):
                console.print(f"✓ 任务 {task_id} 已删除", style="bold green")
            else:
                console.print(f"✗ 删除失败", style="bold red")

        except Exception as e:
            console.print(f"✗ 删除失败: {str(e)}", style="bold red")

    def show_task(self, task_id: int):
        """显示任务详情"""
        try:
            task = self.storage.get_task(task_id)
            if not task:
                console.print(f"✗ 任务 {task_id} 不存在", style="bold red")
                return

            self._display_task(task, detailed=True)

        except Exception as e:
            console.print(f"✗ 获取任务失败: {str(e)}", style="bold red")

    def search_tasks(self, keyword: str):
        """搜索任务"""
        try:
            tasks = self.storage.search_tasks(keyword)

            if not tasks:
                console.print(f"未找到包含 '{keyword}' 的任务", style="yellow")
                return

            console.print(f"\n找到 {len(tasks)} 个匹配的任务:", style="bold")
            for task in tasks:
                self._display_task(task)
                console.print()

        except Exception as e:
            console.print(f"✗ 搜索失败: {str(e)}", style="bold red")

    def show_statistics(self):
        """显示统计信息"""
        try:
            stats = self.storage.get_statistics()

            # 创建统计面板
            stats_text = f"""
📊 总任务数: {stats['total']}
✓ 已完成: {stats['completed']}
○ 待完成: {stats['pending']}
📈 完成率: {stats['completion_rate']:.1f}%

📁 按分类统计:
"""
            for cat, count in stats['by_category'].items():
                stats_text += f"  • {cat}: {count}\n"

            if stats['by_priority']:
                stats_text += "\n⚡ 待办任务优先级分布:\n"
                priority_names = {1: "低", 2: "中", 3: "高", 4: "紧急"}
                for priority, count in sorted(stats['by_priority'].items()):
                    stats_text += f"  • {priority_names.get(priority, priority)}: {count}\n"

            panel = Panel(stats_text, title="统计信息", border_style="blue")
            console.print(panel)

        except Exception as e:
            console.print(f"✗ 获取统计信息失败: {str(e)}", style="bold red")

    def _display_task(self, task: Task, detailed: bool = False):
        """显示单个任务"""
        status_icon = "✓" if task.status == TaskStatus.COMPLETED else "○"
        status_text = "已完成" if task.status == TaskStatus.COMPLETED else "待办"
        status_style = "green" if task.status == TaskStatus.COMPLETED else "yellow"

        priority_names = {1: "低", 2: "中", 3: "高", 4: "紧急"}
        priority_text = priority_names.get(task.priority.value, str(task.priority.value))

        if detailed:
            # 详细视图
            content = f"""
ID: {task.id}
标题: {task.title}
状态: {status_icon} {status_text}
分类: {task.category.value}
优先级: {priority_text}
创建时间: {task.created_at.strftime("%Y-%m-%d %H:%M:%S")}
更新时间: {task.updated_at.strftime("%Y-%m-%d %H:%M:%S")}
"""
            if task.completed_at:
                content += f"完成时间: {task.completed_at.strftime('%Y-%m-%d %H:%M:%S')}\n"

            if task.description:
                content += f"\n描述:\n{task.description}\n"

            panel = Panel(content, title=f"任务详情", border_style=status_style)
            console.print(panel)
        else:
            # 简洁视图
            console.print(
                f"[{task.id}] {status_icon} {task.title} "
                f"({task.category.value}, {priority_text})",
                style=status_style if task.status == TaskStatus.COMPLETED else "white"
            )
