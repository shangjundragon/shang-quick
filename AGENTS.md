# AGENTS.md

## 项目结构

前后端分离的全栈项目：
- **frontend/** — Vue 3 + Vite 8 + Naive UI
- **backend/** — Go 1.25.5 + Gin + SQLite + 自定义 ORM `dbw`

## 环境要求

- **Go** >= 1.23, **Node.js** >= 18, **pnpm** >= 8

## 快速命令

| 操作 | 命令 |
|------|------|
| 后端安装依赖 | `cd backend && go mod tidy` |
| 后端启动 | `cd backend && go run main.go`（需先 `mkdir storage/logs`） |
| 后端测试 | `cd backend && go test ./...` |
| 前端安装依赖 | `cd frontend && pnpm install` |
| 前端开发 | `cd frontend && pnpm dev`（端口 5177，代理 `/api` → `localhost:8084`） |
| 前端构建 | `cd frontend && pnpm build` |

## 前端 (frontend/)

- **包管理器**: pnpm（不可用 npm）
- **框架**: Vue 3 `<script setup>` 语法
- **UI 库**: Naive UI（组件自动导入，API 如 `useMessage`/`useDialog` 自动导入）
- **JSX 支持**: `@vitejs/plugin-vue-jsx` 已启用（`request.jsx` 使用 JSX 扩展名）
- **状态管理**: Pinia 3.x，用户状态在 `store/user.js`
- **别名**: `@` → `./src`
- **HTTP 请求**: Axios 封装在 `src/utils/request.jsx`，baseURL `/api`，自动 Bearer Token、处理 401/403/500
- **权限指令**: `v-permission="['user:add']"`（多权限满足任一即可）
- **全局消息**: `window.$message`（Naive UI message）
- **路由**: 静态定义在 `src/router/index.js`，侧边栏菜单由后端 `/auth/info` 返回的 `menus` 动态渲染
- **可用第三方库**: `lodash-es`, `date-fns`, `uuid`, `@vicons/ionicons5`

## 后端 (backend/)

- **框架**: Gin
- **数据库**: SQLite（`github.com/glebarez/go-sqlite`），运行后生成 `data.db`
- **ORM**: 自定义 `github.com/shangjundragon/dbw` v1.2.7
- **分层**: Controller → Service → Model → Database，中间件在路由层注入
- **配置**: `config/config.yml`（YAML，Viper 加载）
- **中间件顺序**: CORS → TraceLogger → JWTAuth → Permission/OperLog
- **响应格式**: `{ code, message, data, trace_id }`（200 成功, 401 未登录, 403 无权限, 500 错误）
- **事务**: `dbw.ExecuteTx(func(tx *sql.Tx) error { ... }, global_vars.DbConfig.Db)`
- **日志**: `AppDebug: true` 输出控制台，`false` 写入 `storage/logs/`（注意：`config.yml` 中路径前缀为 `/storage/logs/`，但实际拼接 `BasePath`）

### Model 字段标签

| 标签 | 用途 |
|------|------|
| `dbw:"primaryKey"` | 主键 |
| `dbw:"column:xxx"` | 列映射 |
| `dbw:"column:xxx;default:1"` | 默认值 |
| `dbw:"tableLogic"` | 逻辑删除字段（`del_flag`） |
| `dbw:"autoCreateTime:milli"` | 自动创建时间戳 |
| `dbw:"autoUpdateTime:milli"` | 自动更新时间戳 |
| `dbw:"createBy:true" / "updateBy:true"` | 自动填充操作人（Hook） |

## 功能模块（含未在菜单管理预置的）

| 模块 | 权限标识 | 说明 |
|------|----------|------|
| 用户管理 | `user:list/add/edit/delete/resetPwd` | RBAC 核心 |
| 角色管理 | `role:list/add/edit/delete/assign` | 角色-菜单分配 |
| 菜单管理 | `menu:list/add/edit/delete` | 三级类型：目录/菜单/按钮 |
| 部门管理 | `dept:list/add/edit/delete` | 树形 |
| 操作日志 | 仅需登录 | 自动记录 POST 写操作 |
| 文件管理 | `file:list/upload/delete` | 上传/列表/删除 |
| 在线用户 | `onlineUser:list/kick` | 查看/踢下线 |

## API 路由

所有路由定义在 `backend/router/router.go`，以 `/api/v1/` 为前缀。增删改路由使用 `middleware.OperLog` + `middleware.Permission` 中间件链。**以 router.go 为唯一真实来源**。

## 开发约定

- **数据库**: 不使用外键约束
- **回复语言**: 简体中文
- **代码**: 添加必要且简洁的注释（Go 包级注释、导出函数文档、复杂逻辑内联说明、前端关键文件），follow existing conventions
- **Git**: 不自动提交代码，不主动启动项目
- **默认账号**: `admin` / `admin123`

## 新增功能模块步骤

1. **后端**: `model/` → `service/` → `controller/` → `router/router.go` 注册路由+权限
2. **前端**: `src/api/` 定义请求 → `src/views/` 创建页面 → 如需菜单：在「菜单管理」添加菜单项 → 「角色管理」分配权限
3. **按钮权限**: `v-permission="['xxx:add']"`

## 重要文件

- `backend/config/config.yml` — 端口、JWT Secret、数据库、日志、CORS、上传配置
- `backend/router/router.go` — 所有 API 和中间件定义（唯一真实来源）
- `backend/bootstrap/bootstrap.go` — 启动流程（配置加载、日志、数据库初始化）
- `backend/init.sql` — 初始建表 SQL
- `frontend/vite.config.js` — 代理、自动导入、别名、JSX 插件
- `frontend/src/utils/request.jsx` — Axios 请求封装

## 配置文件注意点

- `config.yml` 中 CORS 段（`AllowOrigins`、`AllowCredentials`）和 Upload 段（`MaxSize`、`AllowedExtensions` 等）未在数据库或代码注释中说明，需直接查看配置文件
- `Jwt.Secret` 生产环境必须修改为 >= 32 字符（bootstrap 会校验）
