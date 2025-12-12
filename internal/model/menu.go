package model

import (
	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	gorm.Model
	ParentID  uint     `gorm:"default:0" json:"parentId"` // 树形结构
	Name      string   `gorm:"size:50;not null" json:"name"`
	Path      string   `gorm:"size:200" json:"path"`
	Component string   `gorm:"size:200" json:"component"`
	Type      int      `gorm:"default:0" json:"type"` // 0=Directory, 1=Menu, 2=Button
	Sort      int      `gorm:"default:0" json:"sort"`
	Hidden    bool     `gorm:"default:false" json:"hidden"`
	Roles     []*Role  `gorm:"many2many:role_menus;constraint:OnDelete:CASCADE" json:"roles,omitempty"`
	Children  []*Menu  `gorm:"-" json:"children,omitempty"` // 用于树形返回，不存储到数据库
}
