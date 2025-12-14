package service

import (
	"fmt"
	"time"

	"Authos/internal/model"

	"gorm.io/gorm"
)

// PermissionService 权限服务
type PermissionService struct {
	DB            *gorm.DB
	CasbinService *CasbinService
	RoleService   *RoleService
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(db *gorm.DB, casbinService *CasbinService, roleService *RoleService) *PermissionService {
	return &PermissionService{
		DB:            db,
		CasbinService: casbinService,
		RoleService:   roleService,
	}
}

// Permission 权限结构体
type Permission struct {
	Obj         string    `json:"obj"`
	Act         string    `json:"act"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

// GetAllPermissions 获取所有权限（按应用隔离）
func (s *PermissionService) GetAllPermissions(appID uint) ([]Permission, error) {
	// 获取应用的所有角色
	roles, err := s.RoleService.ListRolesByApp(appID)
	if err != nil {
		return []Permission{}, fmt.Errorf("failed to get roles for app %d: %w", appID, err)
	}

	// 收集所有角色的权限
	var allPermissions []Permission
	for _, role := range roles {
		roleKey := fmt.Sprintf("role:%s", role.UUID)
		policies, _ := s.CasbinService.Enforcer.GetFilteredPolicy(0, roleKey)

		for _, policy := range policies {
			if len(policy) >= 3 {
				permissions := Permission{
					Obj:         policy[1],
					Act:         policy[2],
					Description: "",         // Casbin不存储描述信息
					CreatedAt:   time.Now(), // Casbin不记录创建时间，使用当前时间
				}

				// 检查权限是否已经存在（避免重复）
				exists := false
				for _, p := range allPermissions {
					if p.Obj == permissions.Obj && p.Act == permissions.Act {
						exists = true
						break
					}
				}
				if !exists {
					allPermissions = append(allPermissions, permissions)
				}
			}
		}
	}

	return allPermissions, nil
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(obj, act, description string) error {
	// 验证权限格式
	if obj == "" || act == "" {
		return fmt.Errorf("权限对象和动作不能为空")
	}

	// 检查权限是否已存在
	hasPolicy, _ := s.CasbinService.Enforcer.HasPolicy("role:1", obj, act)
	if hasPolicy {
		return fmt.Errorf("权限已存在")
	}

	// 创建一个默认角色来存储权限
	// 在实际应用中，权限是直接分配给角色的，而不是独立存在
	// 这里我们使用一个默认角色ID 1来存储权限
	_, err := s.CasbinService.Enforcer.AddPolicy("role:1", obj, act)
	if err != nil {
		return fmt.Errorf("添加权限策略失败: %v", err)
	}

	// 重新加载策略
	return s.CasbinService.Enforcer.LoadPolicy()
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(obj, act string) error {
	// 从所有角色中删除指定的权限
	policies, _ := s.CasbinService.Enforcer.GetPolicy()

	for _, policy := range policies {
		if len(policy) >= 3 && policy[1] == obj && policy[2] == act {
			if len(policy[0]) > 5 && policy[0][:5] == "role:" {
				// 删除角色权限
				_, err := s.CasbinService.Enforcer.RemovePolicy(policy[0], policy[1], policy[2])
				if err != nil {
					return err
				}
			}
		}
	}

	// 重新加载策略
	return s.CasbinService.Enforcer.LoadPolicy()
}

// GetRolesForPermission 获取拥有指定权限的角色
func (s *PermissionService) GetRolesForPermission(obj, act string) ([]model.Role, error) {
	// 获取所有拥有此权限的策略
	policies, _ := s.CasbinService.Enforcer.GetFilteredPolicy(1, obj, act)

	var roleUUIDs []string
	for _, policy := range policies {
		if len(policy) >= 1 && len(policy[0]) > 5 && policy[0][:5] == "role:" {
			// 提取角色UUID
			roleUUID := policy[0][5:]
			roleUUIDs = append(roleUUIDs, roleUUID)
		}
	}

	// 查询角色信息
	var roles []model.Role
	if len(roleUUIDs) > 0 {
		if err := s.DB.Where("uuid IN ?", roleUUIDs).Find(&roles).Error; err != nil {
			return nil, err
		}
	}

	return roles, nil
}
