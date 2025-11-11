"""AI Agent 核心实现"""

import os
from typing import List, Dict, Any, Optional
import anthropic
from rich.console import Console
from rich.markdown import Markdown
from rich.panel import Panel

from .tools import TodoTools

console = Console()


class TodoAgent:
    """TodoList AI Agent"""

    def __init__(self, tools: TodoTools, api_key: Optional[str] = None):
        """初始化 Agent

        Args:
            tools: TodoTools 实例
            api_key: Anthropic API Key，如果不提供则从环境变量读取
        """
        self.tools = tools
        self.api_key = api_key or os.getenv("ANTHROPIC_API_KEY")

        if not self.api_key:
            raise ValueError(
                "未找到 ANTHROPIC_API_KEY。请在环境变量中设置或创建 .env 文件。"
            )

        self.client = anthropic.Anthropic(api_key=self.api_key)
        self.conversation_history: List[Dict[str, Any]] = []

        # 系统提示词
        self.system_prompt = """你是一个智能待办事项管理助手。你可以帮助用户管理他们的任务列表。

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
- 如果任务很多，可以先总结再列出重点
"""

    def chat(self, user_message: str) -> str:
        """与 Agent 对话

        Args:
            user_message: 用户消息

        Returns:
            Agent 的回复
        """
        # 添加用户消息到历史
        self.conversation_history.append({
            "role": "user",
            "content": user_message
        })

        try:
            # 调用 Claude API
            response = self.client.messages.create(
                model="claude-3-5-sonnet-20241022",
                max_tokens=4096,
                system=self.system_prompt,
                tools=self.tools.get_tools_definition(),
                messages=self.conversation_history
            )

            # 处理响应
            return self._process_response(response)

        except Exception as e:
            error_msg = f"Agent 错误: {str(e)}"
            console.print(error_msg, style="bold red")
            return error_msg

    def _process_response(self, response: anthropic.types.Message) -> str:
        """处理 API 响应

        Args:
            response: Claude API 响应

        Returns:
            最终的文本回复
        """
        assistant_message = {
            "role": "assistant",
            "content": []
        }

        final_text = ""

        # 处理响应内容
        while response.stop_reason == "tool_use":
            # 收集所有内容块
            for content_block in response.content:
                assistant_message["content"].append(content_block.model_dump())

                if content_block.type == "text":
                    final_text += content_block.text
                elif content_block.type == "tool_use":
                    # 执行工具调用
                    tool_name = content_block.name
                    tool_input = content_block.input
                    tool_use_id = content_block.id

                    console.print(
                        f"[dim]🔧 调用工具: {tool_name}...[/dim]"
                    )

                    # 执行工具
                    tool_result = self.tools.execute_tool(tool_name, tool_input)

                    # 添加工具结果
                    assistant_message["content"].append({
                        "type": "tool_result",
                        "tool_use_id": tool_use_id,
                        "content": str(tool_result)
                    })

            # 添加助手消息到历史
            self.conversation_history.append(assistant_message)

            # 继续对话以获取最终响应
            response = self.client.messages.create(
                model="claude-3-5-sonnet-20241022",
                max_tokens=4096,
                system=self.system_prompt,
                tools=self.tools.get_tools_definition(),
                messages=self.conversation_history
            )

            # 重置 assistant_message 为新的轮次
            assistant_message = {
                "role": "assistant",
                "content": []
            }

        # 收集最终响应的文本
        for content_block in response.content:
            assistant_message["content"].append(content_block.model_dump())
            if content_block.type == "text":
                final_text += content_block.text

        # 添加最终消息到历史
        self.conversation_history.append(assistant_message)

        return final_text

    def start_interactive_session(self):
        """启动交互式对话会话"""
        console.print(Panel(
            "[bold cyan]TodoList AI Agent[/bold cyan]\n\n"
            "我是你的智能待办助手，可以帮你管理任务。\n\n"
            "你可以问我：\n"
            "• 'list' 或 '显示所有任务'\n"
            "• '统计' 或 '总结一下'\n"
            "• '添加任务：写周报'\n"
            "• '完成任务 1'\n"
            "• '搜索包含会议的任务'\n"
            "• 或者用自然语言描述你想做什么\n\n"
            "输入 'exit' 或 'quit' 退出。",
            border_style="cyan"
        ))

        while True:
            try:
                # 获取用户输入
                user_input = console.input("\n[bold green]你:[/bold green] ")

                if not user_input.strip():
                    continue

                # 检查退出命令
                if user_input.lower() in ['exit', 'quit', '退出', 'q']:
                    console.print("\n[cyan]再见！[/cyan]")
                    break

                # 处理快捷命令
                if user_input.lower() in ['list', 'ls', '列表', '显示']:
                    user_input = "显示所有待办任务"
                elif user_input.lower() in ['stats', 'statistics', '统计']:
                    user_input = "显示统计信息和总结"
                elif user_input.lower() in ['help', 'h', '帮助']:
                    console.print(Panel(
                        "可用命令：\n"
                        "• list/ls - 显示所有任务\n"
                        "• stats - 显示统计信息\n"
                        "• help - 显示此帮助\n"
                        "• exit - 退出\n\n"
                        "或者直接用自然语言描述你想做什么，例如：\n"
                        "• '帮我添加一个任务：准备项目演示'\n"
                        "• '完成任务 3'\n"
                        "• '有哪些工作相关的未完成任务？'\n",
                        title="帮助",
                        border_style="blue"
                    ))
                    continue

                # 发送消息给 Agent
                console.print()
                response = self.chat(user_input)

                # 显示 Agent 回复
                console.print(
                    Panel(
                        Markdown(response),
                        title="[bold cyan]Agent[/bold cyan]",
                        border_style="cyan"
                    )
                )

            except KeyboardInterrupt:
                console.print("\n\n[cyan]再见！[/cyan]")
                break
            except EOFError:
                console.print("\n\n[cyan]再见！[/cyan]")
                break
            except Exception as e:
                console.print(f"\n[bold red]错误: {str(e)}[/bold red]")

    def clear_history(self):
        """清空对话历史"""
        self.conversation_history = []
        console.print("[dim]对话历史已清空[/dim]")
