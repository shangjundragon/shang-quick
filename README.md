# Shang Quick Admin

一个前后端分离的轻量级后台管理系统脚手架，基于 **Go + Gin + SQLite** 后端和 **Vue 3 + Vite + Naive UI** 前端构建。开箱即用，适合快速搭建中小型管理后台。

## 目录

- [特性](#特性)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [功能模块](#功能模块)
- [配置说明](#配置说明)
- [开发约定](#开发约定)
- [API 接口规范](#api-接口规范)
- [部署指南](#部署指南)
- [常见问题](#常见问题)

---

## 特性

- **RBAC 权限管理** — 基于角色的访问控制，支持用户、角色、菜单三级权限体系
- **动态菜单** — 后端配置菜单，前端动态渲染侧边栏，支持目录、菜单、按钮三级类型
- **操作日志** — 自动记录关键操作的请求、响应、IP、耗时等信息
- **数据字典** — 部门树形管理、菜单排序、角色分配
- **逻辑删除** — 所有核心数据表支持逻辑删除，保留审计轨迹
- **JWT 认证** — 基于 Token 的无状态认证，支持配置过期时间
- **前端权限指令** — `v-permission="['user:add']"` 按钮级权限控制
- **限流保护** — 登录接口内置 IP+用户名级别的限流防护
- **密码安全** — bcrypt 哈希 + 强密码策略（大小写字母+数字+特殊字符）
- **事务安全** — 关键操作使用数据库事务保证数据一致性

---

## 技术栈

### 后端

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.25.5 | 编程语言 |
| Gin | latest | Web 框架 |
| SQLite | 3.x | 嵌入式数据库（生产可切换 MySQL） |
| dbw | 1.2.7 | 自定义 ORM（链式查询、自动时间戳、逻辑删除） |
| JWT | v5 | Token 认证 |
| bcrypt | - | 密码哈希 |
| Zap | - | 结构化日志 |
| Viper | - | 配置管理 |

### 前端

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.5+ | 渐进式框架 |
| Vite | 8.x | 构建工具 |
| Naive UI | - | UI 组件库 |
| Pinia | 3.x | 状态管理 |
| Axios | - | HTTP 请求 |
| Vue Router | 4.x | 前端路由 |

---

## 快速开始

### 环境要求

- **Go** >= 1.23
- **Node.js** >= 18
- **pnpm** >= 8

### 1. 克隆项目

```bash
git clone <repository-url>
cd shang-quick
```

### 2. 启动后端

```bash
cd backend

# 安装依赖
go mod tidy

# 确保必要目录存在
mkdir -p storage/logs

# 启动服务
go run main.go
```

后端默认监听 `:8084`，首次启动会自动创建 `data.db` SQLite 数据库并初始化表结构和默认数据。

### 3. 启动前端

```bash
cd frontend

# 安装依赖
pnpm install

# 开发模式（带热重载）
pnpm dev
```

前端默认运行在 `http://localhost:5173`，通过 Vite 代理将 `/api` 请求转发到后端 `localhost:8084`。

### 4. 访问系统

打开浏览器访问 `http://localhost:5173`

默认管理员账号：
- 用户名：`admin`
- 密码：`admin123`（首次登录后建议立即修改）

---

## 项目结构

```
shang-quick/
├── backend/                          # 后端 Go 项目
│   ├── config/
│   │   └── config.yml                # 配置文件
│   ├── controller/                   # 控制器层（HTTP Handler）
│   │   ├── auth_controller.go        # 认证相关
│   │   ├── user_controller.go        # 用户管理
│   │   ├── role_controller.go        # 角色管理
│   │   ├── menu_controller.go        # 菜单管理
│   │   ├── dept_controller.go        # 部门管理
│   │   ├── profile_controller.go     # 个人中心
│   │   └── oper_log_controller.go    # 操作日志
│   ├── service/                      # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── role_service.go
│   │   ├── menu_service.go
│   │   ├── dept_service.go
│   │   └── oper_log_service.go
│   ├── model/                        # 数据模型层
│   │   ├── user.go                   # SysUser
│   │   ├── role.go                   # SysRole
│   │   ├── menu.go                   # SysMenu
│   │   ├── dept.go                   # SysDept
│   │   ├── oper_log.go               # SysOperLog
│   │   ├── user_role.go              # SysUserRole（关联表）
│   │   └── role_menu.go              # SysRoleMenu（关联表）
│   ├── middleware/                   # Gin 中间件
│   │   ├── jwt_auth.go               # JWT 认证
│   │   ├── permission.go             # RBAC 权限校验
│   │   ├── oper_log.go               # 操作日志记录
│   │   ├── trace_logger.go           # TraceID 链路追踪
│   │   └── cors.go                   # 跨域处理
│   ├── router/
│   │   └── router.go                 # 路由注册中心
│   ├── pkg/                          # 公共包
│   │   ├── jwt/                      # Token 生成与解析
│   │   ├── password/                 # bcrypt 密码哈希
│   │   ├── cache/                    # 缓存抽象层（内存/Redis）
│   │   ├── ratelimit/                # 限流器
│   │   ├── logger/                   # Zap 日志初始化
│   │   ├── global_vars/              # 全局变量
│   │   ├── res_util/                 # 统一响应封装
│   │   └── req_util/                 # 请求绑定工具
│   ├── bootstrap/
│   │   ├── bootstrap.go              # 启动流程
│   │   └── init_db.go                # 数据库初始化
│   ├── main.go                       # 入口文件
│   ├── go.mod                        # Go 模块
│   ├── init.sql                      # 初始建表 SQL
│   └── data.db                       # SQLite 数据文件（运行时生成）
│
├── frontend/                         # 前端 Vue 项目
│   ├── src/
│   │   ├── api/                      # API 请求层
│   │   ├── store/
│   │   │   └── user.js               # Pinia 用户状态
│   │   ├── router/
│   │   │   └── index.js              # 路由配置
│   │   ├── layout/                   # 布局组件
│   │   │   ├── index.vue             # 主布局
│   │   │   ├── components/
│   │   │   │   ├── Sidebar.vue       # 侧边栏菜单
│   │   │   │   ├── Navbar.vue        # 顶部导航
│   │   │   │   └── AppMain.vue       # 内容区
│   │   ├── views/                    # 页面视图
│   │   │   ├── login/                # 登录页
│   │   │   ├── dashboard/            # 仪表盘
│   │   │   ├── profile/              # 个人中心
│   │   │   └── system/               # 系统管理
│   │   │       ├── user/             # 用户管理
│   │   │       ├── role/             # 角色管理
│   │   │       ├── menu/             # 菜单管理
│   │   │       ├── dept/             # 部门管理
│   │   │       └── operLog/          # 操作日志
│   │   ├── directives/
│   │   │   └── permission.js         # v-permission 指令
│   │   ├── utils/
│   │   │   ├── request.jsx           # Axios 封装
│   │   │   └── format.js             # 格式化工具
│   │   ├── main.js                   # 入口
│   │   └── App.vue                   # 根组件
│   ├── vite.config.js                # Vite 配置
│   └── package.json
│
└── README.md
```

---

## 功能模块

### 用户管理（/system/user）

- 用户 CRUD（增删改查）
- 分配角色（支持多选）
- 状态启用/禁用（开关控制）
- 重置密码（重置为默认密码，建议首次登录后修改）
- 按用户名、手机号、状态、部门筛选
- 显示部门名称和角色名称

**权限标识**：`user:list`, `user:add`, `user:edit`, `user:delete`, `user:resetPwd`

### 角色管理（/system/role）

- 角色 CRUD
- 分配菜单权限（树形勾选，级联选择）
- 删除前检查是否关联用户
- 角色编码唯一（如 `admin`, `editor`）

**权限标识**：`role:list`, `role:add`, `role:edit`, `role:delete`, `role:assign`

### 菜单管理（/system/menu）

- 菜单 CRUD（树形结构）
- 三级菜单类型：目录(0)、菜单(1)、按钮(2)
- 按钮权限标识（如 `user:add`）用于 `v-permission` 指令
- 显示排序、图标、路由路径、组件路径
- 删除前检查是否被角色引用

**权限标识**：`menu:list`, `menu:add`, `menu:edit`, `menu:delete`

### 部门管理（/system/dept）

- 部门 CRUD（树形结构）
- 负责人、联系电话、邮箱
- 删除前检查是否有子部门或关联用户

**权限标识**：`dept:list`, `dept:add`, `dept:edit`, `dept:delete`

### 操作日志（/system/operLog）

- 自动记录 POST 请求的写操作
- 记录字段：标题、操作类型、请求URL、请求数据、响应数据、操作人、IP、耗时、状态、错误信息
- 支持按标题、操作人筛选

### 个人中心（/profile）

- 查看/修改个人信息
- 修改密码（需输入旧密码验证）

---

## 配置说明

配置文件位于 `backend/config/config.yml`：

```yaml
AppDebug: true                    # 调试模式（true=日志输出控制台，false=写入文件）
Port: ":8084"                     # 服务监听端口

# 日志配置
Logs:
  GinLogName: "storage/logs/gin.log"              # Gin 访问日志
  AppFileLogName: "storage/logs/app.log"          # 应用日志
  TextFormat: true                                 # 文本格式（false=json）
  TimePrecision: "millisecond"                     # 时间精度
  MaxSize: 100                                     # 单个日志文件最大 MB
  MaxBackups: 5                                    # 保留备份数
  MaxAge: 30                                       # 保留天数
  Compress: true                                   # 是否压缩

# 数据库配置
Db:
  DriverName: "sqlite"             # 数据库驱动（sqlite/mysql/postgres）
  Dsn: "data.db"                   # 连接字符串

# JWT 配置
Jwt:
  Secret: "your-secret-key-here"   # JWT 签名密钥（生产环境必须修改！）
  ExpireHours: 24                  # Token 过期时间（小时）

# 缓存配置
Cache:
  UseRedis: false                  # 是否使用 Redis（false=内存缓存）
```

### 生产环境注意

1. **修改 JWT Secret**：`Jwt.Secret` 必须修改为随机长字符串
2. **关闭调试模式**：`AppDebug: false`
3. **切换数据库**：SQLite 适合开发，生产建议使用 MySQL/Postgres
4. **配置日志路径**：确保 `storage/logs` 目录存在且有写入权限

---

## 开发约定

### 后端

#### 1. 分层架构

```
HTTP Request → Controller → Service → Model → Database
                      ↓
                Middleware（认证、权限、日志、限流）
```

- **Controller**：处理 HTTP 请求，参数校验，调用 Service，返回响应
- **Service**：业务逻辑，数据库事务，调用 dbw ORM
- **Model**：数据模型，仅包含结构体定义和表名方法
- **Middleware**：横切关注点（认证、权限、日志、跨域）

#### 2. dbw ORM 使用规范

```go
// 基本查询
users, err := dbw.New[model.SysUser](
    dbw.WithConfig(global_vars.DbConfig),
    dbw.WithContext(ctx),
).SelectList()

// 条件查询
users, err := dbw.New[model.SysUser](...).
    Eq("status", 1).
    Like("username", "%admin%").
    OrderByDesc("create_time").
    SelectPage(pageNum, pageSize)

// 事务操作（使用 ExecuteTx）
err := dbw.ExecuteTx(func(tx *sql.Tx) error {
    _, err := dbw.New[model.SysUser](
        dbw.WithConfig(global_vars.DbConfig),
        dbw.WithTx(tx),
        dbw.WithContext(ctx),
    ).Insert(user)
    if err != nil {
        return err
    }
    // ... 其他操作
    return nil
}, global_vars.DbConfig.Db)
```

#### 3. 模型字段标签

```go
type SysUser struct {
    Id         int64   `dbw:"primaryKey" json:"id,string"`                 // 主键（前端用字符串避免精度丢失）
    Username   string  `dbw:"column:username" json:"username"`             // 列映射
    Status     int     `dbw:"column:status;default:1" json:"status"`       // 默认值
    DelFlag    string  `dbw:"column:del_flag;tableLogic" json:"delFlag"`   // 逻辑删除字段
    CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli"`      // 自动创建时间戳
    UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli"`      // 自动更新时间戳
    CreateBy   int64   `dbw:"column:create_by" json:"createBy,string"`     // 创建人（通过 Hook 自动填充）
    UpdateBy   int64   `dbw:"column:update_by" json:"updateBy,string"`     // 更新人（通过 Hook 自动填充）
}
```

#### 4. 响应规范

所有接口统一返回格式：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": { ... },
  "trace_id": "abc123"
}
```

状态码：
- `200` — 成功
- `401` — 未登录/Token 过期
- `403` — 无权限
- `500` — 服务器错误

#### 5. 权限标识约定

- `模块:操作`，如 `user:add`, `user:edit`, `user:delete`, `user:list`
- 按钮类型菜单（`menu_type=2`）的 `perm` 字段即为权限标识
- 超级管理员（`role_code=admin`）拥有所有权限

### 前端

#### 1. 按钮级权限

```vue
<!-- 有权限才显示按钮 -->
<n-button v-permission="['user:add']">新增用户</n-button>

<!-- 多权限（满足任一即可） -->
<n-button v-permission="['user:edit', 'user:add']">编辑</n-button>
```

#### 2. API 调用

```javascript
import { getUserList, addUser, updateUser, deleteUser } from '@/api/user'

// 查询
const res = await getUserList({ pageNum: 1, pageSize: 10, username: 'admin' })
// res = { list: [...], total: 100 }

// 新增/编辑
await addUser({ username: 'test', password: 'Test1234!', ... })

// 删除
await deleteUser({ id: '123' })
```

#### 3. 路由权限

路由在编译时静态定义，但侧边栏菜单根据后端返回的 `menus` 动态渲染。用户无权限的菜单不会显示在侧边栏中。

---

## API 接口规范

### 认证相关

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 登录 | 公开 |
| GET | `/api/v1/auth/info` | 获取当前用户信息 | 需登录 |

### 用户管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/user/list` | 用户列表 | user:list |
| POST | `/api/v1/user/add` | 新增用户 | user:add |
| POST | `/api/v1/user/edit` | 编辑用户 | user:edit |
| POST | `/api/v1/user/changeStatus` | 修改状态 | user:edit |
| POST | `/api/v1/user/resetPwd` | 重置密码 | user:resetPwd |
| POST | `/api/v1/user/delete` | 删除用户 | user:delete |
| GET | `/api/v1/user/roleIds` | 获取用户角色ID | user:list |

### 角色管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/role/list` | 角色列表 | role:list |
| POST | `/api/v1/role/add` | 新增角色 | role:add |
| POST | `/api/v1/role/edit` | 编辑角色 | role:edit |
| POST | `/api/v1/role/delete` | 删除角色 | role:delete |
| GET | `/api/v1/role/menuIds` | 获取角色菜单ID | role:list |
| POST | `/api/v1/role/assignMenu` | 分配菜单 | role:assign |

### 菜单管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/menu/list` | 菜单列表 | menu:list |
| POST | `/api/v1/menu/add` | 新增菜单 | menu:add |
| POST | `/api/v1/menu/edit` | 编辑菜单 | menu:edit |
| POST | `/api/v1/menu/delete` | 删除菜单 | menu:delete |

### 部门管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/dept/list` | 部门列表 | dept:list |
| POST | `/api/v1/dept/add` | 新增部门 | dept:add |
| POST | `/api/v1/dept/edit` | 编辑部门 | dept:edit |
| POST | `/api/v1/dept/delete` | 删除部门 | dept:delete |

### 个人中心

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/profile` | 获取个人信息 | 需登录 |
| POST | `/api/v1/profile/update` | 更新个人信息 | 需登录 |
| POST | `/api/v1/profile/updatePwd` | 修改密码 | 需登录 |

### 操作日志

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/operLog/list` | 操作日志列表 | 需登录 |

---

## 部署指南

### 1. 构建前端

```bash
cd frontend
pnpm install
pnpm build
```

构建产物位于 `frontend/dist/`。

### 2. 构建后端

```bash
cd backend
# Linux/Mac
go build -o shang-quick .

# Windows
go build -o shang-quick.exe .
```

### 3. 生产环境配置

```yaml
# config.yml
AppDebug: false
Port: ":8084"

Jwt:
  Secret: "change-this-to-a-random-long-string-in-production"
  ExpireHours: 24

Db:
  DriverName: "mysql"
  Dsn: "user:password@tcp(127.0.0.1:3306)/shang_quick?charset=utf8mb4&parseTime=True"
```

### 4. 目录准备

```bash
mkdir -p storage/logs
# 确保可写权限
chmod -R 755 storage/
```

### 5. 运行

```bash
./shang-quick
```

### 6. Nginx 反向代理（推荐）

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态资源
    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8084/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### 7. 使用 systemd 管理（Linux）

创建 `/etc/systemd/system/shang-quick.service`：

```ini
[Unit]
Description=Shang Quick Admin
After=network.target

[Service]
Type=simple
User=www
WorkingDirectory=/path/to/backend
ExecStart=/path/to/backend/shang-quick
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable shang-quick
systemctl start shang-quick
```

---

## 常见问题

### Q: 前端启动后提示 "登录已过期" 怎么办？

A: 确保后端服务已正常启动，且前端 `vite.config.js` 中的代理配置正确指向后端地址。

### Q: 如何修改默认管理员密码？

A: 登录系统后，进入「个人中心」→「修改密码」，或管理员在用户管理中重置密码。

### Q: 如何新增一个菜单页面？

A:
1. 在「菜单管理」中新增菜单（类型选择「菜单」）
2. 填写菜单名称、路由路径、组件路径
3. 在「角色管理」中为角色分配该菜单权限
4. 前端在 `src/views/` 下创建对应页面组件
5. 在 `src/router/index.js` 中注册路由（或配置动态路由）

### Q: 数据库如何迁移到 MySQL？

A:
1. 修改 `config.yml` 中 `Db.DriverName` 为 `mysql`
2. 修改 `Db.Dsn` 为 MySQL 连接字符串
3. 运行后端，dbw 会自动建表（需要确保 MySQL 已安装对应驱动：`go get -u github.com/go-sql-driver/mysql`）
4. 手动迁移数据或使用数据库迁移工具

### Q: 密码策略太严格，能否放宽？

A: 可以修改 `backend/pkg/password/password.go` 中的 `ValidatePasswordStrong` 函数，或在使用处改为 `ValidatePassword`（仅检查长度）。

### Q: 如何关闭操作日志？

A: 在 `router/router.go` 中移除对应路由的 `middleware.OperLog(...)` 中间件即可。

### Q: 权限缓存多久刷新一次？

A: 默认缓存 10 分钟。修改 `user_service.go` 中 `cache.GlobalCache.Set(...)` 的 `10*time.Minute` 参数即可调整。

### Q: 出现 "storage/logs 目录不存在" 错误？

A: 在项目根目录执行 `mkdir -p storage/logs`，确保该目录存在且有写入权限。

---

## 许可证

MIT License
