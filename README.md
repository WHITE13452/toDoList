# 📋 TodoList CLI with AI Agent

一个功能强大的命令行待办事项管理工具，集成了 AI Agent 智能助手。

## ✨ 特性

### 核心功能
- ✅ 添加、删除、更新待办事项
- ✅ 标记任务完成/未完成
- ✅ 查看任务列表（支持过滤）
- ✅ 任务搜索
- ✅ 统计信息展示

### 扩展功能
- 📁 任务分类（工作/学习/生活/其他）
- ⚡ 优先级管理（低/中/高/紧急）
- 💾 SQLite 持久化存储
- 🎨 美观的 Rich 终端界面

### AI Agent 功能
- 🤖 自然语言交互
- 🔧 智能工具调用
- 📊 任务统计与总结
- 🔄 批量操作支持
- 💡 智能建议

## 🚀 快速开始

### 安装依赖

```bash
pip install -r requirements.txt
```

### 配置 API Key（可选，用于 AI Agent 功能）

创建 `.env` 文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件，添加你的 Anthropic API Key：

```
ANTHROPIC_API_KEY=your_api_key_here
```

### 使用方法

#### 方式一：传统 CLI 命令

```bash
# 添加任务
python -m src.main add "完成项目文档" -d "包括架构设计和API文档" -c work -p 3

# 列出所有任务
python -m src.main list

# 列出待办任务
python -m src.main list -s pending

# 列出工作相关任务
python -m src.main list -c work

# 标记任务完成
python -m src.main complete 1

# 查看任务详情
python -m src.main show 1

# 搜索任务
python -m src.main search "项目"

# 删除任务
python -m src.main delete 1

# 显示统计信息
python -m src.main stats
```

#### 方式二：AI Agent 交互模式（推荐）

```bash
# 启动 AI Agent
python -m src.main chat
```

在 Agent 模式下，你可以用自然语言交互：

```
你: 显示所有未完成的任务
你: 帮我添加一个任务：准备周会演示
你: 把任务 3 标记为完成
你: 有哪些高优先级的工作任务？
你: 给我一个总结
你: 批量完成所有学习相关的任务
```

快捷命令：
- `list` / `ls` - 显示所有任务
- `stats` - 显示统计信息
- `help` - 显示帮助
- `exit` / `quit` - 退出

## 📁 项目结构

```
todolist/
├── src/
│   ├── __init__.py      # 包初始化
│   ├── main.py          # CLI 主入口
│   ├── models.py        # 数据模型定义
│   ├── storage.py       # SQLite 持久化层
│   ├── cli.py           # CLI 命令实现
│   ├── agent.py         # AI Agent 核心
│   └── tools.py         # Agent 工具定义
├── requirements.txt     # Python 依赖
├── .env.example         # 环境变量示例
├── prd.md              # 产品需求文档
└── README.md           # 本文件
```

## 🛠️ 技术栈

- **语言**: Python 3.8+
- **CLI 框架**: Click
- **终端界面**: Rich
- **持久化**: SQLite
- **AI**: Anthropic Claude API

## 📊 数据模型

每个任务包含以下字段：

- `id`: 任务 ID（自动生成）
- `title`: 任务标题（必填）
- `description`: 任务描述（可选）
- `status`: 状态（pending/completed）
- `category`: 分类（work/study/life/other）
- `priority`: 优先级（1-4）
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `completed_at`: 完成时间

## 🤖 AI Agent 能力

Agent 集成了以下工具：

1. `get_all_tasks` - 获取任务列表
2. `add_task` - 添加新任务
3. `update_task_status` - 更新任务状态
4. `delete_task` - 删除任务
5. `search_tasks` - 搜索任务
6. `get_statistics` - 获取统计信息
7. `get_task_detail` - 获取任务详情
8. `batch_complete_tasks` - 批量完成任务
9. `batch_delete_tasks` - 批量删除任务

## 💡 使用建议

1. **任务分类**: 使用分类功能（work/study/life/other）来组织任务
2. **优先级管理**: 为重要任务设置高优先级，帮助你专注于关键事项
3. **使用 Agent**: Agent 模式可以更自然地管理任务，尤其是需要批量操作或复杂查询时
4. **定期查看统计**: 使用 `stats` 命令了解你的任务完成情况

## 🔒 数据存储

所有数据存储在用户家目录的 `.todolist.db` 文件中（SQLite 数据库）。

位置：`~/.todolist.db`

## 📝 示例场景

### 场景 1: 快速添加任务

```bash
python -m src.main add "买菜" -c life -p 2
python -m src.main add "写代码" -c work -p 3
python -m src.main add "学习 Python" -c study -p 2
```

### 场景 2: 使用 Agent 进行智能管理

```bash
python -m src.main chat

你: 显示所有任务
Agent: [显示任务列表]

你: 哪些任务最紧急？
Agent: [分析并列出高优先级任务]

你: 帮我把所有生活相关的已完成任务删除
Agent: [批量删除操作]
```

### 场景 3: 查看进度

```bash
python -m src.main stats
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可

MIT License