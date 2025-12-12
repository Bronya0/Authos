package service

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"Authos/internal/model"
)

// ExternalAppMenuService 外部应用菜单服务
type ExternalAppMenuService struct {
	DB *gorm.DB
}

// NewExternalAppMenuService 创建外部应用菜单服务实例
func NewExternalAppMenuService(db *gorm.DB) *ExternalAppMenuService {
	return &ExternalAppMenuService{
		DB: db,
	}
}

// CreateMenu 创建菜单
func (s *ExternalAppMenuService) CreateMenu(menu *model.ExternalAppMenu) error {
	return s.DB.Create(menu).Error
}

// UpdateMenu 更新菜单
func (s *ExternalAppMenuService) UpdateMenu(menu *model.ExternalAppMenu) error {
	return s.DB.Updates(menu).Error
}

// DeleteMenu 删除菜单
func (s *ExternalAppMenuService) DeleteMenu(id uint) error {
	return s.DB.Delete(&model.ExternalAppMenu{}, id).Error
}

// GetMenuByID 根据ID获取菜单
func (s *ExternalAppMenuService) GetMenuByID(id uint) (*model.ExternalAppMenu, error) {
	var menu model.ExternalAppMenu
	if err := s.DB.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetMenusByAppID 根据应用ID获取菜单列表
func (s *ExternalAppMenuService) GetMenusByAppID(appID uint) ([]*model.ExternalAppMenu, error) {
	var menus []*model.ExternalAppMenu
	if err := s.DB.Where("app_id = ?", appID).Order("sort asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetMenuTree 获取菜单树
func (s *ExternalAppMenuService) GetMenuTree(appID uint) ([]*model.ExternalAppMenu, error) {
	menus, err := s.GetMenusByAppID(appID)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	menuMap := make(map[uint]*model.ExternalAppMenu)
	var rootMenus []*model.ExternalAppMenu

	// 创建菜单映射
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	// 构建树形结构
	for _, menu := range menus {
		if menu.ParentID == 0 {
			rootMenus = append(rootMenus, menu)
		} else if _, exists := menuMap[menu.ParentID]; exists {
			// 这里需要处理子菜单，但由于模型中没有直接关联，我们需要在返回结果中构建树
		}
	}

	return rootMenus, nil
}

// ExternalAppRoleService 外部应用角色服务
type ExternalAppRoleService struct {
	DB *gorm.DB
}

// NewExternalAppRoleService 创建外部应用角色服务实例
func NewExternalAppRoleService(db *gorm.DB) *ExternalAppRoleService {
	return &ExternalAppRoleService{
		DB: db,
	}
}

// CreateRole 创建角色
func (s *ExternalAppRoleService) CreateRole(role *model.ExternalAppRole) error {
	return s.DB.Create(role).Error
}

// UpdateRole 更新角色
func (s *ExternalAppRoleService) UpdateRole(role *model.ExternalAppRole) error {
	return s.DB.Updates(role).Error
}

// DeleteRole 删除角色
func (s *ExternalAppRoleService) DeleteRole(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除角色
		if err := tx.Delete(&model.ExternalAppRole{}, id).Error; err != nil {
			return err
		}

		// 删除角色菜单关联
		if err := tx.Where("role_id = ?", id).Delete(&model.ExternalAppRoleMenu{}).Error; err != nil {
			return err
		}

		// 删除角色权限关联
		if err := tx.Where("role_id = ?", id).Delete(&model.ExternalAppRolePermission{}).Error; err != nil {
			return err
		}

		// 删除用户角色关联
		if err := tx.Where("role_id = ?", id).Delete(&model.ExternalAppUserRole{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetRoleByID 根据ID获取角色
func (s *ExternalAppRoleService) GetRoleByID(id uint) (*model.ExternalAppRole, error) {
	var role model.ExternalAppRole
	if err := s.DB.Preload("Menus").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRolesByAppID 根据应用ID获取角色列表
func (s *ExternalAppRoleService) GetRolesByAppID(appID uint) ([]*model.ExternalAppRole, error) {
	var roles []*model.ExternalAppRole
	if err := s.DB.Where("app_id = ?", appID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// AssignMenus 为角色分配菜单
func (s *ExternalAppRoleService) AssignMenus(roleID uint, menuIDs []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有菜单关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.ExternalAppRoleMenu{}).Error; err != nil {
			return err
		}

		// 添加新的菜单关联
		for _, menuID := range menuIDs {
			roleMenu := model.ExternalAppRoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			}
			if err := tx.Create(&roleMenu).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ExternalAppUserService 外部应用用户服务
type ExternalAppUserService struct {
	DB *gorm.DB
}

// NewExternalAppUserService 创建外部应用用户服务实例
func NewExternalAppUserService(db *gorm.DB) *ExternalAppUserService {
	return &ExternalAppUserService{
		DB: db,
	}
}

// CreateUser 创建用户
func (s *ExternalAppUserService) CreateUser(user *model.ExternalAppUser) error {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.DB.Create(user).Error
}

// UpdateUser 更新用户
func (s *ExternalAppUserService) UpdateUser(user *model.ExternalAppUser) error {
	// 如果密码不为空，则加密密码
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return s.DB.Updates(user).Error
}

// DeleteUser 删除用户
func (s *ExternalAppUserService) DeleteUser(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除用户
		if err := tx.Delete(&model.ExternalAppUser{}, id).Error; err != nil {
			return err
		}

		// 删除用户角色关联
		if err := tx.Where("user_id = ?", id).Delete(&model.ExternalAppUserRole{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetUserByID 根据ID获取用户
func (s *ExternalAppUserService) GetUserByID(id uint) (*model.ExternalAppUser, error) {
	var user model.ExternalAppUser
	if err := s.DB.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersByAppID 根据应用ID获取用户列表
func (s *ExternalAppUserService) GetUsersByAppID(appID uint) ([]*model.ExternalAppUser, error) {
	var users []*model.ExternalAppUser
	if err := s.DB.Where("app_id = ?", appID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *ExternalAppUserService) GetUserByUsername(appID uint, username string) (*model.ExternalAppUser, error) {
	var user model.ExternalAppUser
	if err := s.DB.Where("app_id = ? AND username = ?", appID, username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// AssignRoles 为用户分配角色
func (s *ExternalAppUserService) AssignRoles(userID uint, roleIDs []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.ExternalAppUserRole{}).Error; err != nil {
			return err
		}

		// 添加新的角色关联
		for _, roleID := range roleIDs {
			userRole := model.ExternalAppUserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// CreateUserWithRoles 创建用户并分配角色
func (s *ExternalAppUserService) CreateUserWithRoles(user *model.ExternalAppUser, roleIDs []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)

		// 创建用户
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 分配角色
		if len(roleIDs) > 0 {
			for _, roleID := range roleIDs {
				userRole := model.ExternalAppUserRole{
					UserID: user.ID,
					RoleID: roleID,
				}
				if err := tx.Create(&userRole).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// UpdateUserWithRoles 更新用户并更新角色
func (s *ExternalAppUserService) UpdateUserWithRoles(user *model.ExternalAppUser, roleIDs []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 如果密码不为空，则加密密码
		if user.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			user.Password = string(hashedPassword)
		} else {
			// 如果密码为空，则不更新密码字段
			user.Password = ""
		}

		// 更新用户信息
		updateData := map[string]interface{}{
			"username": user.Username,
			"app_id":   user.AppID,
		}

		if user.Password != "" {
			updateData["password"] = user.Password
		}

		if user.Status != 0 {
			updateData["status"] = user.Status
		}

		if err := tx.Model(user).Updates(updateData).Error; err != nil {
			return err
		}

		// 删除原有角色关联
		if err := tx.Where("user_id = ?", user.ID).Delete(&model.ExternalAppUserRole{}).Error; err != nil {
			return err
		}

		// 添加新的角色关联
		if len(roleIDs) > 0 {
			for _, roleID := range roleIDs {
				userRole := model.ExternalAppUserRole{
					UserID: user.ID,
					RoleID: roleID,
				}
				if err := tx.Create(&userRole).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// CheckPassword 检查密码
func (s *ExternalAppUserService) CheckPassword(user *model.ExternalAppUser, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// ExternalAppPermissionService 外部应用权限服务
type ExternalAppPermissionService struct {
	DB *gorm.DB
}

// NewExternalAppPermissionService 创建外部应用权限服务实例
func NewExternalAppPermissionService(db *gorm.DB) *ExternalAppPermissionService {
	return &ExternalAppPermissionService{
		DB: db,
	}
}

// CreatePermission 创建权限
func (s *ExternalAppPermissionService) CreatePermission(permission *model.ExternalAppPermission) error {
	return s.DB.Create(permission).Error
}

// UpdatePermission 更新权限
func (s *ExternalAppPermissionService) UpdatePermission(permission *model.ExternalAppPermission) error {
	return s.DB.Updates(permission).Error
}

// DeletePermission 删除权限
func (s *ExternalAppPermissionService) DeletePermission(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除权限
		if err := tx.Delete(&model.ExternalAppPermission{}, id).Error; err != nil {
			return err
		}

		// 删除角色权限关联
		if err := tx.Where("permission_id = ?", id).Delete(&model.ExternalAppRolePermission{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetPermissionByID 根据ID获取权限
func (s *ExternalAppPermissionService) GetPermissionByID(id uint) (*model.ExternalAppPermission, error) {
	var permission model.ExternalAppPermission
	if err := s.DB.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetPermissionsByAppID 根据应用ID获取权限列表
func (s *ExternalAppPermissionService) GetPermissionsByAppID(appID uint) ([]*model.ExternalAppPermission, error) {
	var permissions []*model.ExternalAppPermission
	if err := s.DB.Where("app_id = ?", appID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// AssignPermissions 为角色分配权限
func (s *ExternalAppPermissionService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.ExternalAppRolePermission{}).Error; err != nil {
			return err
		}

		// 添加新的权限关联
		for _, permissionID := range permissionIDs {
			rolePermission := model.ExternalAppRolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
			}
			if err := tx.Create(&rolePermission).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// CheckUserPermission 检查用户权限
func (s *ExternalAppPermissionService) CheckUserPermission(userID uint, resource string, action string) (bool, error) {
	// 获取用户的所有角色
	var user model.ExternalAppUser
	if err := s.DB.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		return false, err
	}

	// 检查用户的所有角色是否有指定权限
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			if permission.Resource == resource && permission.Action == action {
				return true, nil
			}
		}
	}

	return false, nil
}

// CheckUserMenuPermission 检查用户菜单权限
func (s *ExternalAppPermissionService) CheckUserMenuPermission(userID uint, menuID uint) (bool, error) {
	// 获取用户的所有角色
	var user model.ExternalAppUser
	if err := s.DB.Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return false, err
	}

	// 检查用户的所有角色是否有指定菜单权限
	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			if menu.ID == menuID {
				return true, nil
			}
		}
	}

	return false, nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)[:length]
}

// GenerateAppKey 生成应用密钥
func GenerateAppKey() string {
	return GenerateRandomString(16)
}

// GenerateAppSecret 生成应用密钥
func GenerateAppSecret() string {
	return GenerateRandomString(32)
}

// GenerateToken 生成令牌
func GenerateToken() string {
	return GenerateRandomString(64)
}
