package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"Authos/internal/model"
	"Authos/internal/service"

	"github.com/labstack/echo/v4"
)

// ApiPermissionHandler 接口权限处理器
type ApiPermissionHandler struct {
	ApiPermissionService *service.ApiPermissionService
}

// NewApiPermissionHandler 创建接口权限处理器
func NewApiPermissionHandler(apiPermissionService *service.ApiPermissionService) *ApiPermissionHandler {
	return &ApiPermissionHandler{
		ApiPermissionService: apiPermissionService,
	}
}

// ListApiPermissions 列出所有接口权限
func (h *ApiPermissionHandler) ListApiPermissions(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	name := c.QueryParam("name")
	path := c.QueryParam("path")
	method := c.QueryParam("method")

	db := h.ApiPermissionService.DB.Where("app_id = ?", appID)
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if path != "" {
		db = db.Where("path LIKE ?", "%"+path+"%")
	}
	if method != "" {
		db = db.Where("method = ?", method)
	}

	var perms []*model.ApiPermission
	if err := db.Order("id asc").Find(&perms).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取接口权限列表失败"})
	}

	return c.JSON(http.StatusOK, perms)
}

// GetApiPermission 获取接口权限
func (h *ApiPermissionHandler) GetApiPermission(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的权限ID"})
	}

	permission, err := h.ApiPermissionService.GetApiPermission(uint(id), appID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "接口权限不存在"})
	}

	return c.JSON(http.StatusOK, permission)
}

// CreateApiPermission 创建接口权限
func (h *ApiPermissionHandler) CreateApiPermission(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	var req struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "请求参数错误"})
	}

	permission, err := h.ApiPermissionService.CreateApiPermission(appID, req.Key, req.Name, req.Path, req.Method, req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
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

	h.ApiPermissionService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "CREATE",
		Resource:   "API_PERMISSION",
		ResourceID: fmt.Sprintf("%d", permission.ID),
		Content:    fmt.Sprintf("创建接口权限: %s", permission.Name),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusCreated, permission)
}

// UpdateApiPermission 更新接口权限
func (h *ApiPermissionHandler) UpdateApiPermission(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的权限ID"})
	}

	var req struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "请求参数错误"})
	}

	permission, err := h.ApiPermissionService.UpdateApiPermission(uint(id), appID, req.Key, req.Name, req.Path, req.Method, req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
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

	h.ApiPermissionService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "UPDATE",
		Resource:   "API_PERMISSION",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("更新接口权限: %s", permission.Name),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, permission)
}

// DeleteApiPermission 删除接口权限
func (h *ApiPermissionHandler) DeleteApiPermission(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的权限ID"})
	}

	if err := h.ApiPermissionService.DeleteApiPermission(uint(id), appID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	// 记录审计日志
	h.ApiPermissionService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     c.Get("userID").(uint),
		Username:   c.Get("username").(string),
		Action:     "DELETE",
		Resource:   "API_PERMISSION",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("删除接口权限ID: %d", id),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "接口权限删除成功"})
}

// GetApiPermissionsForRole 获取角色的接口权限
func (h *ApiPermissionHandler) GetApiPermissionsForRole(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	roleUUID := c.Param("roleUUID")
	if roleUUID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "角色UUID不能为空"})
	}

	permissions, err := h.ApiPermissionService.GetApiPermissionsForRole(appID, roleUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取角色接口权限失败"})
	}

	return c.JSON(http.StatusOK, permissions)
}

// AddApiPermissionToRole 为角色添加接口权限
func (h *ApiPermissionHandler) AddApiPermissionToRole(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	roleUUID := c.Param("roleUUID")
	if roleUUID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "角色UUID不能为空"})
	}

	var req struct {
		PermissionUUID string `json:"permissionUUID"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "请求参数错误"})
	}

	if err := h.ApiPermissionService.AddApiPermissionToRole(appID, roleUUID, req.PermissionUUID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	// 记录审计日志
	h.ApiPermissionService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     c.Get("userID").(uint),
		Username:   c.Get("username").(string),
		Action:     "ASSIGN",
		Resource:   "ROLE_PERMISSION",
		ResourceID: roleUUID,
		Content:    fmt.Sprintf("为角色分配权限: %s", req.PermissionUUID),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "为角色添加接口权限成功"})
}

// RemoveApiPermissionFromRole 移除角色的接口权限
func (h *ApiPermissionHandler) RemoveApiPermissionFromRole(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	roleUUID := c.Param("roleUUID")
	if roleUUID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "角色UUID不能为空"})
	}

	var req struct {
		PermissionUUID string `json:"permissionUUID"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "请求参数错误"})
	}

	if err := h.ApiPermissionService.RemoveApiPermissionFromRole(appID, roleUUID, req.PermissionUUID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	// 记录审计日志
	h.ApiPermissionService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     c.Get("userID").(uint),
		Username:   c.Get("username").(string),
		Action:     "UNASSIGN",
		Resource:   "ROLE_PERMISSION",
		ResourceID: roleUUID,
		Content:    fmt.Sprintf("移除角色权限: %s", req.PermissionUUID),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "移除角色接口权限成功"})
}
