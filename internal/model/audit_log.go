package model

import (
	"gorm.io/gorm"
)

// AuditLog 审计日志模型
type AuditLog struct {
	gorm.Model
	AppID      uint   `gorm:"index;not null" json:"appId"`      // 所属应用ID, 0表示系统级日志
	UserID     uint   `gorm:"index" json:"userId"`              // 操作人ID
	Username   string `gorm:"size:50" json:"username"`          // 操作人用户名
	Action     string `gorm:"size:50;not null" json:"action"`   // 操作类型: LOGIN, LOGOUT, CREATE, UPDATE, DELETE, ASSIGN
	Resource   string `gorm:"size:50;not null" json:"resource"` // 资源类型: USER, ROLE, MENU, API_PERMISSION, APPLICATION
	ResourceID string `gorm:"size:100" json:"resourceId"`       // 资源ID
	Content    string `gorm:"type:text" json:"content"`         // 操作详细内容 (JSON)
	IP         string `gorm:"size:50" json:"ip"`                // 操作IP
	Status     int    `gorm:"default:1" json:"status"`          // 1=Success, 0=Failed
	ErrorMsg   string `gorm:"type:text" json:"errorMsg"`        // 错误信息
}
