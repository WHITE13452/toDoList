#!/usr/bin/env python3
"""TodoList CLI 主入口"""

import click
import os
from dotenv import load_dotenv

from .storage import TaskStorage
from .cli import CLI
from .tools import TodoTools
from .agent import TodoAgent

# 加载环境变量
load_dotenv()


@click.group()
@click.version_option(version="1.0.0")
def todo():
    """
    📋 TodoList - 智能待办事项管理工具

    支持传统 CLI 命令和 AI Agent 交互两种模式。
    """
    pass


@todo.command()
@click.argument('title')
@click.option('-d', '--description', help='任务描述')
@click.option('-c', '--category',
              type=click.Choice(['work', 'study', 'life', 'other']),
              default='other',
              help='任务分类')
@click.option('-p', '--priority',
              type=click.IntRange(1, 4),
              default=2,
              help='优先级：1(低) 2(中) 3(高) 4(紧急)')
def add(title, description, category, priority):
    """添加新任务"""
    storage = TaskStorage()
    cli = CLI(storage)
    cli.add_task(title, description, category, priority)


@todo.command()
@click.option('-s', '--status',
              type=click.Choice(['pending', 'completed']),
              help='按状态过滤')
@click.option('-c', '--category',
              type=click.Choice(['work', 'study', 'life', 'other']),
              help='按分类过滤')
@click.option('-a', '--all', 'show_all',
              is_flag=True,
              help='显示所有任务')
def list(status, category, show_all):
    """列出任务"""
    storage = TaskStorage()
    cli = CLI(storage)
    cli.list_tasks(status, category, show_all)


@todo.command()
@click.argument('task_id', type=int)
@click.option('-u', '--uncomplete',
              is_flag=True,
              help='标记为未完成')
def complete(task_id, uncomplete):
    """标记任务完成/未完成"""
    storage = TaskStorage()
    cli = CLI(storage)
    cli.complete_task(task_id, uncomplete)


@todo.command()
@click.argument('task_id', type=int)
@click.option('-y', '--yes',
              is_flag=True,
              help='跳过确认')
def delete(task_id, yes):
    """删除任务"""
    storage = TaskStorage()
    cli = CLI(storage)

    if not yes:
        if not click.confirm(f'确定要删除任务 {task_id} 吗？'):
            click.echo('已取消')
            return

    cli.delete_task(task_id)


@todo.command()
@click.argument('task_id', type=int)
def show(task_id):
    """显示任务详情"""
    storage = TaskStorage()
    cli = CLI(storage)
    cli.show_task(task_id)


@todo.command()
@click.argument('keyword')
def search(keyword):
    """搜索任务"""
    storage = TaskStorage()
    cli = CLI(storage)
    cli.search_tasks(keyword)


@todo.command()
def stats():
    """显示统计信息"""
    storage = TaskStorage()
    cli = CLI(storage)
    cli.show_statistics()


@todo.command()
@click.option('--api-key',
              envvar='ANTHROPIC_API_KEY',
              help='Anthropic API Key（或通过环境变量设置）')
def chat(api_key):
    """
    启动 AI Agent 交互模式

    在这个模式下，你可以用自然语言与 AI 助手对话来管理任务。

    示例：
    • "显示所有未完成的任务"
    • "帮我添加一个任务：准备项目演示"
    • "完成任务 3"
    • "有哪些工作相关的任务？"
    • "给我一个总结"
    """
    try:
        storage = TaskStorage()
        tools = TodoTools(storage)
        agent = TodoAgent(tools, api_key=api_key)
        agent.start_interactive_session()
    except ValueError as e:
        click.echo(f"错误: {str(e)}", err=True)
        click.echo("\n请确保设置了 ANTHROPIC_API_KEY 环境变量。", err=True)
        click.echo("你可以创建一个 .env 文件并添加：", err=True)
        click.echo("ANTHROPIC_API_KEY=your_api_key_here", err=True)
    except Exception as e:
        click.echo(f"启动 Agent 失败: {str(e)}", err=True)


def main():
    """主函数"""
    todo()


if __name__ == '__main__':
    main()
