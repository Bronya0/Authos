package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"Authos/internal/model"
)

// UserService 用户服务
type UserService struct {
	DB *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *model.User) error {
	// 检查用户名是否已存在
	var count int64
	if err := s.DB.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	
	if count > 0 {
		return fmt.Errorf("user with username '%s' already exists", user.Username)
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// 开始事务
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// 关联角色
		if len(user.RoleIDs) > 0 {
			var roles []*model.Role
			if err := tx.Where("id IN ?", user.RoleIDs).Find(&roles).Error; err != nil {
				return fmt.Errorf("failed to find roles: %w", err)
			}
			if err := tx.Model(user).Association("Roles").Replace(roles); err != nil {
				return fmt.Errorf("failed to associate roles: %w", err)
			}
		}

		return nil
	})
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *model.User) error {
	// 开始事务
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 更新用户基本信息（不包含密码）
		updateData := map[string]interface{}{
			"Username": user.Username,
			"Status":   user.Status,
		}
		if err := tx.Model(user).Updates(updateData).Error; err != nil {
			return err
		}

		// 关联角色
		if len(user.RoleIDs) > 0 {
			var roles []*model.Role
			if err := tx.Where("id IN ?", user.RoleIDs).Find(&roles).Error; err != nil {
				return err
			}
			if err := tx.Model(user).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateUserPassword 更新用户密码
func (s *UserService) UpdateUserPassword(id uint, password string) error {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.DB.Model(&model.User{}).Where("id = ?", id).Update("password", string(hashedPassword)).Error
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	return s.DB.Delete(&model.User{}, id).Error
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.DB.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}

	// 填充 RoleIDs
	for _, role := range user.Roles {
		user.RoleIDs = append(user.RoleIDs, role.ID)
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.DB.Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers 列出所有用户
func (s *UserService) ListUsers() ([]*model.User, error) {
	var users []*model.User
	if err := s.DB.Preload("Roles").Find(&users).Error; err != nil {
		return nil, err
	}

	// 填充 RoleIDs
	for _, user := range users {
		for _, role := range user.Roles {
			user.RoleIDs = append(user.RoleIDs, role.ID)
		}
	}

	return users, nil
}
