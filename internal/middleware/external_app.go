package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"Authos/internal/service"
)

// ExternalAppAuthMiddleware 外部应用认证中间件
type ExternalAppAuthMiddleware struct {
	externalAppService *service.ExternalAppService
}

// NewExternalAppAuthMiddleware 创建外部应用认证中间件实例
func NewExternalAppAuthMiddleware(externalAppService *service.ExternalAppService) *ExternalAppAuthMiddleware {
	return &ExternalAppAuthMiddleware{
		externalAppService: externalAppService,
	}
}

// Middleware 返回Echo中间件函数
func (m *ExternalAppAuthMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从请求头获取令牌
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "缺少授权令牌"})
			}

			// 检查令牌格式
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "无效的令牌格式"})
			}

			token := tokenParts[1]

			// 验证令牌
			app, err := m.externalAppService.ValidateAppToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "无效的令牌"})
			}

			// 将应用信息存储到上下文中
			c.Set("app", app)

			return next(c)
		}
	}
}