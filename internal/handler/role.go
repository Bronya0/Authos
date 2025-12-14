package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Authos/internal/model"
	"Authos/internal/service"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	RoleService *service.RoleService
}

// NewRoleHandler 创建角色处理器实例
func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{RoleService: roleService}
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c echo.Context) error {
	var role model.Role
	if err := c.Bind(&role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 数据验证
	if role.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Role name is required"})
	}
	if len(role.Name) > 50 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Role name cannot exceed 50 characters"})
	}

	if err := h.RoleService.CreateRole(&role); err != nil {
		// 记录详细错误信息
		log.Printf("Failed to create role: %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Failed to create role: %v", err)})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"role":    role,
		"message": "Role created successfully",
	})
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	var role model.Role
	if err := c.Bind(&role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	role.ID = uint(id)

	if err := h.RoleService.UpdateRole(&role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update role"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"role":    role,
		"message": "Role updated successfully",
	})
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	if err := h.RoleService.DeleteRole(uint(id), appID); err != nil {
		// Log the actual error for debugging
		log.Printf("Error deleting role %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Failed to delete role: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Role deleted successfully"})
}

// GetRole 获取角色
func (h *RoleHandler) GetRole(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	role, err := h.RoleService.GetRoleByID(uint(id), appID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Role not found"})
	}

	return c.JSON(http.StatusOK, role)
}

// ListRoles 列出所有角色
func (h *RoleHandler) ListRoles(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	roles, err := h.RoleService.ListRolesByApp(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get roles"})
	}

	return c.JSON(http.StatusOK, roles)
}

// AssignMenusRequest 分配菜单请求
type AssignMenusRequest struct {
	MenuIDs []uint `json:"menuIds" binding:"required"`
}

// AssignMenus 为角色分配菜单
func (h *RoleHandler) AssignMenus(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	var req AssignMenusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := h.RoleService.AssignMenus(uint(id), appID, req.MenuIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to assign menus"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Menus assigned successfully"})
}

// GetRoleMenus 获取角色菜单
func (h *RoleHandler) GetRoleMenus(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	menus, err := h.RoleService.GetRoleMenus(uint(id), appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get role menus"})
	}

	return c.JSON(http.StatusOK, menus)
}

// UpdateRoleMenus 更新角色菜单
func (h *RoleHandler) UpdateRoleMenus(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	var req AssignMenusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := h.RoleService.AssignMenus(uint(id), appID, req.MenuIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update role menus"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Role menus updated successfully"})
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	Permissions []map[string]string `json:"permissions" binding:"required"`
}

// AssignPermissions 为角色分配 API 权限
func (h *RoleHandler) AssignPermissions(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	var req AssignPermissionsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := h.RoleService.AssignPermissions(uint(id), appID, req.Permissions); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to assign permissions"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Permissions assigned successfully"})
}

// UpdatePermissions 更新角色 API 权限 (PUT方法)
func (h *RoleHandler) UpdatePermissions(c echo.Context) error {
	return h.AssignPermissions(c)
}
