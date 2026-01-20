package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
			// 1. 获取 Token (使用自定义头 X-Authos-Token 以避免与业务 Authorization 冲突)
			authHeader := c.Request().Header.Get("X-Authos-Token")
			
			// 兼容性逻辑：如果 X-Authos-Token 为空，尝试从旧的 Header 获取 (为了平滑过渡)
			if authHeader == "" {
				if sysToken := c.Request().Header.Get("X-System-Token"); sysToken != "" {
					authHeader = sysToken
				} else if appToken := c.Request().Header.Get("X-App-Token"); appToken != "" {
					authHeader = appToken
				}
			}

			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "X-Authos-Token header is required"})
			}

			// 去除 Bearer 前缀
			tokenString := authHeader
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}

			// 2. 解析 Token (使用 MapClaims 以支持多种类型)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(j.JWTConfig.SecretKey), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
			}

			// 3. 识别 Token 类型
			tokenType, _ := claims["type"].(string)

			// 如果没有 type 字段 (旧 Token)，尝试通过特征推断
			if tokenType == "" {
				if _, ok := claims["isAdmin"]; ok {
					tokenType = "system"
				} else if _, ok := claims["appCode"]; ok && claims["username"] == nil {
					tokenType = "app"
				} else {
					tokenType = "user"
				}
			}

			// 4. 根据类型设置上下文
			switch tokenType {
			case "system":
				username, _ := claims["username"].(string)
				c.Set("isSystemAdmin", true)
				c.Set("username", username)
				
				// System Admin 需要从 Header 获取目标 AppID
				appIDStr := c.Request().Header.Get("X-App-ID")
				if appIDStr != "" {
					var appID uint
					if _, err := fmt.Sscanf(appIDStr, "%d", &appID); err == nil && appID > 0 {
						c.Set("appID", appID)
					}
				}

			case "app":
				appIDFloat, _ := claims["appId"].(float64) // JSON 数字通常解析为 float64
				appCode, _ := claims["appCode"].(string)
				
				c.Set("isAppToken", true)
				c.Set("appID", uint(appIDFloat))
				c.Set("appCode", appCode)

			case "user":
				userIDFloat, _ := claims["userId"].(float64)
				username, _ := claims["username"].(string)
				appIDFloat, _ := claims["appId"].(float64)
				
				c.Set("userID", uint(userIDFloat))
				c.Set("username", username)
				c.Set("appID", uint(appIDFloat))

			default:
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unknown token type"})
			}

			return next(c)
		}
	}
}
