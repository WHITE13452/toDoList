# 📋 TodoList CLI with AI Agent

一个功能强大的命令行待办事项管理工具，使用 **Golang** 实现，集成了 AI Agent 智能助手（支持 Qwen API）。

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
- 🎨 美观的终端界面（彩色输出、表格展示）

### AI Agent 功能
- 🤖 自然语言交互
- 🔧 智能工具调用（OpenAI Function Calling）
- 📊 任务统计与总结
- 🔄 批量操作支持
- 💡 智能建议
- 🌐 支持 Qwen API（兼容 OpenAI API）

## 🚀 快速开始

### 环境要求

- Go 1.19+
- GCC（用于编译 SQLite）

### 编译安装

```bash
# 克隆仓库
git clone https://github.com/WHITE13452/toDoList.git
cd toDoList

# 下载依赖
go mod download

# 编译
make build

# 或直接安装到 $GOPATH/bin
make install
```

### 配置 API Key（可选，用于 AI Agent 功能）

创建 `.env` 文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件，添加你的 Qwen API Key：

```env
QWEN_API_KEY=your_qwen_api_key_here
QWEN_API_BASE=https://dashscope.aliyuncs.com/compatible-mode/v1
QWEN_MODEL=qwen-plus
```

### 使用方法

#### 方式一：传统 CLI 命令

```bash
# 添加任务
./bin/todo add "完成项目文档" -d "包括架构设计和API文档" -c work -p 3

# 列出所有任务
./bin/todo list

# 列出待办任务
./bin/todo list -s pending

# 列出工作相关任务
./bin/todo list -c work

# 标记任务完成
./bin/todo complete 1

# 标记任务未完成
./bin/todo complete 1 -u

# 查看任务详情
./bin/todo show 1

# 搜索任务
./bin/todo search "项目"

# 删除任务
./bin/todo delete 1

# 跳过确认直接删除
./bin/todo delete 1 -y

# 显示统计信息
./bin/todo stats
```

#### 方式二：AI Agent 交互模式（推荐）

```bash
# 启动 AI Agent
./bin/todo chat
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
- `clear` - 清空对话历史
- `exit` / `quit` - 退出

## 📁 项目结构

```
todolist/
├── cmd/
│   └── todo/               # CLI 入口和命令
│       ├── main.go         # 主入口
│       ├── root.go         # Cobra 根命令
│       ├── add.go          # 添加命令
│       ├── list.go         # 列表命令
│       ├── complete.go     # 完成命令
│       ├── delete.go       # 删除命令
│       ├── show.go         # 详情命令
│       ├── search.go       # 搜索命令
│       ├── stats.go        # 统计命令
│       └── chat.go         # Agent 交互命令
├── internal/
│   ├── models/             # 数据模型
│   │   └── task.go
│   ├── storage/            # SQLite 持久化
│   │   └── sqlite.go
│   ├── cli/                # CLI 界面辅助
│   │   └── ui.go
│   ├── agent/              # AI Agent 核心
│   │   └── agent.go
│   └── tools/              # Agent 工具定义
│       └── tools.go
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖校验
├── Makefile                # 构建脚本
├── .env.example            # 环境变量示例
├── prd.md                  # 产品需求文档
├── ARCHITECTURE.md         # 架构设计文档
└── README.md               # 本文件
```

## 🛠️ 技术栈

- **语言**: Golang 1.19+
- **CLI 框架**: Cobra
- **终端美化**: fatih/color, tablewriter
- **持久化**: SQLite (go-sqlite3)
- **AI SDK**: go-openai（兼容 Qwen API）

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

Agent 集成了以下 9 个工具（基于 OpenAI Function Calling）：

1. `get_all_tasks` - 获取任务列表（支持过滤）
2. `add_task` - 添加新任务
3. `update_task_status` - 更新任务状态
4. `delete_task` - 删除任务
5. `search_tasks` - 搜索任务
6. `get_statistics` - 获取统计信息
7. `get_task_detail` - 获取任务详情
8. `batch_complete_tasks` - 批量完成任务
9. `batch_delete_tasks` - 批量删除任务

### 为什么使用 Qwen API？

- ✅ 国内访问速度快
- ✅ 支持中文更好
- ✅ 兼容 OpenAI API 格式
- ✅ Function Calling 支持完善
- ✅ 价格实惠

## 💡 使用建议

1. **任务分类**: 使用分类功能（work/study/life/other）来组织任务
2. **优先级管理**: 为重要任务设置高优先级，帮助你专注于关键事项
3. **使用 Agent**: Agent 模式可以更自然地管理任务，尤其是需要批量操作或复杂查询时
4. **定期查看统计**: 使用 `stats` 命令了解你的任务完成情况

## 🔒 数据存储

所有数据存储在用户家目录的 `.todolist.db` 文件中（SQLite 数据库）。

位置：`~/.todolist.db`

你也可以使用 `--db` 参数指定自定义路径：

```bash
./bin/todo --db /path/to/your/db.sqlite list
```

## 📝 示例场景

### 场景 1: 快速添加任务

```bash
./bin/todo add "买菜" -c life -p 2
./bin/todo add "写代码" -c work -p 3
./bin/todo add "学习 Golang" -c study -p 2
```

### 场景 2: 使用 Agent 进行智能管理

```bash
./bin/todo chat

你: 显示所有任务
Agent: [调用 get_all_tasks 工具，显示任务列表]

你: 哪些任务最紧急？
Agent: [分析并列出高优先级任务]

你: 帮我把所有生活相关的已完成任务删除
Agent: [调用 batch_delete_tasks 批量删除]
```

### 场景 3: 查看进度

```bash
./bin/todo stats
```

输出示例：

```
════════════════════════════════════════════════════════════
                    📊 统计信息
════════════════════════════════════════════════════════════
📋 总任务数: 10
✓ 已完成: 6
○ 待完成: 4
📈 完成率: 60.0%

📁 按分类统计:
  • work: 4
  • study: 3
  • life: 3

⚡ 待办任务优先级分布:
  • 紧急: 1
  • 高: 2
  • 中: 1
════════════════════════════════════════════════════════════
```

## 🔧 开发

### 编译项目

```bash
make build
```

### 运行测试

```bash
make test
```

### 清理构建

```bash
make clean
```

### 查看所有命令

```bash
make help
```

## 🌟 为什么用 Golang 重写？

相比 Python 版本的优势：

1. **性能更好**: 编译型语言，启动速度快
2. **单文件部署**: 编译后是独立的二进制文件，无需安装运行时
3. **跨平台**: 轻松编译到不同平台（Linux/macOS/Windows）
4. **并发支持**: Go 的 goroutine 使得未来扩展（如定时提醒）更容易
5. **类型安全**: 编译时检查，减少运行时错误
6. **生态丰富**: Cobra 等成熟的 CLI 框架

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可

MIT License
