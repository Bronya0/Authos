package service

import (
	"Authos/internal/model"

	"gorm.io/gorm"
)

// AuditLogService 审计日志服务
type AuditLogService struct {
	DB *gorm.DB
}

// NewAuditLogService 创建审计日志服务实例
func NewAuditLogService(db *gorm.DB) *AuditLogService {
	return &AuditLogService{DB: db}
}

// Record 记录审计日志
func (s *AuditLogService) Record(log *model.AuditLog) {
	s.DB.Create(log)
}

// ListAuditLogs 列出审计日志
func (s *AuditLogService) ListAuditLogs(appID uint, action, resource, username string) ([]*model.AuditLog, error) {
	var logs []*model.AuditLog
	db := s.DB.Where("app_id = ?", appID)
	if action != "" {
		db = db.Where("action = ?", action)
	}
	if resource != "" {
		db = db.Where("resource = ?", resource)
	}
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if err := db.Order("id desc").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// ListSystemAuditLogs 列出系统级审计日志 (appID = 0)
func (s *AuditLogService) ListSystemAuditLogs(action, resource, username string) ([]*model.AuditLog, error) {
	return s.ListAuditLogs(0, action, resource, username)
}
