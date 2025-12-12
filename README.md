# Authos - RBAC权限管理微服务

Authos 是一个通用的、轻量级的 RBAC 权限管理微服务，使用 Go 语言开发，集成了用户管理、角色管理、菜单权限和 API 鉴权功能。

## 技术栈

- **Language**: Go (Golang) 1.21+
- **Web Framework**: Echo (高性能)
- **Database**: SQLite (使用纯 Go 驱动 `github.com/glebarez/sqlite`)
- **ORM**: GORM v2
- **Auth Engine**: Casbin v2 (核心权限控制)
- **Casbin Adapter**: Gorm Adapter (配合 SQLite)
- **Authentication**: JWT (JSON Web Tokens)
- **Deployment**: 使用 Go `embed` 特性打包配置文件和前端静态资源

## 核心功能

- **认证管理**: 登录、登出，JWT 令牌颁发
- **用户管理**: 增删改查用户，关联角色
- **角色管理**: 增删改查角色，分配菜单和 API 权限
- **菜单管理**: 增删改查菜单，支持树形结构
- **权限控制**: 基于 Casbin 的 RBAC 权限控制，支持 RESTful API 路径匹配
- **API 鉴权**: 提供权限检查 API，支持外部服务调用

## 项目结构

```
Authos/
├── cmd/
│   └── server/          # 应用入口
│       └── main.go
├── internal/
│   ├── model/           # 数据库模型
│   │   ├── user.go
│   │   ├── role.go
│   │   └── menu.go
│   ├── service/         # 业务逻辑
│   │   ├── db.go
│   │   ├── casbin.go
│   │   ├── user.go
│   │   ├── role.go
│   │   └── menu.go
│   ├── handler/         # HTTP 处理函数
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── menu.go
│   │   └── authz.go
│   └── middleware/      # 中间件
│       ├── jwt.go
│       └── model.conf   # Casbin 模型配置
├── pkg/
│   └── utils/           # 工具类
│       └── jwt.go
├── web/
│   └── dist/            # 前端构建产物（预留）
├── go.mod
└── README.md
```

## 编译和运行

### 1. 编译

在项目根目录执行以下命令编译项目：

```bash
go build -o authos ./cmd/server
```

这将生成一个名为 `authos` 的可执行文件（在 Windows 上为 `authos.exe`）。

### 2. 运行

直接运行生成的可执行文件：

```bash
./authos
```

在 Windows 上：

```cmd
authos.exe
```

服务将在 `http://localhost:8080` 上启动。

### 访问管理界面

启动服务后，可以通过以下地址访问管理界面：

- **管理后台**: `http://localhost:8080/admin.html`
- **默认管理员账号**: 
  - 用户名: `admin`
  - 密码: `123456`

管理界面包含以下功能：
- 用户管理（增删改查、分配角色）
- 角色管理（增删改查、分配菜单权限）
- 菜单管理（树形结构管理）
- 系统仪表盘（数据统计）

### 3. 跨平台编译

编译 Linux 版本：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o authos-linux ./cmd/server
```

编译 macOS 版本：

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o authos-darwin ./cmd/server
```

## API 文档

### 认证相关

- **登录**: `POST /login`
  - 请求体: `{"username": "admin", "password": "123456"}`
  - 响应: 返回 JWT 令牌和用户信息

- **登出**: `POST /logout`
  - 响应: 登出成功消息

### 权限检查

- **检查权限**: `POST /api/v1/check`
  - 请求体: `{"userId": 1, "obj": "/api/users", "act": "GET"}`
  - 响应: `{"allowed": true, "message": "Permission checked successfully"}`

- **获取用户导航菜单**: `GET /api/v1/user/nav`
  - 响应: 返回用户有权访问的菜单树

### 用户管理

- **创建用户**: `POST /api/v1/users`
- **获取用户列表**: `GET /api/v1/users`
- **获取单个用户**: `GET /api/v1/users/:id`
- **更新用户**: `PUT /api/v1/users/:id`
- **删除用户**: `DELETE /api/v1/users/:id`

### 角色管理

- **创建角色**: `POST /api/v1/roles`
- **获取角色列表**: `GET /api/v1/roles`
- **获取单个角色**: `GET /api/v1/roles/:id`
- **更新角色**: `PUT /api/v1/roles/:id`
- **删除角色**: `DELETE /api/v1/roles/:id`
- **为角色分配菜单**: `POST /api/v1/roles/:id/menus`
- **为角色分配 API 权限**: `POST /api/v1/roles/:id/permissions`

### 菜单管理

- **创建菜单**: `POST /api/v1/menus`
- **获取菜单列表**: `GET /api/v1/menus`
- **获取菜单树**: `GET /api/v1/menus/tree`
- **获取单个菜单**: `GET /api/v1/menus/:id`
- **更新菜单**: `PUT /api/v1/menus/:id`
- **删除菜单**: `DELETE /api/v1/menus/:id`

## 初始数据

系统启动时，如果数据库为空，会自动创建以下初始数据：

- **超级管理员账号**: 
  - 用户名: `admin`
  - 密码: `123456`
  - 角色: 超级管理员

- **测试角色**: 
  - 代码: `test`
  - 名称: 测试角色

- **初始菜单**: 
  - 系统管理
    - 用户管理
    - 角色管理
    - 菜单管理

## 配置说明

当前版本使用硬编码配置，实际部署时建议使用环境变量或配置文件进行配置：

- `jwtSecret`: JWT 密钥，用于生成和验证 JWT 令牌
- `jwtExpireTime`: JWT 令牌过期时间
- `serverAddr`: 服务器监听地址

## 数据库

系统使用 SQLite 数据库，数据库文件 `auth.db` 会自动创建在当前工作目录下。

## 前端集成

系统预留了使用 Go `embed` 特性托管前端构建产物的代码，只需将前端构建产物放置在 `web/dist` 目录下，然后取消注释 `main.go` 中的静态文件服务代码即可。

## 开发说明

1. 安装依赖：
   ```bash
   go mod tidy
   ```

2. 运行开发服务器：
   ```bash
   go run ./cmd/server
   ```

3. 编译生产版本：
   ```bash
   go build -o authos ./cmd/server
   ```

## 许可证

MIT
