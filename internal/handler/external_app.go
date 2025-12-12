package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"Authos/internal/model"
	"Authos/internal/service"
)

// ExternalAppHandler 外部应用处理器
type ExternalAppHandler struct {
	externalAppService *service.ExternalAppService
}

// NewExternalAppHandler 创建外部应用处理器实例
func NewExternalAppHandler(externalAppService *service.ExternalAppService) *ExternalAppHandler {
	return &ExternalAppHandler{
		externalAppService: externalAppService,
	}
}

// CreateExternalApp 创建外部应用
func (h *ExternalAppHandler) CreateExternalApp(c echo.Context) error {
	var app model.ExternalApp
	if err := c.Bind(&app); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.externalAppService.CreateExternalApp(&app); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "创建应用失败"})
	}

	return c.JSON(http.StatusOK, app)
}

// ListExternalApps 获取外部应用列表
func (h *ExternalAppHandler) ListExternalApps(c echo.Context) error {
	apps, err := h.externalAppService.ListExternalApps()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取应用列表失败"})
	}

	return c.JSON(http.StatusOK, apps)
}

// GetExternalApp 获取外部应用详情
func (h *ExternalAppHandler) GetExternalApp(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的应用ID"})
	}

	app, err := h.externalAppService.GetExternalApp(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "应用不存在"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取应用详情失败"})
	}

	return c.JSON(http.StatusOK, app)
}

// UpdateExternalApp 更新外部应用
func (h *ExternalAppHandler) UpdateExternalApp(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的应用ID"})
	}

	var app model.ExternalApp
	if err := c.Bind(&app); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.externalAppService.UpdateExternalApp(uint(id), &app); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "更新应用失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "应用更新成功"})
}

// DeleteExternalApp 删除外部应用
func (h *ExternalAppHandler) DeleteExternalApp(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的应用ID"})
	}

	if err := h.externalAppService.DeleteExternalApp(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "删除应用失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "应用删除成功"})
}

// CreateAppPermission 创建应用权限
func (h *ExternalAppHandler) CreateAppPermission(c echo.Context) error {
	var permission model.ExternalAppPermission
	if err := c.Bind(&permission); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.externalAppService.CreateAppPermission(&permission); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "创建权限失败"})
	}

	return c.JSON(http.StatusOK, permission)
}

// GetAppPermissions 获取应用权限列表
func (h *ExternalAppHandler) GetAppPermissions(c echo.Context) error {
	appID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的应用ID"})
	}

	permissions, err := h.externalAppService.GetAppPermissions(uint(appID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取权限列表失败"})
	}

	return c.JSON(http.StatusOK, permissions)
}

// UpdateAppPermission 更新应用权限
func (h *ExternalAppHandler) UpdateAppPermission(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的权限ID"})
	}

	var permission model.ExternalAppPermission
	if err := c.Bind(&permission); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.externalAppService.UpdateAppPermission(uint(id), &permission); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "更新权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "权限更新成功"})
}

// DeleteAppPermission 删除应用权限
func (h *ExternalAppHandler) DeleteAppPermission(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的权限ID"})
	}

	if err := h.externalAppService.DeleteAppPermission(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "删除权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "权限删除成功"})
}

// GenerateAppToken 生成应用令牌
func (h *ExternalAppHandler) GenerateAppToken(c echo.Context) error {
	appID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的应用ID"})
	}

	// 默认令牌有效期为24小时
	expiresAt := time.Now().Add(24 * time.Hour)

	token, err := h.externalAppService.CreateAppToken(uint(appID), expiresAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "生成令牌失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_at": expiresAt,
	})
}