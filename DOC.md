# TodoList CLI with AI Agent 项目说明文档

## 1. 技术选型

* **编程语言**: Go 1.19+

  * **理由**: 一直在用，轻量
* **CLI 框架**: Cobra

  * **理由**: Go 社区最流行的 CLI 框架,支持子命令和参数管理
* **数据库**: SQLite

  * **理由**: 轻量级,无需单独部署;数据存储在单个文件中,便于备份
* **AI 集成**: go-openai SDK + 阿里云通义千问

  * **理由**: 兼容 OpenAI API,国内访问速度快,中文支持好

## 2. 项目结构设计

### 整体架构

采用三层架构: CLI 命令层 → 业务逻辑层 → 数据存储层

### 目录结构

```
toDoList/
├── cmd/todo/              # CLI 命令入口
│   ├── main.go           # 主入口
│   ├── root.go           # Cobra 根命令
│   ├── add.go            # 添加任务
│   ├── list.go           # 列表
│   ├── delete.go         # 删除
│   ├── complete.go       # 完成
│   ├── search.go         # 搜索
│   ├── stats.go          # 统计
│   └── chat.go           # AI Agent
│
├── internal/             # 内部实现
│   ├── models/           # 数据模型
│   │   └── task.go
│   ├── storage/          # SQLite 封装
│   │   └── sqlite.go
│   ├── agent/            # AI Agent
│   │   └── agent.go
│   ├── tools/            # Agent 工具集
│   │   └── tools.go
│   └── cli/              # 界面美化
│       └── ui.go
│
├── .env.example          # 环境变量示例
├── Makefile             # 构建脚本
├── go.mod
└── README.md
```

### 模块职责

- **cmd/todo/**: 解析用户输入,调用业务逻辑
- **internal/storage/**: 封装所有数据库操作
- **internal/agent/**: AI Agent 核心逻辑
- **internal/tools/**: Agent 可调用的 9 个工具
- **internal/cli/**: 终端美化输出

## 3. 需求细节与决策

### 基础功能

- **任务描述**: 标题必填,描述可选
- **空输入处理**: 使用 Cobra 的参数验证
- **已完成任务**: 用不同颜色标识(✅绿色/⏳黄色)
- **排序逻辑**: 默认按优先级+创建时间降序,支持 `--sort` 参数自定义

### AI Agent 设计

- **模式选择**: ReAct 模式(支持真实工具调用)
- **工具调用**: 基于 OpenAI Function Calling,最多循环 10 次
- **防死循环**: 超过 10 次工具调用自动终止

### 扩展功能

- **分类**: 工作/学习/生活/其他
- **优先级**: 1-5 级(数值越大越重要)
- **批量操作**: 支持批量完成/删除
- **搜索**: SQLite LIKE 查询

## 4. AI 使用说明

### 使用的 AI 工具

**Claude Code (Anthropic)** - 用于项目初期构建

### AI 输出的修改

- 支持按照关键字删除，弱化taskID
- 增加排序逻辑

## 5. 运行与测试方式

### 环境要求

- Go 1.19+
- GCC (编译 SQLite)

### 本地运行

```bash
# 1. 克隆代码
git clone https://github.com/WHITE13452/toDoList.git
cd toDoList

# 2. 下载依赖
go mod download

# 3. 配置 API Key (可选)
cp .env.example .env
# 编辑 .env 添加 QWEN_API_KEY

# 4. 编译运行
make build
./bin/todo --help
```

### 功能测试

```bash
# 添加任务
todo add "学习 Go" --priority 1

# 列出任务
todo list --sort priority

# AI Agent
todo chat
```

### 已测试环境

- ✅ macOS 13.5 (M1) - Go 1.21
- ✅ Ubuntu 22.04 - Go 1.20

### 已知问题

1. 默认数据库路径 `~/.todolist.db`,可用 `--db` 自定义
2. Agent 调用较慢(3-5秒)

## 6. 总结与反思

### 项目亮点

1. **架构清晰**: 严格的三层分离设计
2. **AI Agent 创新**: 自然语言交互的 CLI 工具,支持真实工具调用
3. **用户体验**: 彩色输出、表格展示、关键词搜索
4. **代码质量**: 完善的错误处理和统一的命名规范

### 如果有更多时间

- 增加单元测试
- agent流式输出的优化，优化用户体验，当前回答需要等待过长时间
- agent并发调用没有依赖的tools，加快react的过程
- 加入缓存，减少重复查询的可能
- 支持任务导出(JSON/CSV)

最大亮点

**🌟 自然语言交互 + 完善的工具调用系统**

传统方式需要记住复杂的命令和参数,Agent 模式让用户用自然语言完成所有操作,这是对传统 CLI 工具的创新。
