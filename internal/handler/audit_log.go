package handler

import (
	"Authos/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuditLogHandler 审计日志处理器
type AuditLogHandler struct {
	AuditLogService *service.AuditLogService
}

// NewAuditLogHandler 创建审计日志处理器实例
func NewAuditLogHandler(auditLogService *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{AuditLogService: auditLogService}
}

// ListAuditLogs 查询当前应用的审计日志列表，支持按操作类型、资源、用户名过滤
func (h *AuditLogHandler) ListAuditLogs(c echo.Context) error {
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	action := c.QueryParam("action")
	resource := c.QueryParam("resource")
	username := c.QueryParam("username")

	logs, err := h.AuditLogService.ListAuditLogs(appID, action, resource, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch logs"})
	}

	return c.JSON(http.StatusOK, logs)
}

// ListSystemAuditLogs 查询全局（系统级）审计日志列表，仅供系统管理员使用
func (h *AuditLogHandler) ListSystemAuditLogs(c echo.Context) error {
	// 系统管理员权限检查通常由中间件处理
	action := c.QueryParam("action")
	resource := c.QueryParam("resource")
	username := c.QueryParam("username")

	logs, err := h.AuditLogService.ListSystemAuditLogs(action, resource, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch logs"})
	}

	return c.JSON(http.StatusOK, logs)
}
