package service

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"github.com/glebarez/sqlite"
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
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.ExternalApp{},
		&model.ExternalAppPermission{},
		&model.ExternalAppToken{},
		&model.ExternalAppMenu{},
		&model.ExternalAppRole{},
		&model.ExternalAppRoleMenu{},
		&model.ExternalAppUser{},
		&model.ExternalAppUserRole{},
		&model.ExternalAppPermission{},
		&model.ExternalAppRolePermission{},
		// CasbinRule 会被 Gorm Adapter 自动迁移
	)
}

// seedData 初始化种子数据
func seedData(db *gorm.DB) error {
	// 检查是否已经有数据
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount > 0 {
		return nil // 已有数据，跳过种子初始化
	}

	// 创建超级管理员角色
	adminRole := &model.Role{
		Code: "admin",
		Name: "超级管理员",
		Sort: 1,
	}
	if err := db.Create(adminRole).Error; err != nil {
		return err
	}

	// 创建测试角色
	testRole := &model.Role{
		Code: "test",
		Name: "测试角色",
		Sort: 2,
	}
	if err := db.Create(testRole).Error; err != nil {
		return err
	}

	// 创建菜单
	menus := []model.Menu{
		{
			Name:      "系统管理",
			Path:      "/system",
			Component: "Layout",
			Type:      0, // Directory
			Sort:      1,
			Hidden:    false,
		},
		{
			ParentID:  1,
			Name:      "用户管理",
			Path:      "/system/user",
			Component: "system/user/index",
			Type:      1, // Menu
			Sort:      1,
			Hidden:    false,
		},
		{
			ParentID:  1,
			Name:      "角色管理",
			Path:      "/system/role",
			Component: "system/role/index",
			Type:      1, // Menu
			Sort:      2,
			Hidden:    false,
		},
		{
			ParentID:  1,
			Name:      "菜单管理",
			Path:      "/system/menu",
			Component: "system/menu/index",
			Type:      1, // Menu
			Sort:      3,
			Hidden:    false,
		},
	}
	if err := db.Create(&menus).Error; err != nil {
		return err
	}

	// 为超级管理员分配菜单
	if err := db.Model(adminRole).Association("Menus").Append(&menus); err != nil {
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
