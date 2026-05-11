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
- **构建工具**: Vite
- **UI 库**: Naive UI
- **状态管理**: Pinia 3.x
- **自动导入**: `unplugin-auto-import` (Vue、naive-ui API) + `unplugin-vue-components` (Naive UI 组件)
- **HTTP 请求**: Axios，基地址 `/api`，代理到后端 `localhost:8084`
- **安装依赖**: `cd frontend && pnpm install`
- **启动命令**: `cd frontend && pnpm dev`
- **构建命令**: `cd frontend && pnpm build`
- **别名**: `@` 指向 `./src`

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

## 开发约定

- **回复语言**: 简体中文
- **数据库**: 不使用外键约束
- **日志**: 开发模式 (`AppDebug: true`) 输出到控制台，生产模式写入文件
- **代码风格**: 遵循项目现有风格，不添加无关注释
- **Git**: 不自动提交代码
- **启动**: 不主动启动项目，除非用户明确要求
- **Go 测试**: `cd backend && go test ./...`（测试文件少，主要在 `pkg/` 下）

## 默认账号

- 用户名：`admin`
- 密码：`admin123`

## 重要文件

- `backend/config/config.yml` — 后端配置文件（端口、数据库、JWT、日志等）
- `frontend/vite.config.js` — Vite 配置（代理、自动导入、别名）
- `backend/init.sql` — 初始建表 SQL（首次启动时执行）
