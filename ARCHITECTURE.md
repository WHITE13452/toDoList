# 🏗️ 架构设计文档

## 系统架构

```
┌─────────────────────────────────────────────────────┐
│                    用户交互层                          │
│  ┌──────────────┐         ┌─────────────────────┐   │
│  │ CLI Commands │         │   AI Agent Chat     │   │
│  │  (Click)     │         │   (Interactive)     │   │
│  └──────────────┘         └─────────────────────┘   │
└──────────┬──────────────────────────┬───────────────┘
           │                          │
           ▼                          ▼
┌─────────────────────┐    ┌──────────────────────┐
│   CLI Controller    │    │    Agent Controller   │
│    (cli.py)         │    │     (agent.py)        │
└─────────┬───────────┘    └──────────┬───────────┘
          │                           │
          │     ┌─────────────────────┤
          │     │                     │
          ▼     ▼                     ▼
    ┌─────────────────┐     ┌──────────────────┐
    │  Business Logic │     │   Tool Executor  │
    │   (storage.py)  │     │    (tools.py)    │
    └────────┬────────┘     └────────┬─────────┘
             │                       │
             └───────────┬───────────┘
                         ▼
              ┌─────────────────────┐
              │   Data Model Layer  │
              │    (models.py)      │
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

### 1. 数据模型层 (models.py)

**职责**: 定义核心数据结构

**核心类**:
- `Task`: 任务实体
- `TaskStatus`: 状态枚举 (pending/completed)
- `TaskCategory`: 分类枚举 (work/study/life/other)
- `Priority`: 优先级枚举 (1-4)

**设计原则**:
- 使用 dataclass 简化代码
- 提供 to_dict/from_dict 方法支持序列化
- 包含业务方法 (mark_completed/mark_pending)

### 2. 持久化层 (storage.py)

**职责**: 数据持久化和 CRUD 操作

**核心类**:
- `TaskStorage`: 存储管理器

**主要方法**:
- `add_task()`: 添加任务
- `get_task()`: 获取单个任务
- `get_all_tasks()`: 获取任务列表（支持过滤）
- `update_task()`: 更新任务
- `delete_task()`: 删除任务
- `search_tasks()`: 搜索任务
- `get_statistics()`: 获取统计信息

**技术选型**:
- SQLite: 轻量级、无需配置、适合单用户场景
- 数据库位置: `~/.todolist.db`

### 3. CLI 控制器 (cli.py)

**职责**: 处理传统命令行交互

**核心类**:
- `CLI`: CLI 命令处理器

**特点**:
- 使用 Rich 库美化输出
- 支持表格、面板等丰富展示
- 提供详细和简洁两种视图

### 4. AI Agent 核心 (agent.py)

**职责**: 处理自然语言交互

**核心类**:
- `TodoAgent`: AI 助手核心

**工作流程**:
```
用户输入 → Claude API → Tool Use → 执行工具 → 返回结果 → 生成回复
```

**特点**:
- 支持多轮对话
- 自动工具调用
- 上下文管理
- 友好的中文交互

### 5. 工具层 (tools.py)

**职责**: 为 Agent 提供可调用的工具

**核心类**:
- `TodoTools`: 工具集合

**工具列表**:
1. `get_all_tasks` - 查询任务
2. `add_task` - 添加任务
3. `update_task_status` - 更新状态
4. `delete_task` - 删除任务
5. `search_tasks` - 搜索任务
6. `get_statistics` - 统计信息
7. `get_task_detail` - 任务详情
8. `batch_complete_tasks` - 批量完成
9. `batch_delete_tasks` - 批量删除

**设计特点**:
- 统一的工具定义格式（符合 Claude API 规范）
- 统一的错误处理
- 返回结构化的 JSON 结果

### 6. 主入口 (main.py)

**职责**: 应用入口和路由

**框架**: Click

**命令结构**:
```
todo
├── add       # 添加任务
├── list      # 列出任务
├── complete  # 标记完成
├── delete    # 删除任务
├── show      # 查看详情
├── search    # 搜索任务
├── stats     # 统计信息
└── chat      # AI Agent 模式
```

## 数据流

### 传统 CLI 模式

```
用户命令 → Click 解析 → CLI Controller → Storage → SQLite
                                            ↓
                     Rich 格式化输出 ← 返回结果
```

### AI Agent 模式

```
用户消息 → Agent → Claude API → Tool Use Decision
                      ↓
            Tool Executor → Storage → SQLite
                      ↓
            Tool Result → Claude API → 生成自然语言回复
                                         ↓
                              Rich 格式化输出 → 用户
```

## 设计决策

### 1. 为什么选择 SQLite？

**优势**:
- 零配置，无需单独的数据库服务
- 轻量级，适合单用户场景
- ACID 事务保证
- 比 JSON 文件更可靠和高效
- 支持 SQL 查询，便于扩展

**对比 JSON**:
- JSON: 简单但并发不安全，大数据量性能差
- MySQL: 过重，需要额外配置和服务

### 2. 为什么使用双模式（CLI + Agent）？

**互补性**:
- CLI 模式: 快速、精确、适合脚本化
- Agent 模式: 灵活、智能、适合探索和复杂操作

**场景示例**:
- 快速添加: CLI 模式
- 批量操作: Agent 模式
- 脚本集成: CLI 模式
- 智能总结: Agent 模式

### 3. Tool 设计原则

**单一职责**: 每个工具只做一件事
**幂等性**: 相同输入应产生相同结果
**错误处理**: 统一的错误返回格式
**可组合**: Agent 可以组合多个工具完成复杂任务

## 扩展性

### 可扩展点

1. **新增任务字段**
   - 修改 `models.py` 的 Task 类
   - 更新 `storage.py` 的数据库 schema
   - 添加迁移逻辑

2. **新增 Agent 工具**
   - 在 `tools.py` 中添加工具定义
   - 实现工具执行逻辑
   - 更新工具文档

3. **新增 CLI 命令**
   - 在 `main.py` 中添加命令
   - 在 `cli.py` 中实现逻辑

4. **集成其他 AI 模型**
   - 修改 `agent.py` 的 API 调用
   - 适配工具调用格式

### 性能优化方向

1. **数据库索引**: 为常用查询字段添加索引
2. **缓存**: 为统计信息添加缓存
3. **批量操作**: 优化批量插入/更新性能
4. **异步 IO**: 使用异步数据库驱动

## 安全考虑

1. **API Key 管理**: 使用环境变量，不硬编码
2. **SQL 注入**: 使用参数化查询
3. **输入验证**: 验证用户输入的合法性
4. **错误处理**: 不泄露敏感信息

## 测试策略

### 单元测试
- models.py: 数据模型逻辑
- storage.py: CRUD 操作
- tools.py: 工具执行逻辑

### 集成测试
- CLI 命令端到端测试
- Agent 对话流程测试

### 手动测试
- 用户体验测试
- Agent 智能程度评估
