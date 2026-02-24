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

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}
	role.AppID = appID

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

	// 记录审计日志
	userIDInterface := c.Get("userID")
	usernameInterface := c.Get("username")
	var userID uint
	var username string

	if userIDInterface != nil {
		if u, ok := userIDInterface.(uint); ok {
			userID = u
		} else if f, ok := userIDInterface.(float64); ok {
			userID = uint(f)
		}
	}

	if usernameInterface != nil {
		if s, ok := usernameInterface.(string); ok {
			username = s
		}
	}

	h.RoleService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "CREATE",
		Resource:   "ROLE",
		ResourceID: fmt.Sprintf("%d", role.ID),
		Content:    fmt.Sprintf("创建角色: %s", role.Name),
		IP:         c.RealIP(),
		Status:     1,
	})

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

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	var role model.Role
	if err := c.Bind(&role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	role.ID = uint(id)
	role.AppID = appID

	if err := h.RoleService.UpdateRole(&role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update role"})
	}

	// 记录审计日志
	userIDInterface := c.Get("userID")
	usernameInterface := c.Get("username")
	var userID uint
	var username string
	if userIDInterface != nil {
		if u, ok := userIDInterface.(uint); ok {
			userID = u
		} else if f, ok := userIDInterface.(float64); ok {
			userID = uint(f)
		}
	}
	if usernameInterface != nil {
		if s, ok := usernameInterface.(string); ok {
			username = s
		}
	}
	h.RoleService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "UPDATE",
		Resource:   "ROLE",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("更新角色: %s", role.Name),
		IP:         c.RealIP(),
		Status:     1,
	})

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

	// 记录审计日志
	userIDInterface := c.Get("userID")
	usernameInterface := c.Get("username")
	var userID uint
	var username string

	if userIDInterface != nil {
		if u, ok := userIDInterface.(uint); ok {
			userID = u
		} else if f, ok := userIDInterface.(float64); ok {
			userID = uint(f)
		}
	}

	if usernameInterface != nil {
		if s, ok := usernameInterface.(string); ok {
			username = s
		}
	}

	h.RoleService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "DELETE",
		Resource:   "ROLE",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("删除角色ID: %d", id),
		IP:         c.RealIP(),
		Status:     1,
	})

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

	name := c.QueryParam("name")

	db := h.RoleService.DB.Preload("Menus").Where("app_id = ?", appID)
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	var roles []*model.Role
	if err := db.Order("id asc").Find(&roles).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get roles"})
	}

	for _, role := range roles {
		// 菜单信息
		role.MenuCount = len(role.Menus)
		previewCount := 3
		if role.MenuCount < previewCount {
			previewCount = role.MenuCount
		}
		role.MenuPreview = make([]string, 0, previewCount)
		for i := 0; i < previewCount; i++ {
			role.MenuPreview = append(role.MenuPreview, role.Menus[i].Name)
		}

		// API 权限信息
		if role.IsSuperAdmin {
			// 超级管理员拥有所有权限
			var allPermissions []model.ApiPermission
			if err := h.RoleService.DB.Where("app_id = ?", appID).Find(&allPermissions).Error; err == nil {
				role.ApiPermCount = len(allPermissions)

				apiPreviewCount := 3
				if role.ApiPermCount < apiPreviewCount {
					apiPreviewCount = role.ApiPermCount
				}
				role.ApiPermPreview = make([]string, 0, apiPreviewCount)
				for i := 0; i < apiPreviewCount; i++ {
					role.ApiPermPreview = append(role.ApiPermPreview, allPermissions[i].Name)
				}
			} else {
				role.ApiPermCount = 0
				role.ApiPermPreview = []string{}
			}
		} else {
			// 普通角色通过 Casbin 获取权限
			roleKey := fmt.Sprintf("role:%s", role.UUID)
			if h.RoleService.CasbinService != nil && h.RoleService.CasbinService.Enforcer != nil {
				policies, _ := h.RoleService.CasbinService.Enforcer.GetFilteredPolicy(0, roleKey)
				role.ApiPermCount = len(policies)

				apiPreviewCount := 3
				if role.ApiPermCount < apiPreviewCount {
					apiPreviewCount = role.ApiPermCount
				}
				role.ApiPermPreview = make([]string, 0, apiPreviewCount)
				for i := 0; i < apiPreviewCount; i++ {
					// 检查 policies[i] 的长度，防止索引越界
					if len(policies[i]) > 2 {
						role.ApiPermPreview = append(role.ApiPermPreview, fmt.Sprintf("%s %s", policies[i][2], policies[i][1]))
					} else if len(policies[i]) > 1 {
						role.ApiPermPreview = append(role.ApiPermPreview, policies[i][1])
					}
				}
			} else {
				// Casbin 服务未初始化，记录日志并跳过权限信息
				log.Printf("Warning: CasbinService or Enforcer is nil when listing roles")
				role.ApiPermCount = 0
				role.ApiPermPreview = []string{}
			}
		}
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

	// 记录审计日志
	userIDInterface := c.Get("userID")
	usernameInterface := c.Get("username")
	var userID uint
	var username string

	if userIDInterface != nil {
		if u, ok := userIDInterface.(uint); ok {
			userID = u
		} else if f, ok := userIDInterface.(float64); ok {
			userID = uint(f)
		}
	}

	if usernameInterface != nil {
		if s, ok := usernameInterface.(string); ok {
			username = s
		}
	}

	h.RoleService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "UPDATE",
		Resource:   "ROLE_PERMISSION",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("分配角色接口权限, 角色ID: %d, 权限数量: %d", id, len(req.Permissions)),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Permissions assigned successfully"})
}

// UpdatePermissions 更新角色 API 权限 (PUT方法)
func (h *RoleHandler) UpdatePermissions(c echo.Context) error {
	return h.AssignPermissions(c)
}
