package model

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password string    `gorm:"size:100;not null" json:"-"`
	Status   int       `gorm:"default:1" json:"status"` // 1=Enable, 0=Disable
	AppID    uint      `gorm:"not null" json:"appId"`   // 所属应用ID
	Roles    []*Role   `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE" json:"roles,omitempty"`
	RoleIDs  []uint    `gorm:"-" json:"roleIds,omitempty"` // 用于回显，不存储到数据库
	App      *Application `gorm:"foreignKey:AppID" json:"app,omitempty"`
}
