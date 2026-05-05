# AGENTS.md

## 项目结构

前后端分离的全栈项目：
- **frontend/** — Vue 3 + Vite + Naive UI
- **backend/** — Go + Gin + SQLite

## 前端 (frontend/)

- **包管理器**: pnpm
- **框架**: Vue 3 (`<script setup>` 语法)
- **构建工具**: Vite
- **UI 库**: Naive UI
- **自动导入**: `unplugin-auto-import` (Vue、naive-ui API) + `unplugin-vue-components` (Naive UI 组件)
- **HTTP 请求**: Axios，基地址 `/api`，代理到后端 `localhost:8084`
- **启动命令**: `cd frontend && npm run dev`
- **构建命令**: `cd frontend && npm run build`
- **别名**: `@` 指向 `./src`

## 后端 (backend/)

- **语言**: Go 1.25.5
- **框架**: Gin
- **数据库**: SQLite (使用 `github.com/glebarez/go-sqlite`)
- **ORM**: 自定义 `github.com/shangjundragon/dbw`
- **配置**: `config/config.yml` (YAML 格式)
- **端口**: `:8084`
- **启动命令**: `cd backend && go run main.go`
- **必要目录**: `config/`、`storage/logs/` 必须存在，否则启动失败
- **数据库文件**: `data.db` (SQLite 数据文件)

## 开发约定

- **回复语言**: 简体中文
- **数据库**: 不使用外键约束
- **日志**: 开发模式 (`AppDebug: true`) 输出到控制台，生产模式写入文件
- **代码风格**: 遵循项目现有风格，不添加无关注释
- **Git**: 不自动提交代码
- **启动**: 不主动启动项目，除非用户明确要求
