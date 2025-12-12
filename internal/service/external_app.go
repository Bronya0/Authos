package service

import (
	"fmt"
	"time"

	"Authos/internal/model"

	"gorm.io/gorm"
)

// ExternalAppService 外部应用服务
type ExternalAppService struct {
	db *gorm.DB
}

// NewExternalAppService 创建外部应用服务实例
func NewExternalAppService(db *gorm.DB) *ExternalAppService {
	return &ExternalAppService{db: db}
}

// CreateExternalApp 创建外部应用
func (s *ExternalAppService) CreateExternalApp(app *model.ExternalApp) error {
	// 生成应用密钥
	app.AppKey = GenerateAppKey()
	app.AppSecret = GenerateAppSecret()

	return s.db.Create(app).Error
}

// ListExternalApps 获取外部应用列表
func (s *ExternalAppService) ListExternalApps() ([]model.ExternalApp, error) {
	var apps []model.ExternalApp
	err := s.db.Find(&apps).Error
	return apps, err
}

// GetExternalApp 根据ID获取外部应用
func (s *ExternalAppService) GetExternalApp(id uint) (*model.ExternalApp, error) {
	var app model.ExternalApp
	err := s.db.First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetExternalAppByKey 根据AppKey获取外部应用
func (s *ExternalAppService) GetExternalAppByKey(appKey string) (*model.ExternalApp, error) {
	var app model.ExternalApp
	err := s.db.Where("app_key = ?", appKey).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// UpdateExternalApp 更新外部应用
func (s *ExternalAppService) UpdateExternalApp(id uint, app *model.ExternalApp) error {
	return s.db.Model(&model.ExternalApp{}).Where("id = ?", id).Updates(app).Error
}

// DeleteExternalApp 删除外部应用
func (s *ExternalAppService) DeleteExternalApp(id uint) error {
	return s.db.Delete(&model.ExternalApp{}, id).Error
}

// CreateAppPermission 创建应用权限
func (s *ExternalAppService) CreateAppPermission(permission *model.ExternalAppPermission) error {
	return s.db.Create(permission).Error
}

// GetAppPermissions 获取应用权限列表
func (s *ExternalAppService) GetAppPermissions(appID uint) ([]model.ExternalAppPermission, error) {
	var permissions []model.ExternalAppPermission
	err := s.db.Where("app_id = ?", appID).Find(&permissions).Error
	return permissions, err
}

// UpdateAppPermission 更新应用权限
func (s *ExternalAppService) UpdateAppPermission(id uint, permission *model.ExternalAppPermission) error {
	return s.db.Model(&model.ExternalAppPermission{}).Where("id = ?", id).Updates(permission).Error
}

// DeleteAppPermission 删除应用权限
func (s *ExternalAppService) DeleteAppPermission(id uint) error {
	return s.db.Delete(&model.ExternalAppPermission{}, id).Error
}

// CreateAppToken 创建应用令牌
func (s *ExternalAppService) CreateAppToken(appID uint, expiresAt time.Time) (string, error) {
	token := GenerateToken()

	appToken := &model.ExternalAppToken{
		AppID:     appID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	err := s.db.Create(appToken).Error
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateAppToken 验证应用令牌
func (s *ExternalAppService) ValidateAppToken(token string) (*model.ExternalApp, error) {
	var appToken model.ExternalAppToken
	err := s.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&appToken).Error
	if err != nil {
		return nil, err
	}

	var app model.ExternalApp
	err = s.db.First(&app, appToken.AppID).Error
	if err != nil {
		return nil, err
	}

	if app.Status != 1 {
		return nil, fmt.Errorf("应用已禁用")
	}

	return &app, nil
}
