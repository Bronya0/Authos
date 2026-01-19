package service

import (
	"fmt"

	"Authos/internal/model"

	"gorm.io/gorm"
)

type ConfigDictionaryService struct {
	DB *gorm.DB
}

func NewConfigDictionaryService(db *gorm.DB) *ConfigDictionaryService {
	return &ConfigDictionaryService{DB: db}
}

func (s *ConfigDictionaryService) GetConfigDictionary(id uint, appID uint) (*model.ConfigDictionary, error) {
	var item model.ConfigDictionary
	if err := s.DB.Where("id = ? AND app_id = ?", id, appID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ConfigDictionaryService) CreateConfigDictionary(appID uint, key, value, desc string) (*model.ConfigDictionary, error) {
	if key == "" {
		return nil, fmt.Errorf("字典key不能为空")
	}

	var existing model.ConfigDictionary
	if err := s.DB.Where("key = ? AND app_id = ?", key, appID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("字典key已存在: %s", key)
	}

	item := model.ConfigDictionary{
		Key:   key,
		Value: value,
		Desc:  desc,
		AppID: appID,
	}
	if err := s.DB.Create(&item).Error; err != nil {
		return nil, fmt.Errorf("创建配置字典失败: %v", err)
	}
	return &item, nil
}

func (s *ConfigDictionaryService) UpdateConfigDictionary(id uint, appID uint, key, value, desc string) (*model.ConfigDictionary, error) {
	item, err := s.GetConfigDictionary(id, appID)
	if err != nil {
		return nil, err
	}

	if key == "" {
		return nil, fmt.Errorf("字典key不能为空")
	}

	var existing model.ConfigDictionary
	if err := s.DB.Where("key = ? AND id != ? AND app_id = ?", key, id, appID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("字典key已存在: %s", key)
	}

	item.Key = key
	item.Value = value
	item.Desc = desc

	if err := s.DB.Save(item).Error; err != nil {
		return nil, fmt.Errorf("更新配置字典失败: %v", err)
	}
	return item, nil
}

func (s *ConfigDictionaryService) DeleteConfigDictionary(id uint, appID uint) error {
	if err := s.DB.Where("id = ? AND app_id = ?", id, appID).Delete(&model.ConfigDictionary{}).Error; err != nil {
		return fmt.Errorf("删除配置字典失败: %v", err)
	}
	return nil
}
