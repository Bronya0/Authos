package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"Authos/internal/model"
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

	h.ApplicationService.DB.Create(&model.AuditLog{
		AppID:      0, // 系统级操作
		UserID:     userID,
		Username:   username,
		Action:     "CREATE",
		Resource:   "APPLICATION",
		ResourceID: fmt.Sprintf("%d", app.ID),
		Content:    fmt.Sprintf("创建应用: %s (%s)", app.Name, app.Code),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"app":       app,
		"message":   "Application created successfully",
		"appId":     app.ID,        // 返回应用ID
		"appUuid":   app.UUID,      // 返回应用UUID
		"secretKey": app.SecretKey, // 返回应用密钥
	})
}

// ListApplications 列出所有应用
func (h *ApplicationHandler) ListApplications(c echo.Context) error {
	name := c.QueryParam("name")
	code := c.QueryParam("code")

	var apps []*model.Application
	db := h.ApplicationService.DB
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		db = db.Where("code LIKE ?", "%"+code+"%")
	}

	if err := db.Order("id asc").Find(&apps).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取应用列表失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"apps":  apps,
		"total": len(apps),
	})
}

// GetApplication 获取应用详情
func (h *ApplicationHandler) GetApplication(c echo.Context) error {
	id := c.Param("id")

	app, err := h.ApplicationService.GetApplicationByIDWithoutSecret(id)
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

	h.ApplicationService.DB.Create(&model.AuditLog{
		AppID:      0,
		UserID:     userID,
		Username:   username,
		Action:     "UPDATE",
		Resource:   "APPLICATION",
		ResourceID: id,
		Content:    fmt.Sprintf("更新应用: %s", req.Name),
		IP:         c.RealIP(),
		Status:     1,
	})

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

	h.ApplicationService.DB.Create(&model.AuditLog{
		AppID:      0,
		UserID:     userID,
		Username:   username,
		Action:     "DELETE",
		Resource:   "APPLICATION",
		ResourceID: id,
		Content:    fmt.Sprintf("删除应用ID: %s", id),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Application deleted successfully"})
}

// GetApplicationByCode 根据代码获取应用
func (h *ApplicationHandler) GetApplicationByCode(c echo.Context) error {
	code := c.Param("code")

	app, err := h.ApplicationService.GetApplicationByCodeWithoutSecret(code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Application not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"app": app,
	})
}
