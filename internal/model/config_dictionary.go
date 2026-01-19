package model

import "gorm.io/gorm"

type ConfigDictionary struct {
	gorm.Model
	Key   string       `gorm:"size:100;not null;index:idx_config_dictionary_app_key" json:"key"`
	Value string       `gorm:"size:500;not null" json:"value"`
	Desc  string       `gorm:"size:500" json:"desc"`
	AppID uint         `gorm:"not null;index:idx_config_dictionary_app_key" json:"appId"`
	App   *Application `gorm:"foreignKey:AppID" json:"app,omitempty"`
}
