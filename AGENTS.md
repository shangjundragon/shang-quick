# AGENTS.md

## 项目结构

前后端分离的全栈项目：
- **frontend/** — Vue 3 + Vite + Naive UI
- **backend/** — Go + Gin + SQLite

## 环境要求

- **Go** >= 1.23
- **Node.js** >= 18
- **pnpm** >= 8（前端必须使用 pnpm，不可用 npm）

## 前端 (frontend/)

- **包管理器**: pnpm
- **框架**: Vue 3 (`<script setup>` 语法)
- **构建工具**: Vite 8.x
- **UI 库**: Naive UI
- **状态管理**: Pinia 3.x
- **自动导入**: `unplugin-auto-import` (Vue、naive-ui API) + `unplugin-vue-components` (Naive UI 组件)
- **HTTP 请求**: Axios，基地址 `/api`，代理到后端 `localhost:8084`
- **安装依赖**: `cd frontend && pnpm install`
- **启动命令**: `cd frontend && pnpm dev`
- **构建命令**: `cd frontend && pnpm build`
- **别名**: `@` 指向 `./src`

### 前端开发要点

- **请求封装**: `src/utils/request.jsx`，自动处理 401/403/500，带 Bearer Token
- **权限指令**: `v-permission="['user:add']"`，多权限满足任一即可
- **全局消息**: `window.$message`（Naive UI 的 message）
- **路由**: 静态定义在 `src/router/index.js`，侧边栏菜单由后端 `/auth/info` 返回的 `menus` 动态渲染

## 后端 (backend/)

- **语言**: Go 1.25.5
- **框架**: Gin
- **数据库**: SQLite (使用 `github.com/glebarez/go-sqlite`)
- **ORM**: 自定义 `github.com/shangjundragon/dbw`
- **配置**: `config/config.yml` (YAML 格式)
- **端口**: `:8084`
- **安装依赖**: `cd backend && go mod tidy`
- **启动前准备**: `mkdir -p storage/logs`（必须存在，否则启动失败）
- **启动命令**: `cd backend && go run main.go`
- **首次启动**: 自动创建 `data.db` SQLite 数据库并初始化表结构和默认数据
- **数据库文件**: `data.db` (运行时生成，不需要手动创建)

### 后端开发要点

- **分层架构**: HTTP Request → Controller → Service → Model → Database
- **路由注册**: `router/router.go`，中间件执行顺序：CORS → TraceLogger → JWTAuth → Permission/OperLog
- **响应规范**: 统一格式 `{ code, message, data, trace_id }`，code 200 成功、401 未登录、403 无权限、500 错误
- **模型字段标签**:
  - `dbw:"primaryKey"` — 主键
  - `dbw:"column:xxx"` — 列映射
  - `dbw:"tableLogic"` — 逻辑删除字段
  - `dbw:"autoCreateTime:milli"` — 自动创建时间戳
  - `dbw:"autoUpdateTime:milli"` — 自动更新时间戳
  - `dbw:"createBy:true" / "updateBy:true"` — 自动填充操作人（通过 Hook）
- **事务操作**: 使用 `dbw.ExecuteTx(func(tx *sql.Tx) error { ... }, global_vars.DbConfig.Db)`
- **日志**: 开发模式 (`AppDebug: true`) 输出到控制台，生产模式写入 `storage/logs/`
- **测试**: `cd backend && go test ./...`（测试文件主要在 `pkg/` 下）

## 开发约定

- **回复语言**: 简体中文
- **数据库**: 不使用外键约束
- **代码风格**: 遵循项目现有风格，不添加无关注释
- **Git**: 不自动提交代码
- **启动**: 不主动启动项目，除非用户明确要求

## 默认账号

- 用户名：`admin`
- 密码：`admin123`

## 新增功能模块步骤

1. **后端**:
   - `model/` 定义数据模型
   - `service/` 编写业务逻辑
   - `controller/` 编写 HTTP Handler
   - `router/router.go` 注册路由和权限中间件

2. **前端**:
   - `src/views/` 创建页面组件
   - `src/api/` 定义 API 请求函数
   - 如需菜单：在「菜单管理」中添加菜单项，然后在「角色管理」分配权限
   - 按钮权限：使用 `v-permission="['xxx:add']"`

## 重要文件

- `backend/config/config.yml` — 后端配置文件（端口、数据库、JWT、日志等）
- `frontend/vite.config.js` — Vite 配置（代理、自动导入、别名）
- `backend/init.sql` — 初始建表 SQL（首次启动时执行）
- `backend/router/router.go` — 路由注册中心，定义所有 API 和中间件
- `backend/bootstrap/bootstrap.go` — 启动流程（配置加载、日志、数据库初始化）
