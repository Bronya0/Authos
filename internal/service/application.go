package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"gorm.io/gorm"

	"Authos/internal/model"
)

// ApplicationService 应用服务
type ApplicationService struct {
	DB *gorm.DB
}

// NewApplicationService 创建应用服务实例
func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{DB: db}
}

// generateSecretKey 生成应用密钥
func (s *ApplicationService) generateSecretKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

// CreateApplication 创建应用
func (s *ApplicationService) CreateApplication(name, code, description string) (*model.Application, error) {
	fmt.Printf("CreateApplication called with name=%s, code=%s, description=%s\n", name, code, description)

	// 检查应用代码是否已存在（包括软删除的记录）
	var count int64
	if err := s.DB.Unscoped().Model(&model.Application{}).Where("code = ?", code).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check application existence: %w", err)
	}

	fmt.Printf("Found %d applications with code '%s'\n", count, code)
	if count > 0 {
		return nil, fmt.Errorf("application with code '%s' already exists", code)
	}

	// 创建应用
	app := &model.Application{
		Name:        name,
		Code:        code,
		Description: description,
		Status:      1, // 默认启用
	}

	// 生成密钥和UUID
	app.SecretKey = s.generateSecretKey()

	if err := s.DB.Create(app).Error; err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}
	return app, nil
}

// UpdateApplication 更新应用
func (s *ApplicationService) UpdateApplication(id, name, code, description string, status int) (*model.Application, error) {
	// 转换为uint
	appID := s.parseID(id)

	// 检查代码冲突（排除当前应用）
	var count int64
	if err := s.DB.Model(&model.Application{}).Where("code = ? AND id != ?", code, appID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check application existence: %w", err)
	}

	if count > 0 {
		return nil, fmt.Errorf("application with code '%s' already exists", code)
	}

	// 更新应用
	app := &model.Application{
		Model: gorm.Model{
			ID: appID,
		},
		Name:        name,
		Code:        code,
		Description: description,
		Status:      status,
	}

	if err := s.DB.Updates(app).Error; err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	// 返回更新后的应用
	return s.GetApplicationByID(id)
}

// DeleteApplication 删除应用
func (s *ApplicationService) DeleteApplication(id string) error {
	fmt.Printf("DeleteApplication called with ID: %s\n", id)
	appID := s.parseID(id)
	fmt.Printf("Parsed ID: %d\n", appID)
	if appID == 0 {
		return fmt.Errorf("invalid application ID: %s", id)
	}

	// 先检查应用是否存在
	var app model.Application
	if err := s.DB.First(&app, appID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("application with ID %s not found", id)
		}
		return err
	}

	fmt.Printf("Application found: %+v\n", app)

	// 开启事务进行级联删除
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除关联表数据 (User-Role, Role-Menu)
		// 注意：由于SQLite/GORM可能没有建立严格的级联删除，手动清理更安全
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE app_id = ?)", appID).Error; err != nil {
			return fmt.Errorf("failed to delete user_roles: %w", err)
		}
		if err := tx.Exec("DELETE FROM role_menus WHERE role_id IN (SELECT id FROM roles WHERE app_id = ?)", appID).Error; err != nil {
			return fmt.Errorf("failed to delete role_menus: %w", err)
		}

		// 2. 删除各实体表数据
		// 使用 Unscoped() 确保物理删除

		// 删除用户
		if err := tx.Unscoped().Where("app_id = ?", appID).Delete(&model.User{}).Error; err != nil {
			return fmt.Errorf("failed to delete users: %w", err)
		}

		// 删除角色
		if err := tx.Unscoped().Where("app_id = ?", appID).Delete(&model.Role{}).Error; err != nil {
			return fmt.Errorf("failed to delete roles: %w", err)
		}

		// 删除菜单
		if err := tx.Unscoped().Where("app_id = ?", appID).Delete(&model.Menu{}).Error; err != nil {
			return fmt.Errorf("failed to delete menus: %w", err)
		}

		// 删除API权限
		if err := tx.Unscoped().Where("app_id = ?", appID).Delete(&model.ApiPermission{}).Error; err != nil {
			return fmt.Errorf("failed to delete api permissions: %w", err)
		}

		// 删除配置字典
		if err := tx.Unscoped().Where("app_id = ?", appID).Delete(&model.ConfigDictionary{}).Error; err != nil {
			return fmt.Errorf("failed to delete config dictionaries: %w", err)
		}

		// 删除审计日志 (可选，视需求而定，这里选择一并清理以保持干净)
		if err := tx.Unscoped().Where("app_id = ?", appID).Delete(&model.AuditLog{}).Error; err != nil {
			return fmt.Errorf("failed to delete audit logs: %w", err)
		}

		// 3. 最后删除应用本身
		result := tx.Unscoped().Delete(&model.Application{}, appID)
		if result.Error != nil {
			return fmt.Errorf("failed to delete application: %w", result.Error)
		}
		
		fmt.Printf("Delete result: %d rows affected\n", result.RowsAffected)
		return nil
	})
}

// GetApplicationByID 根据ID获取应用
func (s *ApplicationService) GetApplicationByID(id string) (*model.Application, error) {
	appID := s.parseID(id)
	var app model.Application
	if err := s.DB.First(&app, appID).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// GetApplicationByUUID 根据UUID获取应用
func (s *ApplicationService) GetApplicationByUUID(uuid string) (*model.Application, error) {
	var app model.Application
	if err := s.DB.Where("uuid = ?", uuid).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// GetApplicationByIDWithoutSecret 根据ID获取应用（不包含密钥）
func (s *ApplicationService) GetApplicationByIDWithoutSecret(id string) (*model.Application, error) {
	app, err := s.GetApplicationByID(id)
	if err != nil {
		return nil, err
	}

	// 创建不包含密钥的副本
	appCopy := *app
	appCopy.SecretKey = ""

	return &appCopy, nil
}

// parseID 解析字符串ID为uint
func (s *ApplicationService) parseID(idStr string) uint {
	var id uint
	// 尝试解析为uint
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		// 如果解析失败，返回0，表示无效ID
		return 0
	}
	return id
}

// GetApplicationByCode 根据代码获取应用
func (s *ApplicationService) GetApplicationByCode(code string) (*model.Application, error) {
	var app model.Application
	if err := s.DB.Where("code = ?", code).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// GetApplicationByCodeWithoutSecret 根据代码获取应用（不包含密钥）
func (s *ApplicationService) GetApplicationByCodeWithoutSecret(code string) (*model.Application, error) {
	app, err := s.GetApplicationByCode(code)
	if err != nil {
		return nil, err
	}

	// 创建不包含密钥的副本
	appCopy := *app
	appCopy.SecretKey = ""

	return &appCopy, nil
}

// ListApplications 列出所有应用
func (s *ApplicationService) ListApplications() ([]*model.Application, error) {
	var apps []*model.Application
	if err := s.DB.Order("id asc").Find(&apps).Error; err != nil {
		return nil, err
	}

	// 清除所有应用的密钥
	for _, app := range apps {
		app.SecretKey = ""
	}

	return apps, nil
}

// ResetSecretKey 重置应用密钥
func (s *ApplicationService) ResetSecretKey(id uint) (string, error) {
	app, err := s.GetApplicationByID(fmt.Sprintf("%d", id))
	if err != nil {
		return "", fmt.Errorf("failed to get application: %w", err)
	}

	newSecretKey := s.generateSecretKey()
	app.SecretKey = newSecretKey

	if err := s.DB.Save(app).Error; err != nil {
		return "", fmt.Errorf("failed to update secret key: %w", err)
	}

	return newSecretKey, nil
}
