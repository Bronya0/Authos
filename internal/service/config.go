package service

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Log    LogConfig    `yaml:"log"`
	System SystemConfig `yaml:"system"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type LogConfig struct {
	Dir      string `yaml:"dir"`
	Filename string `yaml:"filename"`
}

type SystemConfig struct {
	AdminUsername string `yaml:"adminUsername"`
	AdminPassword string `yaml:"adminPassword"`
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Log.Dir == "" {
		config.Log.Dir = "logs"
	}
	if config.Log.Filename == "" {
		config.Log.Filename = "authos.log"
	}

	// 设置系统管理员默认值
	if config.System.AdminUsername == "" {
		config.System.AdminUsername = "admin"
	}
	// 注意：密码如果为空，后续逻辑应生成随机密码，这里不设默认值以免覆盖用户意图

	return &config, nil
}
