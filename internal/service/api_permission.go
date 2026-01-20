package service

import (
	"fmt"
	"strings"

	"Authos/internal/model"

	"gorm.io/gorm"
)

// ApiPermissionService 接口权限服务
type ApiPermissionService struct {
	DB            *gorm.DB
	CasbinService *CasbinService
	RoleService   *RoleService
}

// NewApiPermissionService 创建接口权限服务实例
func NewApiPermissionService(db *gorm.DB, casbinService *CasbinService, roleService *RoleService) *ApiPermissionService {
	return &ApiPermissionService{
		DB:            db,
		CasbinService: casbinService,
		RoleService:   roleService,
	}
}

// GetAllApiPermissions 获取所有接口权限（按应用隔离）
func (s *ApiPermissionService) GetAllApiPermissions(appID uint) ([]model.ApiPermission, error) {
	var permissions []model.ApiPermission
	if err := s.DB.Where("app_id = ?", appID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetApiPermission 根据ID获取接口权限（按应用隔离）
func (s *ApiPermissionService) GetApiPermission(id uint, appID uint) (*model.ApiPermission, error) {
	var permission model.ApiPermission
	if err := s.DB.Where("id = ? AND app_id = ?", id, appID).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetApiPermissionByUUID 根据UUID获取接口权限（按应用隔离）
func (s *ApiPermissionService) GetApiPermissionByUUID(appID uint, uuid string) (*model.ApiPermission, error) {
	var permission model.ApiPermission
	if err := s.DB.Where("uuid = ? AND app_id = ?", uuid, appID).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetApiPermissionByPathAndMethod 根据路径和方法获取接口权限（支持前缀匹配，优先具体方法，其次 *）（按应用隔离）
func (s *ApiPermissionService) GetApiPermissionByPathAndMethod(appID uint, path, method string) (*model.ApiPermission, error) {
	// 统一规范路径，避免尾部斜杠造成的不一致
	reqPath := strings.TrimRight(path, "/")
	if reqPath == "" {
		reqPath = "/"
	}

	var permissions []model.ApiPermission

	// 一次性查出当前应用下当前方法和 * 方法的所有配置，减少多次 IO
	if err := s.DB.Where("app_id = ? AND method IN ?",
		appID,
		[]string{method, model.HTTP_ALL},
	).Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("查询接口权限失败: %v", err)
	}

	var (
		bestMatch       *model.ApiPermission
		bestPathLength  int
		bestMethodScore int // 2: 精确方法; 1: HTTP_ALL
	)

	for i := range permissions {
		p := &permissions[i]
		cfgPath := strings.TrimRight(p.Path, "/")
		if cfgPath == "" {
			cfgPath = "/"
		}

		// 路径前缀匹配：配置的 path 必须是请求 path 的前缀
		if !strings.HasPrefix(reqPath, cfgPath) {
			continue
		}

		// method 匹配优先级：先精确方法，再 HTTP_ALL（*）
		methodScore := 0
		if p.Method == method {
			methodScore = 2
		} else if p.Method == model.HTTP_ALL {
			methodScore = 1
		}
		if methodScore == 0 {
			continue
		}

		// 选择更长的前缀（更具体的路径），并且在相同前缀长度下优先方法更精确
		if len(cfgPath) > bestPathLength || (len(cfgPath) == bestPathLength && methodScore > bestMethodScore) {
			bestMatch = p
			bestPathLength = len(cfgPath)
			bestMethodScore = methodScore
		}
	}

	if bestMatch == nil {
		return nil, fmt.Errorf("接口权限未找到: path=%s method=%s", path, method)
	}

	return bestMatch, nil
}

// CreateApiPermission 创建接口权限（按应用隔离）
func (s *ApiPermissionService) CreateApiPermission(appID uint, key, name, path, method, description string) (*model.ApiPermission, error) {
	if key == "" || name == "" || path == "" || method == "" {
		return nil, fmt.Errorf("权限标识、名称、接口路径和HTTP方法不能为空")
	}

	// 验证HTTP方法是否有效
	validMethods := model.GetAllHttpMethods()
	methodValid := false
	for _, validMethod := range validMethods {
		if method == validMethod {
			methodValid = true
			break
		}
	}
	if !methodValid {
		return nil, fmt.Errorf("无效的HTTP方法: %s", method)
	}

	// 检查权限标识是否已存在（同一应用内）
	var existingPermission model.ApiPermission
	if err := s.DB.Where("key = ? AND app_id = ?", key, appID).First(&existingPermission).Error; err == nil {
		return nil, fmt.Errorf("权限标识已存在: %s", key)
	}

	// 创建权限
	permission := model.ApiPermission{
		Key:         key,
		Name:        name,
		Path:        path,
		Method:      method,
		Description: description,
		AppID:       appID,
	}

	if err := s.DB.Create(&permission).Error; err != nil {
		return nil, fmt.Errorf("创建接口权限失败: %v", err)
	}

	return &permission, nil
}

// UpdateApiPermission 更新接口权限（按应用隔离）
func (s *ApiPermissionService) UpdateApiPermission(id uint, appID uint, key, name, path, method, description string) (*model.ApiPermission, error) {
	// 获取现有权限
	permission, err := s.GetApiPermission(id, appID)
	if err != nil {
		return nil, err
	}

	oldKey := permission.Key

	// 验证HTTP方法是否有效
	validMethods := model.GetAllHttpMethods()
	methodValid := false
	for _, validMethod := range validMethods {
		if method == validMethod {
			methodValid = true
			break
		}
	}
	if !methodValid {
		return nil, fmt.Errorf("无效的HTTP方法: %s", method)
	}

	// 检查权限标识是否已存在（排除当前权限）
	var existingPermission model.ApiPermission
	if err := s.DB.Where("key = ? AND id != ? AND app_id = ?", key, id, appID).First(&existingPermission).Error; err == nil {
		return nil, fmt.Errorf("权限标识已存在: %s", key)
	}

	// 更新权限
	permission.Key = key
	permission.Name = name
	permission.Path = path
	permission.Method = method
	permission.Description = description

	if err := s.DB.Save(permission).Error; err != nil {
		return nil, fmt.Errorf("更新接口权限失败: %v", err)
	}

	if oldKey != key {
		policies, _ := s.CasbinService.Enforcer.GetFilteredPolicy(1, oldKey)
		for _, policy := range policies {
			args := make([]interface{}, len(policy))
			for i, v := range policy {
				args[i] = v
			}
			_, _ = s.CasbinService.Enforcer.RemovePolicy(args...)
			policy[1] = key
			newArgs := make([]interface{}, len(policy))
			for i, v := range policy {
				newArgs[i] = v
			}
			_, err = s.CasbinService.Enforcer.AddPolicy(newArgs...)
			if err != nil {
				return nil, fmt.Errorf("更新权限策略失败: %v", err)
			}
		}
		if err := s.CasbinService.Enforcer.LoadPolicy(); err != nil {
			return nil, fmt.Errorf("重新加载权限策略失败: %v", err)
		}
	}

	return permission, nil
}

// DeleteApiPermission 删除接口权限（按应用隔离）
func (s *ApiPermissionService) DeleteApiPermission(id uint, appID uint) error {
	// 获取权限
	permission, err := s.GetApiPermission(id, appID)
	if err != nil {
		return err
	}

	// 从Casbin中删除所有与此权限相关的策略
	policies, _ := s.CasbinService.Enforcer.GetFilteredPolicy(1, permission.Key)
	for _, policy := range policies {
		// 将[]string转换为[]interface{}
		args := make([]interface{}, len(policy))
		for i, v := range policy {
			args[i] = v
		}
		s.CasbinService.Enforcer.RemovePolicy(args...)
	}

	// 删除权限记录
	if err := s.DB.Delete(permission).Error; err != nil {
		return fmt.Errorf("删除接口权限失败: %v", err)
	}

	// 重新加载策略
	s.CasbinService.Enforcer.LoadPolicy()

	return nil
}

// GetRolesForApiPermission 获取拥有指定接口权限的角色（按应用隔离）
func (s *ApiPermissionService) GetRolesForApiPermission(appID uint, permissionUUID string) ([]model.Role, error) {
	// 获取权限
	permission, err := s.GetApiPermissionByUUID(appID, permissionUUID)
	if err != nil {
		return nil, err
	}

	// 获取所有拥有此权限的策略
	policies, _ := s.CasbinService.Enforcer.GetFilteredPolicy(1, permission.Key)

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

// AddApiPermissionToRole 为角色添加接口权限（按应用隔离）
func (s *ApiPermissionService) AddApiPermissionToRole(appID uint, roleUUID, permissionUUID string) error {
	// 获取权限
	permission, err := s.GetApiPermissionByUUID(appID, permissionUUID)
	if err != nil {
		return fmt.Errorf("接口权限不存在: %v", err)
	}

	// 检查权限是否已存在
	rolePrefix := fmt.Sprintf("role:%s", roleUUID)
	hasPolicy, _ := s.CasbinService.Enforcer.HasPolicy(rolePrefix, permission.Key, model.HTTP_ALL)
	if hasPolicy {
		return fmt.Errorf("角色已拥有此权限")
	}

	// 添加权限策略
	_, err = s.CasbinService.Enforcer.AddPolicy(rolePrefix, permission.Key, model.HTTP_ALL)
	if err != nil {
		return fmt.Errorf("添加权限策略失败: %v", err)
	}

	// 重新加载策略
	return s.CasbinService.Enforcer.LoadPolicy()
}

// RemoveApiPermissionFromRole 移除角色的接口权限
func (s *ApiPermissionService) RemoveApiPermissionFromRole(appID uint, roleUUID, permissionUUID string) error {
	// 获取权限
	permission, err := s.GetApiPermissionByUUID(appID, permissionUUID)
	if err != nil {
		return fmt.Errorf("接口权限不存在: %v", err)
	}

	// 移除权限策略
	rolePrefix := fmt.Sprintf("role:%s", roleUUID)
	_, err = s.CasbinService.Enforcer.RemovePolicy(rolePrefix, permission.Key, model.HTTP_ALL)
	if err != nil {
		return fmt.Errorf("移除权限策略失败: %v", err)
	}

	// 重新加载策略
	return s.CasbinService.Enforcer.LoadPolicy()
}

// GetApiPermissionsForRole 获取角色的接口权限（按应用隔离）
func (s *ApiPermissionService) GetApiPermissionsForRole(appID uint, roleUUID string) ([]model.ApiPermission, error) {
	// 获取角色的所有权限策略
	rolePrefix := fmt.Sprintf("role:%s", roleUUID)
	policies, _ := s.CasbinService.Enforcer.GetFilteredPolicy(0, rolePrefix)

	var permissionKeys []string

	for _, policy := range policies {
		if len(policy) >= 3 {
			permissionKeys = append(permissionKeys, policy[1])
		}
	}

	// 查询权限信息
	var permissions []model.ApiPermission
	if len(permissionKeys) > 0 {
		if err := s.DB.Where("key IN ?", permissionKeys).Find(&permissions).Error; err != nil {
			return nil, err
		}
	}

	return permissions, nil
}
