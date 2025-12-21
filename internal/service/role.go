package service

import (
	"fmt"
	"strings"
	"unicode"

	"gorm.io/gorm"

	"Authos/internal/model"
)

// RoleService 角色服务
type RoleService struct {
	DB            *gorm.DB
	CasbinService *CasbinService
}

// generateRoleCode 根据角色名称生成角色编码
func generateRoleCode(name string) string {
	// 将中文字符转换为拼音或直接使用英文字符
	// 这里简单地将非字母数字字符替换为空格，然后将空格替换为下划线
	var result strings.Builder

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) {
			result.WriteRune('_')
		}
		// 忽略其他字符
	}

	code := result.String()

	// 如果生成的编码为空，使用默认值
	if code == "" {
		code = "role"
	}

	return code
}

// NewRoleService 创建角色服务实例
func NewRoleService(db *gorm.DB, casbinService *CasbinService) *RoleService {
	return &RoleService{
		DB:            db,
		CasbinService: casbinService,
	}
}

// CreateRole 创建角色（按应用隔离）
func (s *RoleService) CreateRole(role *model.Role) error {
	// UUID会在BeforeCreate钩子中自动生成，无需手动处理

	if err := s.DB.Create(role).Error; err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// UpdateRole 更新角色（按应用隔离）
func (s *RoleService) UpdateRole(role *model.Role) error {
	return s.DB.Where("id = ? AND app_id = ?", role.ID, role.AppID).Updates(role).Error
}

// DeleteRole 删除角色（按应用隔离）
func (s *RoleService) DeleteRole(id uint, appID uint) error {
	// 先获取角色UUID
	role, err := s.GetRoleByID(id, appID)
	if err != nil {
		return fmt.Errorf("failed to get role with ID %d: %w", id, err)
	}

	// 先删除相关的 Casbin 策略（在事务外）
	roleKey := fmt.Sprintf("role:%s", role.UUID)

	// 移除用户-角色关联
	if err := s.CasbinService.RemoveFilteredPolicy(1, roleKey); err != nil {
		return fmt.Errorf("failed to remove user-role policies for role %s: %w", role.UUID, err)
	}
	// 移除角色-权限关联
	if err := s.CasbinService.RemoveFilteredPolicy(0, roleKey); err != nil {
		return fmt.Errorf("failed to remove role-permission policies for role %s: %w", role.UUID, err)
	}

	// 重新加载策略
	if err := s.CasbinService.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to reload Casbin policies after deleting role %s: %w", role.UUID, err)
	}

	// 在事务中删除角色
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除角色
		if err := tx.Where("id = ? AND app_id = ?", id, appID).Delete(&model.Role{}).Error; err != nil {
			return fmt.Errorf("failed to delete role with ID %d: %w", id, err)
		}
		return nil
	})
}

// GetRoleByID 根据ID获取角色（按应用隔离）
func (s *RoleService) GetRoleByID(id uint, appID uint) (*model.Role, error) {
	var role model.Role
	if err := s.DB.Preload("Menus").Where("id = ? AND app_id = ?", id, appID).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByUUID 根据UUID获取角色
func (s *RoleService) GetRoleByUUID(uuid string) (*model.Role, error) {
	var role model.Role
	if err := s.DB.Where("uuid = ?", uuid).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRolesByApp 列出指定应用的所有角色
func (s *RoleService) ListRolesByApp(appID uint) ([]*model.Role, error) {
	var roles []*model.Role
	// 预加载菜单以统计数量和预览
	if err := s.DB.Preload("Menus").Where("app_id = ?", appID).Order("id asc").Find(&roles).Error; err != nil {
		return nil, err
	}

	for _, role := range roles {
		// 菜单信息
		role.MenuCount = len(role.Menus)
		previewCount := 3
		if role.MenuCount < previewCount {
			previewCount = role.MenuCount
		}
		role.MenuPreview = make([]string, 0, previewCount)
		for i := 0; i < previewCount; i++ {
			role.MenuPreview = append(role.MenuPreview, role.Menus[i].Name)
		}

		// API 权限信息 (通过 Casbin)
		roleKey := fmt.Sprintf("role:%s", role.UUID)
		policies, _ := s.CasbinService.Enforcer.GetFilteredPolicy(0, roleKey)
		role.ApiPermCount = len(policies)

		apiPreviewCount := 3
		if role.ApiPermCount < apiPreviewCount {
			apiPreviewCount = role.ApiPermCount
		}
		role.ApiPermPreview = make([]string, 0, apiPreviewCount)
		for i := 0; i < apiPreviewCount; i++ {
			// policies[i] 为 [sub, obj, act]
			role.ApiPermPreview = append(role.ApiPermPreview, fmt.Sprintf("%s %s", policies[i][2], policies[i][1]))
		}
	}

	return roles, nil
}

// AssignMenus 为角色分配菜单（按应用隔离）
func (s *RoleService) AssignMenus(roleID uint, appID uint, menuIDs []uint) error {
	// 开始事务
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 获取角色
		var role model.Role
		if err := tx.Where("id = ? AND app_id = ?", roleID, appID).First(&role).Error; err != nil {
			return err
		}

		// 获取菜单（确保属于同一应用）
		var menus []*model.Menu
		if err := tx.Where("id IN ? AND app_id = ?", menuIDs, appID).Find(&menus).Error; err != nil {
			return err
		}

		// 替换角色菜单关联
		if err := tx.Model(&role).Association("Menus").Replace(menus); err != nil {
			return err
		}

		return nil
	})
}

// GetRoleMenus 获取角色的菜单权限（按应用隔离）
func (s *RoleService) GetRoleMenus(roleID uint, appID uint) ([]*model.Menu, error) {
	var role model.Role
	if err := s.DB.Preload("Menus").Where("id = ? AND app_id = ?", roleID, appID).First(&role).Error; err != nil {
		return nil, err
	}
	return role.Menus, nil
}

// AssignPermissions 为角色分配 API 权限（按应用隔离）
func (s *RoleService) AssignPermissions(roleID uint, appID uint, permissions []map[string]string) error {
	// 先获取角色UUID
	role, err := s.GetRoleByID(roleID, appID)
	if err != nil {
		return fmt.Errorf("failed to get role with ID %d: %w", roleID, err)
	}

	// 开始事务
	return s.DB.Transaction(func(tx *gorm.DB) error {
		roleKey := fmt.Sprintf("role:%s", role.UUID)

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
