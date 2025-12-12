package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"Authos/internal/service"
)

// AuthzHandler 权限处理器
type AuthzHandler struct {
	CasbinService *service.CasbinService
	MenuService   *service.MenuService
}

// NewAuthzHandler 创建权限处理器实例
func NewAuthzHandler(casbinService *service.CasbinService, menuService *service.MenuService) *AuthzHandler {
	return &AuthzHandler{
		CasbinService: casbinService,
		MenuService:   menuService,
	}
}

// CheckPermissionReq 权限检查请求
type CheckPermissionReq struct {
	UserID uint   `json:"userId" binding:"required"`
	Obj    string `json:"obj" binding:"required"`
	Act    string `json:"act" binding:"required"`
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
