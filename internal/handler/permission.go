package handler

import (
	"net/http"
	"strconv"

	"Authos/internal/service"

	"github.com/labstack/echo/v4"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	PermissionService *service.PermissionService
	RoleService       *service.RoleService
}

// NewPermissionHandler 创建权限处理器实例
func NewPermissionHandler(permissionService *service.PermissionService, roleService *service.RoleService) *PermissionHandler {
	return &PermissionHandler{
		PermissionService: permissionService,
		RoleService:       roleService,
	}
}

// ListPermissions 获取权限列表
func (h *PermissionHandler) ListPermissions(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	permissions, err := h.PermissionService.GetAllPermissions(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取权限列表失败"})
	}

	return c.JSON(http.StatusOK, permissions)
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c echo.Context) error {
	var req struct {
		Obj         string `json:"obj" binding:"required"`
		Act         string `json:"act" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.PermissionService.CreatePermission(req.Obj, req.Act, req.Description); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "创建权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "权限创建成功"})
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c echo.Context) error {
	obj := c.QueryParam("obj")
	act := c.QueryParam("act")

	if obj == "" || act == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "缺少必要参数"})
	}

	if err := h.PermissionService.DeletePermission(obj, act); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "删除权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "权限删除成功"})
}

// GetPermissionRoles 获取权限关联的角色
func (h *PermissionHandler) GetPermissionRoles(c echo.Context) error {
	obj := c.QueryParam("obj")
	act := c.QueryParam("act")

	if obj == "" || act == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "缺少必要参数"})
	}

	roles, err := h.PermissionService.GetRolesForPermission(obj, act)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取权限角色失败"})
	}

	return c.JSON(http.StatusOK, roles)
}

// AssignPermissionToRole 为角色分配权限
func (h *PermissionHandler) AssignPermissionToRole(c echo.Context) error {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的角色ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	var req []struct {
		Obj string `json:"obj" binding:"required"`
		Act string `json:"act" binding:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	// 转换为角色服务需要的格式
	var permissions []map[string]string
	for _, perm := range req {
		permissions = append(permissions, map[string]string{
			"obj": perm.Obj,
			"act": perm.Act,
		})
	}

	if err := h.RoleService.AssignPermissions(uint(roleID), appID, permissions); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "分配权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "权限分配成功"})
}
