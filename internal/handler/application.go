package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"Authos/internal/service"
)

// ApplicationHandler 应用处理器
type ApplicationHandler struct {
	ApplicationService *service.ApplicationService
}

// NewApplicationHandler 创建应用处理器实例
func NewApplicationHandler(applicationService *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		ApplicationService: applicationService,
	}
}

// CreateApplicationRequest 创建应用请求
type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

// UpdateApplicationRequest 更新应用请求
type UpdateApplicationRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

// CreateApplication 创建应用
func (h *ApplicationHandler) CreateApplication(c echo.Context) error {
	var req CreateApplicationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 创建应用
	app, err := h.ApplicationService.CreateApplication(req.Name, req.Code, req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"app":      app,
		"message":  "Application created successfully",
		"secretKey": app.SecretKey, // 注意：仅在创建时返回密钥
	})
}

// ListApplications 列出所有应用
func (h *ApplicationHandler) ListApplications(c echo.Context) error {
	apps, err := h.ApplicationService.ListApplications()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get applications"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"apps": apps,
		"total": len(apps),
	})
}

// GetApplication 获取应用详情
func (h *ApplicationHandler) GetApplication(c echo.Context) error {
	id := c.Param("id")
	
	app, err := h.ApplicationService.GetApplicationByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Application not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"app": app,
	})
}

// UpdateApplication 更新应用
func (h *ApplicationHandler) UpdateApplication(c echo.Context) error {
	id := c.Param("id")
	var req UpdateApplicationRequest
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	app, err := h.ApplicationService.UpdateApplication(id, req.Name, req.Code, req.Description, req.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"app":     app,
		"message": "Application updated successfully",
	})
}

// DeleteApplication 删除应用
func (h *ApplicationHandler) DeleteApplication(c echo.Context) error {
	id := c.Param("id")
	
	if err := h.ApplicationService.DeleteApplication(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Application deleted successfully"})
}

// GetApplicationByCode 根据代码获取应用
func (h *ApplicationHandler) GetApplicationByCode(c echo.Context) error {
	code := c.Param("code")
	
	app, err := h.ApplicationService.GetApplicationByCode(code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Application not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"app": app,
	})
}