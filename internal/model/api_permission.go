package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApiPermission 接口权限模型
type ApiPermission struct {
	gorm.Model
	UUID        string       `gorm:"uniqueIndex;size:36;not null" json:"uuid"` // 唯一标识, UUID格式
	Name        string       `gorm:"size:100;not null" json:"name"`            // 权限名称
	Path        string       `gorm:"size:200;not null" json:"path"`            // 接口路径
	Method      string       `gorm:"size:10;not null" json:"method"`           // HTTP方法
	Description string       `gorm:"size:255" json:"description"`              // 描述
	AppID       uint         `gorm:"not null" json:"appId"`                    // 所属应用ID
	App         *Application `gorm:"foreignKey:AppID" json:"app,omitempty"`
}

// BeforeCreate GORM钩子，在创建前生成UUID
func (p *ApiPermission) BeforeCreate(tx *gorm.DB) error {
	if p.UUID == "" {
		p.UUID = uuid.New().String()
	}
	return nil
}

// HTTP方法常量
const (
	HTTP_ALL     = "*"
	HTTP_GET     = "GET"
	HTTP_POST    = "POST"
	HTTP_PUT     = "PUT"
	HTTP_DELETE  = "DELETE"
	HTTP_PATCH   = "PATCH"
	HTTP_HEAD    = "HEAD"
	HTTP_OPTIONS = "OPTIONS"
)

// 获取所有HTTP方法
func GetAllHttpMethods() []string {
	return []string{
		HTTP_ALL,
		HTTP_GET,
		HTTP_POST,
		HTTP_PUT,
		HTTP_DELETE,
		HTTP_PATCH,
		HTTP_HEAD,
		HTTP_OPTIONS,
	}
}
