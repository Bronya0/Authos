package model

import (
	"time"
)

// ExternalApp 外部应用模型
type ExternalApp struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`                // 应用名称
	Description string     `json:"description"`                         // 应用描述
	AppKey      string     `json:"app_key" gorm:"uniqueIndex;not null"` // 应用密钥
	AppSecret   string     `json:"app_secret" gorm:"not null"`          // 应用密钥
	Status      int        `json:"status" gorm:"default:1"`             // 状态 1:启用 0:禁用
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// ExternalAppToken 外部应用令牌模型
type ExternalAppToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AppID     uint      `json:"app_id" gorm:"not null"`            // 应用ID
	Token     string    `json:"token" gorm:"uniqueIndex;not null"` // 令牌
	ExpiresAt time.Time `json:"expires_at"`                        // 过期时间
	CreatedAt time.Time `json:"created_at"`
}
