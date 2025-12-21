package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	gorm.Model
	UUID           string       `gorm:"uniqueIndex;size:36;not null" json:"uuid"` // 唯一标识, UUID格式
	Name           string       `gorm:"size:50;not null" json:"name"`
	AppID          uint         `gorm:"not null" json:"appId"` // 所属应用ID
	Users          []*User      `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE" json:"users,omitempty"`
	Menus          []*Menu      `gorm:"many2many:role_menus;constraint:OnDelete:CASCADE" json:"menus,omitempty"`
	App            *Application `gorm:"foreignKey:AppID" json:"app,omitempty"`
	MenuCount      int          `gorm:"-" json:"menuCount"`      // 用于列表展示
	ApiPermCount   int          `gorm:"-" json:"apiPermCount"`   // 用于列表展示
	MenuPreview    []string     `gorm:"-" json:"menuPreview"`    // 菜单预览（前几个名称）
	ApiPermPreview []string     `gorm:"-" json:"apiPermPreview"` // 接口预览（前几个名称）
}

// BeforeCreate GORM钩子，在创建前生成UUID
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.UUID == "" {
		r.UUID = uuid.New().String()
	}
	return nil
}
