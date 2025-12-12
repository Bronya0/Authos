package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"Authos/pkg/utils"
)

// JWTMiddleware JWT 认证中间件
type JWTMiddleware struct {
	JWTConfig *utils.JWTConfig
}

// NewJWTMiddleware 创建 JWT 中间件实例
func NewJWTMiddleware(jwtConfig *utils.JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{
		JWTConfig: jwtConfig,
	}
}

// Middleware 返回 JWT 中间件函数
func (j *JWTMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从请求头获取 Authorization
		authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header is required"})
			}

			// 解析 Bearer 令牌
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid authorization header format"})
			}

			// 验证令牌
			claims, err := j.JWTConfig.ParseToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
			}

			// 将用户信息存储到上下文
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)

			return next(c)
		}
	}
}
