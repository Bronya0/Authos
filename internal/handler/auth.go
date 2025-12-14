package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"Authos/internal/service"
	"Authos/pkg/utils"
)

// getAppIDFromToken 从 JWT token 中获取应用ID
func getAppIDFromToken(c echo.Context) (uint, error) {
	appID, ok := c.Get("appID").(uint)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "App ID not found in token")
	}
	return appID, nil
}

// AuthHandler 认证处理器
type AuthHandler struct {
	UserService        *service.UserService
	ApplicationService *service.ApplicationService
	JWTConfig          *utils.JWTConfig
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService *service.UserService, applicationService *service.ApplicationService, jwtConfig *utils.JWTConfig) *AuthHandler {
	return &AuthHandler{
		UserService:        userService,
		ApplicationService: applicationService,
		JWTConfig:          jwtConfig,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	AppCode  string `json:"appCode" binding:"required"` // 应用代码，用于多租户
}

// SystemLoginRequest 系统管理员登录请求
type SystemLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AppLoginRequest 应用登录请求
type AppLoginRequest struct {
	AppID     uint   `json:"appId" binding:"required"`
	AppSecret string `json:"appSecret" binding:"required"`
}

// SystemLogin 系统管理员登录接口
func (h *AuthHandler) SystemLogin(c echo.Context) error {
	var req SystemLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 添加调试日志
	fmt.Printf("SystemLogin attempt: username=%s, password=%s\n", req.Username, req.Password)

	// 验证系统管理员账号（硬编码的内置管理员）
	// 注意：此处使用硬编码密码仅用于开发/测试环境，生产环境应使用数据库存储或配置文件
	if req.Username != "admin" || req.Password != "admin123" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid admin credentials"})
	}

	// 生成系统管理员JWT令牌
	token, err := h.JWTConfig.GenerateSystemToken(req.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"user":    map[string]interface{}{"username": req.Username},
		"message": "System login successful",
	})
}

// AppLogin 应用登录接口
func (h *AuthHandler) AppLogin(c echo.Context) error {
	var req AppLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 获取应用信息
	app, err := h.ApplicationService.GetApplicationByID(fmt.Sprintf("%d", req.AppID))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid application ID"})
	}

	// 检查应用密钥
	if app.SecretKey != req.AppSecret {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid application secret"})
	}

	// 检查应用状态
	if app.Status == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Application is disabled"})
	}

	// 生成应用JWT令牌
	token, err := h.JWTConfig.GenerateAppToken(app.ID, app.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"app":     app,
		"message": "App login successful",
	})
}

// Login 登录接口（多租户）
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 获取应用信息
	app, err := h.UserService.GetApplicationByCode(req.AppCode)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid application code"})
	}

	// 检查应用状态
	if app.Status == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Application is disabled"})
	}

	// 查找用户（按应用隔离）
	user, err := h.UserService.GetUserByUsername(req.Username, app.ID)
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

	// 生成JWT令牌（包含应用ID）
	token, err := h.JWTConfig.GenerateToken(user.ID, user.Username, app.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"user":    user,
		"app":     app,
		"message": "Login successful",
	})
}

// Logout 登出接口
func (h *AuthHandler) Logout(c echo.Context) error {
	// JWT是无状态的，登出只需客户端删除令牌即可
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}
