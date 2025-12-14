package service

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"Authos/internal/model"
)

// DBService 数据库服务
type DBService struct {
	DB *gorm.DB
}

// NewDBService 创建数据库服务实例
func NewDBService() (*DBService, error) {
	// 数据库文件路径
	dbPath := "auth.db"

	// 配置 GORM 日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)

	// 连接到 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移模型
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	// 初始化种子数据
	if err := seedData(db); err != nil {
		return nil, fmt.Errorf("failed to seed data: %w", err)
	}

	return &DBService{
		DB: db,
	}, nil
}

// autoMigrate 自动迁移模型
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Application{},
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.ApiPermission{},
		// CasbinRule 会被 Gorm Adapter 自动迁移
	)
}

// seedData 初始化种子数据
func seedData(db *gorm.DB) error {
	// 检查是否已经有数据
	var appCount int64
	db.Model(&model.Application{}).Count(&appCount)
	if appCount > 0 {
		return nil // 已有数据，跳过种子初始化
	}

	// 创建默认应用
	defaultApp := &model.Application{
		Name:        "默认应用",
		Code:        "default",
		SecretKey:   "default-secret-key",
		Status:      1,
		Description: "系统默认应用",
	}
	if err := db.Create(defaultApp).Error; err != nil {
		return err
	}

	// 创建超级管理员角色
	adminRole := &model.Role{
		Name:  "超级管理员",
		AppID: defaultApp.ID,
	}
	if err := db.Create(adminRole).Error; err != nil {
		return err
	}

	// 创建测试角色
	testRole := &model.Role{
		Name:  "测试角色",
		AppID: defaultApp.ID,
	}
	if err := db.Create(testRole).Error; err != nil {
		return err
	}

	// 创建超级管理员用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminUser := &model.User{
		Username: "admin",
		Password: string(hashedPassword),
		Status:   1,
		AppID:    defaultApp.ID,
	}
	if err := db.Create(adminUser).Error; err != nil {
		return err
	}

	// 为超级管理员用户分配角色
	if err := db.Model(adminUser).Association("Roles").Append(adminRole); err != nil {
		return err
	}

	return nil
}
