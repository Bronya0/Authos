package handler

import (
	"net/http"
	"strconv"

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

	permissions, err := h.ApiPermissionService.GetAllApiPermissions(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取接口权限列表失败"})
	}

	return c.JSON(http.StatusOK, permissions)
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
		Name        string `json:"name"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "请求参数错误"})
	}

	permission, err := h.ApiPermissionService.CreateApiPermission(appID, req.Name, req.Path, req.Method, req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

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
		Name        string `json:"name"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "请求参数错误"})
	}

	permission, err := h.ApiPermissionService.UpdateApiPermission(uint(id), appID, req.Name, req.Path, req.Method, req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

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

	return c.JSON(http.StatusOK, map[string]string{"message": "移除角色接口权限成功"})
}
