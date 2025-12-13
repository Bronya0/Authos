package service

import (
	"fmt"

	"gorm.io/gorm"

	"Authos/internal/model"
)

// RoleService 角色服务
type RoleService struct {
	DB            *gorm.DB
	CasbinService *CasbinService
}

// NewRoleService 创建角色服务实例
func NewRoleService(db *gorm.DB, casbinService *CasbinService) *RoleService {
	return &RoleService{
		DB:            db,
		CasbinService: casbinService,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(role *model.Role) error {
	// 检查角色代码是否已存在
	var count int64
	if err := s.DB.Model(&model.Role{}).Where("code = ?", role.Code).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	
	if count > 0 {
		return fmt.Errorf("role with code '%s' already exists", role.Code)
	}

	if err := s.DB.Create(role).Error; err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(role *model.Role) error {
	return s.DB.Updates(role).Error
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	// 先删除相关的 Casbin 策略（在事务外）
	roleKey := fmt.Sprintf("role:%d", id)
	
	// 移除用户-角色关联
	if err := s.CasbinService.RemoveFilteredPolicy(1, roleKey); err != nil {
		return fmt.Errorf("failed to remove user-role policies for role %d: %w", id, err)
	}
	// 移除角色-权限关联
	if err := s.CasbinService.RemoveFilteredPolicy(0, roleKey); err != nil {
		return fmt.Errorf("failed to remove role-permission policies for role %d: %w", id, err)
	}

	// 重新加载策略
	if err := s.CasbinService.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to reload Casbin policies after deleting role %d: %w", id, err)
	}

	// 在事务中删除角色
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除角色
		if err := tx.Delete(&model.Role{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete role with ID %d: %w", id, err)
		}
		return nil
	})
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := s.DB.Preload("Menus").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByCode 根据Code获取角色
func (s *RoleService) GetRoleByCode(code string) (*model.Role, error) {
	var role model.Role
	if err := s.DB.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRoles 列出所有角色
func (s *RoleService) ListRoles() ([]*model.Role, error) {
	var roles []*model.Role
	if err := s.DB.Order("sort asc").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// AssignMenus 为角色分配菜单
func (s *RoleService) AssignMenus(roleID uint, menuIDs []uint) error {
	// 开始事务
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 获取角色
		var role model.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		// 获取菜单
		var menus []*model.Menu
		if err := tx.Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
			return err
		}

		// 替换角色菜单关联
		if err := tx.Model(&role).Association("Menus").Replace(menus); err != nil {
			return err
		}

		return nil
	})
}

// AssignPermissions 为角色分配 API 权限
func (s *RoleService) AssignPermissions(roleID uint, permissions []map[string]string) error {
	// 开始事务
	return s.DB.Transaction(func(tx *gorm.DB) error {
		roleKey := fmt.Sprintf("role:%d", roleID)

		// 先删除该角色的所有权限
		if err := s.CasbinService.RemoveFilteredPolicy(0, roleKey); err != nil {
			return err
		}

		// 添加新权限
		for _, perm := range permissions {
			obj := perm["obj"]
			act := perm["act"]
			if err := s.CasbinService.AddPolicy(roleKey, obj, act); err != nil {
				return err
			}
		}

		// 重新加载策略
		return s.CasbinService.LoadPolicy()
	})
}
