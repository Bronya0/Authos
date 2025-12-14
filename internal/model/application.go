package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Application 应用模型（租户）
type Application struct {
	gorm.Model
	UUID      string `gorm:"uniqueIndex;size:36;not null" json:"uuid"` // 唯一标识, UUID格式
	Name      string `gorm:"size:100;not null" json:"name"`           // 应用名称
	Code      string `gorm:"uniqueIndex;size:50;not null" json:"code"` // 应用代码
	SecretKey string `gorm:"size:100;not null" json:"-"`               // 应用密钥（不返回给前端）
	Status    int    `gorm:"default:1" json:"status"`                  // 1=Enable, 0=Disable
	Description string `gorm:"size:255" json:"description"`            // 描述

	// 关联关系
	Users []*User `gorm:"foreignKey:AppID" json:"users,omitempty"`
}

// BeforeCreate GORM钩子，在创建前生成UUID
func (a *Application) BeforeCreate(tx *gorm.DB) error {
	if a.UUID == "" {
		a.UUID = uuid.New().String()
	}
	return nil
}