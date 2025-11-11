# 🏗️ 架构设计文档 (Golang 实现)

## 系统架构

```
┌─────────────────────────────────────────────────────┐
│                    用户交互层                          │
│  ┌──────────────┐         ┌─────────────────────┐   │
│  │ CLI Commands │         │   AI Agent Chat     │   │
│  │   (Cobra)    │         │   (Interactive)     │   │
│  └──────────────┘         └─────────────────────┘   │
└──────────┬──────────────────────────┬───────────────┘
           │                          │
           ▼                          ▼
┌─────────────────────┐    ┌──────────────────────┐
│   CLI Controller    │    │    Agent Controller   │
│  (cmd/todo/*.go)    │    │   (agent/agent.go)   │
└─────────┬───────────┘    └──────────┬───────────┘
          │                           │
          │     ┌─────────────────────┤
          │     │                     │
          ▼     ▼                     ▼
    ┌─────────────────┐     ┌──────────────────┐
    │  Business Logic │     │   Tool Executor  │
    │ (storage/*.go)  │     │  (tools/*.go)    │
    └────────┬────────┘     └────────┬─────────┘
             │                       │
             └───────────┬───────────┘
                         ▼
              ┌─────────────────────┐
              │   Data Model Layer  │
              │   (models/*.go)     │
              └──────────┬──────────┘
                         │
                         ▼
              ┌─────────────────────┐
              │   Persistence Layer │
              │   SQLite Database   │
              │   (~/.todolist.db)  │
              └─────────────────────┘
```

## 模块说明

### 1. 数据模型层 (internal/models/task.go)

**职责**: 定义核心数据结构

**核心类型**:
- `Task`: 任务实体结构体
- `TaskStatus`: 状态类型 (pending/completed)
- `TaskCategory`: 分类类型 (work/study/life/other)
- `Priority`: 优先级类型 (1-4)
- `Statistics`: 统计信息结构体

**设计特点**:
- 使用 Go struct 定义清晰的数据模型
- JSON 标签支持序列化/反序列化
- 方法：`MarkCompleted()`, `MarkPending()`, `NewTask()`
- 类型安全的枚举常量

### 2. 持久化层 (internal/storage/sqlite.go)

**职责**: 数据持久化和 CRUD 操作

**核心结构**:
- `Storage`: 存储管理器，封装 SQLite 数据库操作

**主要方法**:
- `New()`: 创建存储实例，初始化数据库
- `AddTask()`: 添加任务
- `GetTask()`: 获取单个任务
- `GetAllTasks()`: 获取任务列表（支持过滤）
- `UpdateTask()`: 更新任务
- `DeleteTask()`: 删除任务
- `SearchTasks()`: 搜索任务
- `GetStatistics()`: 获取统计信息
- `Close()`: 关闭数据库连接

**技术选型**:
- SQLite (github.com/mattn/go-sqlite3): 轻量级、无需配置
- 数据库位置: `~/.todolist.db`
- 使用索引优化查询性能
- sql.NullTime 处理可选时间字段

### 3. CLI 界面层 (internal/cli/ui.go)

**职责**: 提供美观的终端界面

**主要函数**:
- `PrintSuccess()`: 打印成功消息（绿色）
- `PrintError()`: 打印错误消息（红色）
- `PrintInfo()`: 打印信息（青色）
- `PrintTask()`: 打印单个任务（支持简洁/详细模式）
- `PrintTaskTable()`: 表格形式打印任务列表
- `PrintStatistics()`: 打印统计信息
- `PrintAgentWelcome()`: Agent 欢迎界面
- `PrintAgentResponse()`: Agent 响应展示

**使用的库**:
- `fatih/color`: 彩色输出
- `olekukonko/tablewriter`: 表格展示

### 4. AI Agent 核心 (internal/agent/agent.go)

**职责**: 处理自然语言交互

**核心结构**:
- `Agent`: AI 助手核心
- `Config`: Agent 配置（API Key, Base URL, Model）

**工作流程**:
```
用户输入 → Chat() → OpenAI API → Tool Use?
                                    ↓ Yes
                          ExecuteTool() → 返回结果
                                    ↓
                          继续对话 → 最终文本响应
```

**主要方法**:
- `New()`: 创建 Agent 实例
- `Chat()`: 处理用户消息，返回 AI 响应
- `ClearHistory()`: 清空对话历史
- `GetHistory()`: 获取对话历史

**API 集成**:
- 使用 `sashabaranov/go-openai` SDK
- 支持 Qwen API（OpenAI 兼容）
- Function Calling 实现工具调用

### 5. 工具层 (internal/tools/tools.go)

**职责**: 为 Agent 提供可调用的工具

**核心结构**:
- `TodoTools`: 工具集合，封装所有可用工具

**工具列表**（9 个）:
1. `get_all_tasks` - 查询任务（支持过滤）
2. `add_task` - 添加任务
3. `update_task_status` - 更新状态
4. `delete_task` - 删除任务
5. `search_tasks` - 搜索任务
6. `get_statistics` - 统计信息
7. `get_task_detail` - 任务详情
8. `batch_complete_tasks` - 批量完成
9. `batch_delete_tasks` - 批量删除

**主要方法**:
- `GetToolDefinitions()`: 返回 OpenAI Function Calling 格式的工具定义
- `ExecuteTool()`: 执行指定工具，返回 JSON 结果

**设计特点**:
- 统一的工具定义格式（OpenAI 标准）
- JSON Schema 参数验证
- 统一的错误处理和结果格式

### 6. 命令层 (cmd/todo/*.go)

**职责**: CLI 命令定义和路由

**框架**: Cobra

**命令结构**:
```
todo (root.go)
├── add       (add.go)       # 添加任务
├── list      (list.go)      # 列出任务
├── complete  (complete.go)  # 标记完成
├── delete    (delete.go)    # 删除任务
├── show      (show.go)      # 查看详情
├── search    (search.go)    # 搜索任务
├── stats     (stats.go)     # 统计信息
└── chat      (chat.go)      # AI Agent 模式
```

**特点**:
- 每个命令独立文件，清晰易维护
- 统一的错误处理
- 持久化 flags（如 --db）
- PreRun/PostRun 钩子管理资源

## 数据流

### 传统 CLI 模式

```
用户命令 → Cobra 解析 → 命令 Handler
                           ↓
                    Storage 操作 → SQLite
                           ↓
                    UI 格式化 ← 返回结果
                           ↓
                        终端输出
```

### AI Agent 模式

```
用户消息 → Agent.Chat() → OpenAI API
                             ↓
               Tool Use Decision → tools.ExecuteTool()
                             ↓
                    Storage 操作 → SQLite
                             ↓
               Tool Result → OpenAI API → 生成自然语言
                                           ↓
                                 UI 格式化 → 终端输出
```

## 设计决策

### 1. 为什么使用 Golang？

**优势**:
- **性能**: 编译型语言，启动快、运行快
- **部署**: 单个二进制文件，无需运行时
- **跨平台**: 一次编译，到处运行
- **并发**: 原生 goroutine 支持，便于扩展
- **类型安全**: 编译时检查，减少运行时错误
- **生态**: Cobra、go-sqlite3 等成熟库

**对比 Python**:
- Python: 需要运行时，依赖管理复杂
- Go: 单文件部署，依赖编译进二进制

### 2. 为什么选择 SQLite？

**优势**:
- 零配置，无需单独的数据库服务
- 轻量级，适合单用户场景
- ACID 事务保证
- 比 JSON 文件更可靠和高效
- 支持 SQL 查询，便于扩展
- 跨平台兼容性好

**对比其他方案**:
- JSON: 简单但并发不安全，大数据量性能差
- MySQL/PostgreSQL: 过重，需要额外配置和服务

### 3. 为什么使用 Qwen API？

**优势**:
- 国内访问速度快，无需代理
- 中文理解能力强
- 兼容 OpenAI API 格式
- Function Calling 支持完善
- 价格实惠

**实现方式**:
- 使用 `go-openai` SDK
- 只需修改 BaseURL 即可切换到 Qwen
- 工具定义完全兼容 OpenAI 格式

### 4. 项目结构设计

采用标准的 Go 项目布局：
- `cmd/`: 应用程序入口
- `internal/`: 私有应用代码（不对外暴露）
- `internal/models/`: 数据模型
- `internal/storage/`: 存储层
- `internal/agent/`: AI Agent
- `internal/tools/`: 工具层
- `internal/cli/`: UI 辅助

**优点**:
- 清晰的模块划分
- 符合 Go 社区标准
- 易于测试和维护

## 扩展性

### 可扩展点

1. **新增任务字段**
   - 修改 `models/task.go` 的 Task 结构
   - 更新 `storage/sqlite.go` 的数据库 schema
   - 添加数据库迁移逻辑

2. **新增 Agent 工具**
   - 在 `tools/tools.go` 中添加工具定义
   - 实现工具执行逻辑
   - OpenAI API 会自动识别新工具

3. **新增 CLI 命令**
   - 在 `cmd/todo/` 创建新的命令文件
   - 使用 Cobra 注册命令
   - 调用已有的存储层接口

4. **切换 AI 模型**
   - 修改 `agent/agent.go` 的 API 调用
   - 支持多种兼容 OpenAI 格式的 API
   - 只需修改配置即可切换

5. **添加提醒功能**
   - 使用 goroutine 实现后台定时器
   - 利用 Go 的并发优势
   - 可考虑使用 cron 库

### 性能优化方向

1. **数据库优化**
   - 为常用查询字段添加索引（已实现）
   - 使用预编译语句（prepared statements）
   - 批量操作优化

2. **Agent 优化**
   - 缓存常用查询结果
   - 实现流式输出
   - 并发处理多个工具调用

3. **编译优化**
   - 使用 `-ldflags="-s -w"` 减小二进制大小
   - UPX 压缩（可选）
   - 交叉编译支持多平台

## 安全考虑

1. **API Key 管理**: 使用环境变量，不硬编码
2. **SQL 注入**: 使用参数化查询（已实现）
3. **输入验证**: 验证用户输入的合法性
4. **错误处理**: 不泄露敏感信息
5. **文件权限**: 数据库文件设置合适的权限

## 测试策略

### 单元测试
- `models`: 数据模型逻辑
- `storage`: CRUD 操作
- `tools`: 工具执行逻辑

### 集成测试
- CLI 命令端到端测试
- Agent 对话流程测试
- 数据库迁移测试

### 手动测试
- 用户体验测试
- Agent 智能程度评估
- 跨平台兼容性测试

## 依赖管理

使用 Go Modules（go.mod）管理依赖：

```go
module github.com/WHITE13452/toDoList

go 1.19

require (
    github.com/spf13/cobra       # CLI 框架
    github.com/mattn/go-sqlite3   # SQLite 驱动
    github.com/sashabaranov/go-openai  # OpenAI/Qwen API
    github.com/joho/godotenv      # .env 文件支持
    github.com/fatih/color        # 彩色输出
    github.com/olekukonko/tablewriter  # 表格展示
)
```

## 构建和部署

### 本地开发
```bash
make build     # 编译
make test      # 测试
make clean     # 清理
```

### 生产部署
```bash
make install   # 安装到 $GOPATH/bin
```

### 跨平台编译
```bash
# Linux
GOOS=linux GOARCH=amd64 make build

# macOS
GOOS=darwin GOARCH=amd64 make build

# Windows
GOOS=windows GOARCH=amd64 make build
```

## 未来规划

1. **Web 界面**: 使用 Gin 框架提供 REST API 和 Web UI
2. **多用户支持**: 添加用户认证和数据隔离
3. **云同步**: 支持多设备数据同步
4. **插件系统**: 允许第三方扩展功能
5. **移动端**: 开发配套的移动应用
6. **提醒通知**: 集成桌面通知和邮件提醒
7. **数据分析**: 提供更丰富的统计和可视化
