import json
from django.http import JsonResponse
from django.conf import settings

# 假设 authos_client.py 在同一目录下，或者在 Python Path 中
try:
    from .authos_client import AuthosClient
except ImportError:
    from authos_client import AuthosClient

# ==========================================
# 配置初始化
# ==========================================
# 注意：这是一个示例文件，需要在你的 Django 项目中正确配置 settings
# 假设这是你的 Django 项目 settings.py 中的配置
# AUTHOS_HOST = "http://localhost:8080"
# AUTHOS_APP_CODE = "your_app_code"
# AUTHOS_APP_SECRET = "your_app_secret"

# 初始化客户端 (建议在 Django 的 apps.py 或 settings.py 中初始化单例)
authos_client = AuthosClient(
    host="http://localhost:8080", 
    app_code="example_app",     
    app_secret="example_secret"  
)

# ==========================================
# 1. 登录接口示例 (Django View)
# ==========================================

def login_view(request):
    """
    用户登录视图 (Controller)
    前端 POST 用户名和密码到此接口
    """
    if request.method != 'POST':
        return JsonResponse({"message": "Method not allowed"}, status=405)

    try:
        data = json.loads(request.body)
        username = data.get('username')
        password = data.get('password')
    except json.JSONDecodeError:
        return JsonResponse({"message": "Invalid JSON"}, status=400)

    if not username or not password:
        return JsonResponse({"message": "Username and password required"}, status=400)

    # 关键步骤：调用 Authos Client 进行透传登录
    # 租户后端此时扮演“代理”角色
    resp_data, status_code = authos_client.proxy_login(username, password)

    if status_code == 200:
        # 登录成功
        # 1. 拿到 Token (resp_data['token'])
        # 2. 可以选择把 Token 直接返回给前端，或者存入 HttpOnly Cookie
        return JsonResponse(resp_data)
    else:
        # 登录失败，透传错误信息 (如用户名密码错误，或 AppSecret 错误)
        return JsonResponse(resp_data, status=status_code)


# ==========================================
# 2. 鉴权中间件示例 (Django Middleware)
# ==========================================

class AuthosMiddleware:
    """
    权限拦截中间件
    拦截进入后端的所有请求，向权限系统确认是否有权访问
    """
    def __init__(self, get_response):
        self.get_response = get_response
        # 定义不需要鉴权的白名单路径 (如登录接口、健康检查)
        self.white_list = [
            '/api/login',
            '/health',
            '/favicon.ico',
        ]

    def __call__(self, request):
        path = request.path
        
        # 1. 白名单放行
        if path in self.white_list:
            return self.get_response(request)

        # 2. 获取 Token
        # 通常约定前端将 Token 放在 Header: Authorization: Bearer <token>
        auth_header = request.headers.get('Authorization', '')
        if not auth_header.startswith('Bearer '):
            return JsonResponse({"message": "Missing or invalid Authorization header"}, status=401)
        
        token = auth_header.split(' ')[1]

        # 3. 调用 Authos 进行鉴权
        # 将 "当前请求的路径" 和 "请求方法" 发给 Authos
        method = request.method
        
        # 注意：这里可能会有性能损耗，生产环境可以考虑加一层本地缓存 (Redis/Memory)
        # 缓存 key 可以是 f"{token}:{path}:{method}"
        allowed, details = authos_client.check_access(token, path, method)

        if not allowed:
            # 鉴权失败 (无权限 或 Token无效)
            return JsonResponse({
                "message": "Forbidden", 
                "details": details.get("message", "Access denied")
            }, status=403)

        # 4. (可选) 将用户信息注入 request
        # 这样在后续的 View 中，可以通过 request.user_id 获取当前用户ID
        if "userId" in details:
            request.user_id = details["userId"]

        # 5. 鉴权通过，放行给下一个中间件或 View
        response = self.get_response(request)
        return response
