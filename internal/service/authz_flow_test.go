package service

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"Authos/internal/model"
)

func newTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(
		&model.Application{},
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.ApiPermission{},
		&model.AuditLog{},
	); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	return db
}

func TestPermissionBindingAndCheckByKey(t *testing.T) {
	db := newTestDB(t)

	casbinService, err := NewCasbinService(db)
	if err != nil {
		t.Fatalf("failed to create casbin service: %v", err)
	}

	apiPermissionService := NewApiPermissionService(db, casbinService, nil)

	app := &model.Application{
		Name:      "test-app",
		Code:      "test-app",
		SecretKey: "secret",
		Status:    1,
	}
	if err := db.Create(app).Error; err != nil {
		t.Fatalf("failed to create application: %v", err)
	}

	role := &model.Role{
		Name:  "test-role",
		AppID: app.ID,
	}
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("failed to create role: %v", err)
	}

	user := &model.User{
		Username: "test-user",
		Password: "password",
		Status:   1,
		AppID:    app.ID,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if err := db.Model(user).Association("Roles").Append(role); err != nil {
		t.Fatalf("failed to associate role to user: %v", err)
	}

	permission, err := apiPermissionService.CreateApiPermission(
		app.ID,
		"user:create",
		"创建用户",
		"/api/v1/users",
		model.HTTP_ALL,
		"",
	)
	if err != nil {
		t.Fatalf("failed to create api permission: %v", err)
	}

	if err := apiPermissionService.AddApiPermissionToRole(app.ID, role.UUID, permission.UUID); err != nil {
		t.Fatalf("failed to add permission to role: %v", err)
	}

	allowed, err := casbinService.CheckPermission(user.ID, "user:create", model.HTTP_ALL)
	if err != nil {
		t.Fatalf("check permission error: %v", err)
	}
	if !allowed {
		t.Fatalf("expected permission allowed for user:create")
	}

	if err := apiPermissionService.RemoveApiPermissionFromRole(app.ID, role.UUID, permission.UUID); err != nil {
		t.Fatalf("failed to remove permission from role: %v", err)
	}

	allowed, err = casbinService.CheckPermission(user.ID, "user:create", model.HTTP_ALL)
	if err != nil {
		t.Fatalf("check permission error after remove: %v", err)
	}
	if allowed {
		t.Fatalf("expected permission denied for user:create after remove")
	}
}
