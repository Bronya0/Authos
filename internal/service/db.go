package service

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
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
		&model.AuditLog{},
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
	uuid := uuid.New().String()
	// 创建默认应用
	defaultApp := &model.Application{
		UUID:        uuid,
		Name:        "默认应用",
		Code:        "default",
		SecretKey:   uuid,
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

	// 交互式创建管理员账户
	var adminUsername, adminPassword string

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("检测到没有管理员账户，请选择创建方式：")
	fmt.Println("1. 自动生成 (默认)")
	fmt.Println("2. 手动输入")
	fmt.Print("请输入选项 (1/2): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "2" {
		// 手动输入
		for {
			fmt.Print("请输入管理员用户名 (默认为 admin): ")
			inputUsername, _ := reader.ReadString('\n')
			inputUsername = strings.TrimSpace(inputUsername)
			if inputUsername == "" {
				adminUsername = "admin"
			} else {
				adminUsername = inputUsername
			}

			fmt.Print("请输入管理员密码: ")
			inputPassword, _ := reader.ReadString('\n')
			inputPassword = strings.TrimSpace(inputPassword)
			if inputPassword == "" {
				fmt.Println("密码不能为空，请重新输入")
				continue
			}
			adminPassword = inputPassword
			break
		}
	} else {
		// 自动生成
		adminUsername = "admin"
		adminPassword = generateSecurePassword()
	}

	// 创建管理员用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminUser := &model.User{
		Username: adminUsername,
		Password: string(hashedPassword),
		Status:   1,
		AppID:    defaultApp.ID,
	}
	if err := db.Create(adminUser).Error; err != nil {
		return err
	}

	// 打印初始管理员密码到日志
	log.Printf("========================================")
	log.Printf("初始管理员账户已创建")
	log.Printf("用户名: %s", adminUsername)
	if choice != "2" {
		log.Printf("密码: %s", adminPassword)
	} else {
		log.Printf("密码: (已手动设置)")
	}
	log.Printf("========================================")

	// 为超级管理员用户分配角色
	if err := db.Model(adminUser).Association("Roles").Append(adminRole); err != nil {
		return err
	}

	return nil
}

// generateSecurePassword 生成一个复杂度较高的密码
func generateSecurePassword() string {
	rand.Seed(time.Now().UnixNano())

	// 定义密码字符集
	upperLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerLetters := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?"

	// 密码长度 16 位
	length := 16

	// 确保包含各种字符类型
	password := make([]byte, length)

	// 随机选择各类字符
	password[0] = upperLetters[rand.Intn(len(upperLetters))]
	password[1] = lowerLetters[rand.Intn(len(lowerLetters))]
	password[2] = digits[rand.Intn(len(digits))]
	password[3] = specialChars[rand.Intn(len(specialChars))]

	// 剩余的位置随机填充
	allChars := upperLetters + lowerLetters + digits + specialChars
	for i := 4; i < length; i++ {
		password[i] = allChars[rand.Intn(len(allChars))]
	}

	// 打乱密码顺序
	for i := len(password) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}
