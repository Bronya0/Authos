推荐的标准流程：
1、登录

概述：用户输入 -> App后端透传(AppID+Secret+账号密码) -> 权限系统校验 -> 返回Token -> 前端存储
详细：

	用户输入：用户在 租户 App 的前端 输入 用户名 和 密码。
	后端透传：
		租户 App 前端将账号密码发给 租户 App 后端。
		租户 App 后端 此时扮演“客户端”角色，带上自己的 App ID + App Secret（作为身份凭证）以及用户提交的 用户名 + 密码，调用 权限系统 的 /oauth/token（或类似登录接口）。
	权限系统校验：
		校验 App ID 和 Secret 是否合法。
		在指定 App ID 的隔离域下，查找该用户名并校验密码。
		校验通过后，权限系统生成一个 Token。
	返回 Token：权限系统将 Token 返回给租户 App 后端，租户 App 后端再透传给前端。
	租户前端挂到会话上，请求接口时都带着。

2、用户访问

概述：前端带Token -> App中间件 -> 调用权限系统(AppID+Secret+Token+Path+method) -> 权限系统返回结果 -> 中间件放行/拒绝

详情：
前端带上token给app后端中间件
app后端中间件将app id、app secret、token、访问path、method发给权限系统的校验接口
权限系统返回用户是否有权限、超时、或者用户不存在等
app后端中间件根据返回，选择是否放行或者重定向登录等


# 3、涉及接口
假设提前建好了app galaxy，以及下属test用户：test/test。将生成的appCode、appSecret保存好

1、登录：
POST /api/public/proxy-login
{
  "appCode": "galaxy",
  "appSecret": "xxx",
  "username": "test",
  "password": "test"
}

返回token：
{
  "app": {
  },
  "message": "Proxy login successful",
  "token": "xxx",
  "user": {
  }
}

2、校验用户访问接口：

先在页面上或者通过接口为test用户分配好角色，并绑定权限，例如：/api/v1/users *
接着在中间件里校验用户访问的接口：

POST /api/public/check-access
{
  "appCode": "galaxy",
  "appSecret": "xxx",
  "token": "xxx",
  "obj": "/api/v1/users/1",
  "act": "GET"
}

通过则返回 allowed为true，否则为false：
{
  "allowed": true,
  "message": "Permission checked successfully",
  "userId": 2
}