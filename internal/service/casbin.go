package service

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"

	"Authos/internal/model"
)

// CasbinService Casbin 服务
type CasbinService struct {
	Enforcer *casbin.Enforcer
	DB       *gorm.DB
}

// NewCasbinService 创建 Casbin 服务实例
func NewCasbinService(db *gorm.DB) (*CasbinService, error) {
	// 使用 Gorm Adapter 连接到 SQLite 数据库
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	// 优先在根目录查找 model.conf
	modelPath := "model.conf"
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer:%s %w", modelPath, err)
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load casbin policy: %w", err)
	}

	return &CasbinService{
		Enforcer: enforcer,
		DB:       db,
	}, nil
}

// CheckPermission 检查权限
type CheckPermissionReq struct {
	UserID uint   `json:"userId" binding:"required"`
	Obj    string `json:"obj" binding:"required"`
	Act    string `json:"act" binding:"required"`
}

// CheckPermission 检查用户是否具有指定资源的操作权限。
// 若用户持有超级管理员角色则直接放行，否则交由 Casbin 策略判断。
func (s *CasbinService) CheckPermission(userId uint, obj, act string) (bool, error) {
	var user model.User
	if err := s.DB.Preload("Roles").First(&user, userId).Error; err != nil {
		return false, err
	}

	for _, role := range user.Roles {
		// 超级管理员角色直接放行，无需经过 Casbin 策略
		if role.IsSuperAdmin {
			return true, nil
		}

		roleKey := fmt.Sprintf("role:%s", role.UUID)
		allowed, err := s.Enforcer.Enforce(roleKey, obj, act)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}

	return false, nil
}

// LoadPolicy 重新加载策略
func (s *CasbinService) LoadPolicy() error {
	return s.Enforcer.LoadPolicy()
}

// AddPolicy 添加策略
func (s *CasbinService) AddPolicy(sub, obj, act string) error {
	_, err := s.Enforcer.AddPolicy(sub, obj, act)
	return err
}

// RemovePolicy 删除策略
func (s *CasbinService) RemovePolicy(sub, obj, act string) error {
	_, err := s.Enforcer.RemovePolicy(sub, obj, act)
	return err
}

// RemoveFilteredPolicy 根据过滤条件删除策略
func (s *CasbinService) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) error {
	_, err := s.Enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
	return err
}
