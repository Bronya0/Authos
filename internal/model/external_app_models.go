package model

import (
	"time"
)

// ExternalAppMenu 外部应用菜单模型
type ExternalAppMenu struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	AppID     uint       `json:"app_id" gorm:"not null"`      // 应用ID
	ParentID  uint       `json:"parent_id"`                   // 父菜单ID
	Name      string     `json:"name" gorm:"not null"`        // 菜单名称
	Path      string     `json:"path"`                        // 菜单路径
	Component string     `json:"component"`                   // 组件路径
	Type      int        `json:"type" gorm:"default:1"`       // 菜单类型 0:目录 1:菜单 2:按钮
	Sort      int        `json:"sort" gorm:"default:0"`       // 排序
	Icon      string     `json:"icon"`                        // 图标
	Hidden    bool       `json:"hidden" gorm:"default:false"` // 是否隐藏
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// ExternalAppRole 外部应用角色模型
type ExternalAppRole struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	AppID       uint       `json:"app_id" gorm:"not null"`  // 应用ID
	Name        string     `json:"name" gorm:"not null"`    // 角色名称
	Code        string     `json:"code" gorm:"not null"`    // 角色编码
	Description string     `json:"description"`             // 角色描述
	Status      int        `json:"status" gorm:"default:1"` // 状态 1:启用 0:禁用
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Menus       []ExternalAppMenu       `json:"menus" gorm:"many2many:external_app_role_menus;"`
	Permissions []ExternalAppPermission `json:"permissions" gorm:"many2many:external_app_role_permissions;"`
}

// ExternalAppRoleMenu 外部应用角色菜单关联表
type ExternalAppRoleMenu struct {
	RoleID uint `json:"role_id" gorm:"primaryKey"`
	MenuID uint `json:"menu_id" gorm:"primaryKey"`
}

// ExternalAppUser 外部应用用户模型
type ExternalAppUser struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	AppID     uint       `json:"app_id" gorm:"not null"`   // 应用ID
	Username  string     `json:"username" gorm:"not null"` // 用户名
	Password  string     `json:"-" gorm:"not null"`        // 密码
	RealName  string     `json:"real_name"`                // 真实姓名
	Email     string     `json:"email"`                    // 邮箱
	Phone     string     `json:"phone"`                    // 手机号
	Status    int        `json:"status" gorm:"default:1"`  // 状态 1:启用 0:禁用
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Roles []ExternalAppRole `json:"roles" gorm:"many2many:external_app_user_roles;"`
}

// ExternalAppUserRole 外部应用用户角色关联表
type ExternalAppUserRole struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`
	RoleID uint `json:"role_id" gorm:"primaryKey"`
}

// ExternalAppPermission 外部应用权限模型
type ExternalAppPermission struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	AppID       uint       `json:"app_id" gorm:"not null"`   // 应用ID
	Name        string     `json:"name" gorm:"not null"`     // 权限名称
	Code        string     `json:"code" gorm:"not null"`     // 权限编码
	Type        int        `json:"type" gorm:"default:2"`    // 权限类型 1:菜单 2:接口
	Resource    string     `json:"resource" gorm:"not null"` // 资源
	Action      string     `json:"action" gorm:"not null"`   // 操作
	Description string     `json:"description"`              // 权限描述
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// ExternalAppRolePermission 外部应用角色权限关联表
type ExternalAppRolePermission struct {
	RoleID       uint `json:"role_id" gorm:"primaryKey"`
	PermissionID uint `json:"permission_id" gorm:"primaryKey"`
}
