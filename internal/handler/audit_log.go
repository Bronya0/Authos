package handler

import (
	"Authos/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuditLogHandler struct {
	AuditLogService *service.AuditLogService
}

func NewAuditLogHandler(auditLogService *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{AuditLogService: auditLogService}
}

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
