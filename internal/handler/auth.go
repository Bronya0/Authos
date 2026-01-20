package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"Authos/internal/model"
	"Authos/internal/service"
	"Authos/pkg/utils"
)

// getAppIDFromToken 从 JWT token 或请求头中获取应用ID
func getAppIDFromToken(c echo.Context) (uint, error) {

	// 尝试从上下文获取 appID（JWT token）
	var tokenAppID uint
	appIDInterface := c.Get("appID")
	if appIDInterface != nil {
		// 尝试类型断言
		switch v := appIDInterface.(type) {
		case uint:
			tokenAppID = v
		case float64:
			tokenAppID = uint(v)
		case int:
			tokenAppID = uint(v)
		default:
			fmt.Printf("getAppIDFromToken: Type assertion failed for appID from token, actual type: %T\n", appIDInterface)
		}
	}

	// 尝试从请求头获取 appID
	var headerAppID uint
	appIDStr := c.Request().Header.Get("X-App-ID")
	if appIDStr != "" {
		if _, err := fmt.Sscanf(appIDStr, "%d", &headerAppID); err == nil && headerAppID > 0 {
			fmt.Printf("getAppIDFromToken: Got appID from header: %d\n", headerAppID)
		}
	}

	// 逻辑判断：
	// 1. 如果 Token 中的 AppID 为 1 (系统管理员)，且 Header 中有指定的 AppID，则使用 Header 中的 AppID
	// 2. 否则，优先使用 Token 中的 AppID
	// 3. 如果都没有，则返回错误

	if tokenAppID == 1 && headerAppID > 0 {
		fmt.Printf("getAppIDFromToken: System admin switching context to appID: %d\n", headerAppID)
		return headerAppID, nil
	}

	if tokenAppID > 0 {
		return tokenAppID, nil
	}

	if headerAppID > 0 {
		// 兼容逻辑：如果没有 Token 但有 Header (且未被鉴权中间件拦截)，则使用 Header
		return headerAppID, nil
	}

	fmt.Printf("getAppIDFromToken: appID not found in token or header\n")
	return 0, echo.NewHTTPError(http.StatusUnauthorized, "App ID not found in token or header")
}

// AuthHandler 认证处理器
type AuthHandler struct {
	UserService        *service.UserService
	ApplicationService *service.ApplicationService
	AuditLogService    *service.AuditLogService
	JWTConfig          *utils.JWTConfig
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService *service.UserService, applicationService *service.ApplicationService, auditLogService *service.AuditLogService, jwtConfig *utils.JWTConfig) *AuthHandler {
	return &AuthHandler{
		UserService:        userService,
		ApplicationService: applicationService,
		AuditLogService:    auditLogService,
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

// ProxyLoginRequest 代理登录请求
type ProxyLoginRequest struct {
	AppCode   string `json:"appCode" binding:"required"`
	AppSecret string `json:"appSecret" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// SystemLogin 系统管理员登录接口
func (h *AuthHandler) SystemLogin(c echo.Context) error {
	var req SystemLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 添加调试日志
	fmt.Printf("SystemLogin attempt: username=%s\n", req.Username)

	// 查询系统管理员账号（username为admin）
	// 系统管理员不与特定应用关联，appID为1（系统默认应用）
	user, err := h.UserService.GetUserByUsernameForSystem(req.Username)
	if err != nil {
		fmt.Printf("SystemLogin: user not found or not admin: %v\n", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid admin credentials"})
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Printf("SystemLogin: password mismatch\n")
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid admin credentials"})
	}

	// 获取应用信息（使用user的AppID）
	var app model.Application
	if err := h.UserService.DB.First(&app, user.AppID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Application not found"})
	}

	// 生成系统管理员JWT令牌（包含应用信息）
	token, err := h.JWTConfig.GenerateToken(user.ID, user.Username, app.ID, app.UUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	// 记录审计日志
	h.AuditLogService.Record(&model.AuditLog{
		AppID:    0, // 系统管理员登录记录为系统日志
		UserID:   user.ID,
		Username: user.Username,
		Action:   "SYSTEM_LOGIN",
		Resource: "APPLICATION",
		IP:       c.RealIP(),
		Status:   1,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"user":    user,
		"app":     app,
		"message": "System login successful",
	})
}

// AppLogin 应用登录接口
func (h *AuthHandler) AppLogin(c echo.Context) error {
	var req struct {
		AppUUID   string `json:"appUuid" binding:"required"` // 使用UUID而不是AppID
		AppSecret string `json:"appSecret" binding:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 获取应用信息（通过UUID）
	app, err := h.ApplicationService.GetApplicationByUUID(req.AppUUID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid application UUID"})
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
	token, err := h.JWTConfig.GenerateAppToken(app.ID, app.UUID, app.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	// 记录审计日志
	h.AuditLogService.Record(&model.AuditLog{
		AppID:    app.ID,
		Action:   "APP_LOGIN",
		Resource: "APPLICATION",
		IP:       c.RealIP(),
		Status:   1,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"app":     app,
		"message": "App login successful",
	})
}

// ProxyLogin 代理登录接口（后端透传模式）
func (h *AuthHandler) ProxyLogin(c echo.Context) error {
	var req ProxyLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 1. 验证应用身份 (AppCode + Secret)
	app, err := h.ApplicationService.GetApplicationByCode(req.AppCode)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid application code"})
	}

	if app.SecretKey != req.AppSecret {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid application secret"})
	}

	if app.Status == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Application is disabled"})
	}

	// 2. 验证用户身份
	user, err := h.UserService.GetUserByUsername(req.Username, app.ID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	if user.Status == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User is disabled"})
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	// 3. 生成 Token
	token, err := h.JWTConfig.GenerateToken(user.ID, user.Username, app.ID, app.UUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	// 4. 记录审计日志
	h.AuditLogService.Record(&model.AuditLog{
		AppID:    app.ID,
		UserID:   user.ID,
		Username: user.Username,
		Action:   "PROXY_LOGIN",
		Resource: "USER",
		IP:       c.RealIP(),
		Status:   1,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"user":    user,
		"app":     app,
		"message": "Proxy login successful",
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

	// 生成JWT令牌（包含应用ID和UUID）
	token, err := h.JWTConfig.GenerateToken(user.ID, user.Username, app.ID, app.UUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	// 记录审计日志
	h.AuditLogService.Record(&model.AuditLog{
		AppID:    app.ID,
		UserID:   user.ID,
		Username: user.Username,
		Action:   "LOGIN",
		Resource: "USER",
		IP:       c.RealIP(),
		Status:   1,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"user":    user,
		"app":     app,
		"message": "Login successful",
	})
}

// Logout 登出接口
func (h *AuthHandler) Logout(c echo.Context) error {
	// 尝试记录审计日志
	appIDInterface := c.Get("appID")
	userIDInterface := c.Get("userID")
	usernameInterface := c.Get("username")

	if appIDInterface != nil && userIDInterface != nil {
		h.AuditLogService.Record(&model.AuditLog{
			AppID:    appIDInterface.(uint),
			UserID:   userIDInterface.(uint),
			Username: usernameInterface.(string),
			Action:   "LOGOUT",
			Resource: "USER",
			IP:       c.RealIP(),
			Status:   1,
		})
	}

	// JWT是无状态的，登出只需客户端删除令牌即可
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}

// GetDashboardStats 获取仪表盘统计数据
func (h *AuthHandler) GetDashboardStats(c echo.Context) error {
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	var userCount int64
	var roleCount int64
	var menuCount int64
	var apiCount int64
	var auditCount int64

	h.UserService.DB.Model(&model.User{}).Where("app_id = ?", appID).Count(&userCount)
	h.UserService.DB.Model(&model.Role{}).Where("app_id = ?", appID).Count(&roleCount)
	h.UserService.DB.Model(&model.Menu{}).Where("app_id = ?", appID).Count(&menuCount)
	h.UserService.DB.Model(&model.ApiPermission{}).Where("app_id = ?", appID).Count(&apiCount)
	h.UserService.DB.Model(&model.AuditLog{}).Where("app_id = ?", appID).Count(&auditCount)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users":     userCount,
		"roles":     roleCount,
		"menus":     menuCount,
		"apiPerms":  apiCount,
		"auditLogs": auditCount,
	})
}
