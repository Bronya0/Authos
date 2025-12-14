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
			// 尝试从多个请求头获取令牌
			authHeader := c.Request().Header.Get("Authorization")
			systemAuthHeader := c.Request().Header.Get("X-System-Token")
			appAuthHeader := c.Request().Header.Get("X-App-Token")

			var token string
			var tokenType string

			// 确定令牌类型和值
			if systemAuthHeader != "" {
				// 系统管理员令牌
				parts := strings.SplitN(systemAuthHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
					tokenType = "system"
				}
			} else if appAuthHeader != "" {
				// 应用令牌
				parts := strings.SplitN(appAuthHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
					tokenType = "app"
				}
			} else if authHeader != "" {
				// 传统令牌
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
					tokenType = "user"
				}
			}

			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header is required"})
			}

			// 根据令牌类型验证
			switch tokenType {
			case "system":
				claims, err := j.JWTConfig.ParseSystemToken(token)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired system token"})
				}
				// 将系统管理员信息存储到上下文
				c.Set("isSystemAdmin", true)
				c.Set("username", claims.Username)
			case "app":
				claims, err := j.JWTConfig.ParseAppToken(token)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired app token"})
				}
				// 将应用信息存储到上下文
				c.Set("isAppToken", true)
				c.Set("appID", claims.AppID)
				c.Set("appCode", claims.AppCode)
			case "user":
				claims, err := j.JWTConfig.ParseToken(token)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
				}
				// 将用户信息存储到上下文
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("appID", claims.AppID)
			default:
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unknown token type"})
			}

			return next(c)
		}
	}
}
