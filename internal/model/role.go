package model

import (
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	gorm.Model
	Code  string    `gorm:"uniqueIndex;size:50;not null" json:"code"` // 唯一标识, 如 admin
	Name  string    `gorm:"size:50;not null" json:"name"`
	Sort  int       `gorm:"default:0" json:"sort"`
	Users []*User   `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE" json:"users,omitempty"`
	Menus []*Menu   `gorm:"many2many:role_menus;constraint:OnDelete:CASCADE" json:"menus,omitempty"`
}
