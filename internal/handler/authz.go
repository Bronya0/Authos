package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"Authos/internal/model"
	"Authos/internal/service"
)

// AuthzHandler 权限处理器
type AuthzHandler struct {
	CasbinService        *service.CasbinService
	MenuService          *service.MenuService
	ApplicationService   *service.ApplicationService
	ApiPermissionService *service.ApiPermissionService
	JWTConfig            *service.JWTConfig
}

// NewAuthzHandler 创建权限处理器实例
func NewAuthzHandler(casbinService *service.CasbinService, menuService *service.MenuService, applicationService *service.ApplicationService, apiPermissionService *service.ApiPermissionService, jwtConfig *service.JWTConfig) *AuthzHandler {
	return &AuthzHandler{
		CasbinService:        casbinService,
		MenuService:          menuService,
		ApplicationService:   applicationService,
		ApiPermissionService: apiPermissionService,
		JWTConfig:            jwtConfig,
	}
}

// CheckPermissionReq 权限检查请求
type CheckPermissionReq struct {
	UserID uint   `json:"userId" binding:"required"`
	Obj    string `json:"obj" binding:"required"`
	Act    string `json:"act" binding:"required"`
}

type CheckPermissionByKeyReq struct {
	UserID     uint   `json:"userId" binding:"required"`
	Permission string `json:"permission" binding:"required"`
}

// CheckPermission 检查权限
func (h *AuthzHandler) CheckPermission(c echo.Context) error {
	var req CheckPermissionReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 调用 Casbin 检查权限
	allowed, err := h.CasbinService.CheckPermission(req.UserID, req.Obj, req.Act)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to check permission"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"allowed": allowed,
		"message": "Permission checked successfully",
	})
}

// CheckPermissionWithSecretReq 统一鉴权请求（带Secret）
type CheckPermissionWithSecretReq struct {
	AppCode   string `json:"appCode" binding:"required"`
	AppSecret string `json:"appSecret" binding:"required"`
	Token     string `json:"token" binding:"required"`
	Obj       string `json:"obj" binding:"required"` // 访问路径
	Act       string `json:"act" binding:"required"` // 访问方法
}

// CheckPermissionWithSecret 统一鉴权接口
func (h *AuthzHandler) CheckPermissionWithSecret(c echo.Context) error {
	var req CheckPermissionWithSecretReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 1. 验证应用身份
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

	// 2. 解析并验证 Token
	claims, err := h.JWTConfig.ParseToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}

	// 验证 Token 是否属于该应用
	if claims.AppID != app.ID {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token does not belong to this application"})
	}

	// 3. 根据路径和方法解析对应的接口权限（支持 * 通配方法）
	permission, err := h.ApiPermissionService.GetApiPermissionByPathAndMethod(app.ID, req.Obj, req.Act)
	if err != nil {
		// 如果未找到对应权限，直接视为无权限
		log.Printf("permission not found for appID=%d path=%s method=%s: %v", app.ID, req.Obj, req.Act, err)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"allowed": false,
			"userId":  claims.UserID,
			"message": "Permission not found",
		})
	}

	// 4. 使用权限标识 + 请求方法 交给 Casbin 检查（策略里方法为 * 时也可匹配）
	log.Printf("Checking permission for userID: %d, key: %s, path: %s, act: %s", claims.UserID, permission.Key, req.Obj, req.Act)
	allowed, err := h.CasbinService.CheckPermission(claims.UserID, permission.Key, req.Act)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to check permission"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"allowed": allowed,
		"userId":  claims.UserID,
		"message": "Permission checked successfully",
	})
}

func (h *AuthzHandler) CheckPermissionByKey(c echo.Context) error {
	var req CheckPermissionByKeyReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	allowed, err := h.CasbinService.CheckPermission(req.UserID, req.Permission, model.HTTP_ALL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to check permission"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"allowed": allowed,
		"message": "Permission checked successfully",
	})
}

// GetUserNav 获取用户导航菜单
func (h *AuthzHandler) GetUserNav(c echo.Context) error {
	// 从上下文获取用户ID
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not authenticated"})
	}

	// 获取用户菜单树
	menuTree, err := h.MenuService.GetUserMenuTree(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get user menu"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"nav":     menuTree,
		"message": "User menu retrieved successfully",
	})
}
