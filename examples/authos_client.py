import requests

class AuthosClient:
    """
    负责与权限系统 (Authos) 进行通信的客户端
    """
    def __init__(self, host, app_code, app_secret):
        self.host = host.rstrip('/')
        self.app_code = app_code
        self.app_secret = app_secret

    def proxy_login(self, username, password):
        """
        后端透传登录：将用户提交的账号密码，加上App身份凭证，发给权限系统
        """
        url = f"{self.host}/api/public/proxy-login"
        payload = {
            "appCode": self.app_code,
            "appSecret": self.app_secret,
            "username": username,
            "password": password
        }
        
        try:
            # 发起请求
            response = requests.post(url, json=payload, timeout=5)
            # 返回 JSON 数据和 HTTP 状态码
            return response.json(), response.status_code
        except requests.RequestException as e:
            # 处理网络错误等异常
            return {"message": f"Authos service error: {str(e)}"}, 500

    def check_access(self, token, path, method):
        """
        统一鉴权：将Token和请求信息，加上App身份凭证，发给权限系统校验
        """
        url = f"{self.host}/api/public/check-access"
        payload = {
            "appCode": self.app_code,
            "appSecret": self.app_secret,
            "token": token,
            "obj": path,
            "act": method
        }

        try:
            response = requests.post(url, json=payload, timeout=5)
            
            if response.status_code == 200:
                data = response.json()
                # 假设 Authos 返回结构中有 "allowed": true/false
                return data.get("allowed", False), data
            
            # 如果状态码不是 200 (如 401 Token 无效)，则视为无权限
            return False, response.json()
            
        except requests.RequestException:
            return False, {"message": "Authos service unavailable"}
