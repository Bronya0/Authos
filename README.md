# Authos 使用全流程说明（面向接入方系统）

本文只关注“业务系统如何使用 Authos 做统一认证和权限校验”的完整流程，分为：

1. 启动 Authos 服务  
2. 系统管理员初始化应用  
3. 业务系统获取应用凭证并登录  
4. 在 Authos 中配置角色、接口权限和逻辑权限 Key  
5. 在 Authos 中创建业务用户并绑定角色  
6. 业务用户登录获取 userToken  
7. 业务系统在每次请求前调用 Authos 校验权限  

---

## 1. 启动 Authos 服务

### 1.1 编译

在项目根目录执行：

```bash
go build -o authos ./cmd/server
```

Windows 下会生成 `authos.exe`。

### 1.2 启动服务

```bash
./authos
```

或在 Windows：

```cmd
authos.exe
```

默认监听地址：

- 后端 API: `http://localhost:8080`

首次启动时，如果数据库中没有数据，会自动：

- 创建一个“默认应用”
- 创建一个“超级管理员角色”
- 引导创建一个管理员账号（用户名可自定义，密码你自己输入或自动生成）

> 注意：种子数据是通过 [internal/service/db.go](file:///f:/coding/go/src/Authos/internal/service/db.go) 自动初始化的。

---

## 2. 系统管理员初始化应用

### 2.1 系统管理员登录

接口（无需前置 Token）：

- `POST /api/public/system-login`

请求体示例：

```json
{
  "username": "admin",
  "password": "你的管理员密码"
}
```

返回示例（关键字段）：

```json
{
  "token": "系统管理员JWT",
  "user": { "id": 1, "username": "admin" },
  "app": { "id": 1, "uuid": "app-uuid", "code": "default" }
}
```

使用方式：

- 把返回的 `token` 存为 `systemToken`
- 后续管理端调用带上请求头：

```http
X-System-Token: Bearer <systemToken>
```

前端内置逻辑可以参考 [web-vue3/src/api/index.js](file:///f:/coding/go/src/Authos/web-vue3/src/api/index.js#L11-L48)。

### 2.2 创建业务应用（业务系统 B）

系统管理员使用 `systemToken` 在管理界面或通过接口创建应用：

- `POST /api/v1/applications`

关键字段：

- `name`: 应用名称，例如 `"业务系统B"`
- `code`: 应用代码，用于用户登录时区分租户，例如 `"biz-b"`

创建成功后，你会拿到：

- `id`
- `uuid`
- `code`
- `secretKey`（应用密钥）

业务系统后续登录将使用 `uuid + secretKey`。

---

## 3. 业务系统获取应用凭证并登录

### 3.1 应用登录获取 appToken

接口：

- `POST /api/public/app-login`

请求体示例：

```json
{
  "appUuid": "应用的UUID",
  "appSecret": "应用的secretKey"
}
```

返回示例：

```json
{
  "token": "appToken",
  "app": {
    "id": 2,
    "uuid": "应用UUID",
    "code": "biz-b"
  }
}
```

使用方式：

- 业务系统在后台持有 `appToken`，每次调用 Authos 业务 API 时带上：

```http
X-App-Token: Bearer <appToken>
X-App-ID: <app.id>
```

前端示例可参考：

- 拦截器中设置 `X-App-Token` 和 `X-App-ID`  
  [web-vue3/src/api/index.js](file:///f:/coding/go/src/Authos/web-vue3/src/api/index.js#L11-L44)

---

## 4. 在 Authos 中配置角色、接口权限和逻辑权限 Key

这一步是“把业务权限模型落到 Authos”的关键。

### 4.1 定义接口权限（带逻辑权限 Key）

接口（需 `userToken` 或 `systemToken + appToken`，下文假设在管理后台操作）：

- `POST /api/v1/api-permissions`

请求体示例：

```json
{
  "key": "user:create",
  "name": "创建用户",
  "path": "/api/v1/users",
  "method": "*",
  "description": "业务系统B创建用户的接口权限"
}
```

字段含义：

- `key`：**逻辑权限标识**，业务系统后续校验就用它（与具体 HTTP 路径解耦）
- `path` + `method`：用于在 Authos 控制台展示或未来做网关级校验

前端配置入口参考：

- 权限列表界面 [web-vue3/src/views/Permissions.vue](file:///f:/coding/go/src/Authos/web-vue3/src/views/Permissions.vue)

### 4.2 创建角色

接口：

- `POST /api/v1/roles`

请求体示例：

```json
{
  "name": "用户管理员"
}
```

角色属于某个应用（通过 Token 的 appID 关联），典型角色如：

- `UserAdmin`：管理用户
- `Auditor`：只读审计日志

### 4.3 给角色分配接口权限

有两种方式：

1. 通过“角色接口权限分配弹窗”：
   - 角色列表 → 点击“权限”按钮 → 出现接口权限列表 → 勾选权限 → 保存
   - 实际调用的是：
     - `GET /api/v1/api-permissions/roles/:roleUUID`
     - `POST /api/v1/api-permissions/roles/:roleUUID`
     - `DELETE /api/v1/api-permissions/roles/:roleUUID`
   - 对应前端组件：  
     [web-vue3/src/components/RoleApiPermissionModal.vue](file:///f:/coding/go/src/Authos/web-vue3/src/components/RoleApiPermissionModal.vue)

2. 直接调用后端接口（适合自动化脚本）：

   - 添加权限：

     ```http
     POST /api/v1/api-permissions/roles/{roleUUID}
     Content-Type: application/json

     {
       "permissionUUID": "接口权限的UUID"
     }
     ```

   - 移除权限：

     ```http
     DELETE /api/v1/api-permissions/roles/{roleUUID}
     Content-Type: application/json

     {
       "permissionUUID": "接口权限的UUID"
     }
     ```

绑定成功后，内部会通过 Casbin 生成策略：

```text
p, role:<roleUUID>, <permissionKey>, *
```

例如：

```text
p, role:xxxx-uuid, user:create, *
```

---

## 5. 在 Authos 中创建业务用户并绑定角色

业务用户由 Authos 管理，外部系统只需要保存 `userId` 即可。

### 5.1 创建用户

接口：

- `POST /api/v1/users`

请求体示例：

```json
{
  "username": "alice",
  "password": "P@ssw0rd",
  "status": 1,
  "roleIds": [1, 2]
}
```

说明：

- `roleIds` 是角色的数据库 ID（非 UUID），同一应用内校验。
- 服务端会：
  - 哈希密码
  - 在事务中通过 `Association("Roles").Replace(roles)` 绑定用户和角色  
    参见 [internal/service/user.go](file:///f:/coding/go/src/Authos/internal/service/user.go#L24-L64)。

前端入口：

- 用户管理界面 [web-vue3/src/views/Users.vue](file:///f:/coding/go/src/Authos/web-vue3/src/views/Users.vue)

### 5.2 更新用户角色

接口：

- `PUT /api/v1/users/{id}`

请求体与创建类似，带上新的 `roleIds` 即可重新绑定。

---

## 6. 业务用户登录获取 userToken

业务用户的登录入口：

- `POST /api/public/login`

请求体示例：

```json
{
  "username": "alice",
  "password": "P@ssw0rd",
  "appCode": "biz-b"
}
```

说明：

- `appCode`：应用代码，用于区分租户（在创建应用时设置）

返回示例：

```json
{
  "token": "userToken",
  "user": { "id": 10, "username": "alice" },
  "app": { "id": 2, "uuid": "应用UUID", "code": "biz-b" }
}
```

业务系统前端或后端需要保存：

- `userToken`（作为 `Authorization: Bearer` 头）
- `user.id`（后续调用权限检查时要带的 `userId`）

前端参考：

- userToken 的注入见 [web-vue3/src/api/index.js](file:///f:/coding/go/src/Authos/web-vue3/src/api/index.js#L27-L30)

---

## 7. 业务系统每次请求前调用 Authos 校验权限

接入方式的核心就是这一条：**业务系统在执行关键业务前，调用 Authos 校验某个用户是否拥有某个逻辑权限 Key**。

### 7.1 检查逻辑权限 Key

接口：

- `POST /api/v1/auth/check-permission`

请求体示例：

```json
{
  "userId": 10,
  "permission": "user:create"
}
```

返回示例：

```json
{
  "allowed": true,
  "message": "Permission checked successfully"
}
```

解释：

- Authos 内部会：
  - 通过 `userId` 查出用户并预加载其角色
  - 使用每个角色的 `UUID` 构造 `sub = "role:<roleUUID>"`
  - 使用 `permission` 作为 `obj`，`"*"` 作为 `act`
  - 调用 Casbin 逐个校验，任一角色命中即返回 `allowed = true`

### 7.2 业务系统伪代码（Go 示例）

示例：在业务系统 B 中封装一个简单客户端：

```go
type AuthosClient struct {
	BaseURL string
	Client  *http.Client
}

type CheckPermissionResp struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message"`
}

func (c *AuthosClient) CheckPermission(ctx context.Context, userID uint, permission string) (bool, error) {
	body := map[string]interface{}{
		"userId":     userID,
		"permission": permission,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/v1/auth/check-permission", bytes.NewReader(b))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("authos check failed, status=%d", resp.StatusCode)
	}

	var r CheckPermissionResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}
	return r.Allowed, nil
}
```

业务代码中使用：

```go
allowed, err := authosClient.CheckPermission(ctx, userID, "user:create")
if err != nil {
    // 处理 Authos 异常，比如降级为拒绝
}
if !allowed {
    // 拒绝这次业务操作
}
// 继续执行真正的业务逻辑
```

你可以用同样的方式在任何语言（Node、Java、Python 等）中封装一个轻量客户端，核心就是：

1. 记住 Authos 的访问地址  
2. 每次带上 `userId` 和约定好的 `permissionKey` 调用 `/api/v1/auth/check-permission`  

---

## 8. 整体时序总览

从“运维/平台方”到“业务系统 B”到“最终业务用户”的完整时序如下：

1. 运维启动 Authos 服务
2. 系统管理员 `system-login` 登录管理控制台
3. 管理控制台中创建“业务系统 B”应用，拿到 `appUuid + secretKey + code`
4. 业务系统 B 在配置中心写入这组凭证
5. 业务系统 B 启动时或按需调用 `app-login`，拿到 `appToken`，请求时设置 `X-App-Token` 和 `X-App-ID`
6. 管理员在控制台中：
   - 创建角色
   - 创建接口权限，并为每个接口定义统一的 `key`
   - 将接口权限分配给角色
7. 管理员在控制台中为应用创建业务用户，并绑定角色
8. 业务用户从业务系统的登录入口发起登录，业务系统将登录请求转发给 Authos 的 `/public/login`（带 `appCode`），获得 `userToken`
9. 业务系统前端保存 userToken，调用自身后端时带上 userToken
10. 业务系统后端在执行关键操作前调用 Authos 的 `/api/v1/auth/check-permission`，基于 `userId + permissionKey` 判断是否放行

做到这 10 步，业务系统就完全接入了 Authos 的统一认证与权限服务。
