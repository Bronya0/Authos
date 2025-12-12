package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"Authos/internal/service"
	"Authos/pkg/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	UserService *service.UserService
	JWTConfig   *utils.JWTConfig
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService *service.UserService, jwtConfig *utils.JWTConfig) *AuthHandler {
	return &AuthHandler{
		UserService: userService,
		JWTConfig:   jwtConfig,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录接口
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 查找用户
	user, err := h.UserService.GetUserByUsername(req.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	// 检查用户状态
	if user.Status == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User is disabled"})
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	// 生成JWT令牌
	token, err := h.JWTConfig.GenerateToken(user.ID, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":  token,
		"user":   user,
		"message": "Login successful",
	})
}

// Logout 登出接口
func (h *AuthHandler) Logout(c echo.Context) error {
	// JWT是无状态的，登出只需客户端删除令牌即可
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}
