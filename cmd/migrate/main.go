package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"Authos/internal/model"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("data/app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 开始迁移
	fmt.Println("Starting migration from Code to UUID...")

	// 1. 添加UUID列
	if err := db.Exec("ALTER TABLE roles ADD COLUMN uuid VARCHAR(36)").Error; err != nil {
		fmt.Println("UUID column might already exist or there was an error:", err)
	}

	// 2. 为所有现有角色生成UUID
	var roles []model.Role
	if err := db.Find(&roles).Error; err != nil {
		log.Fatal("Failed to fetch roles:", err)
	}

	for _, role := range roles {
		// 生成UUID
		newUUID := uuid.New().String()
		
		// 更新角色
		if err := db.Model(&role).Update("uuid", newUUID).Error; err != nil {
			log.Printf("Failed to update role %d: %v", role.ID, err)
		} else {
			fmt.Printf("Updated role %s with UUID %s\n", role.Name, newUUID)
		}
	}

	// 3. 更新Casbin策略，将role:ID替换为role:UUID
	policies, err := db.Raw("SELECT ptype, v0, v1, v2 FROM casbin_rule").Rows()
	if err != nil {
		log.Fatal("Failed to fetch Casbin policies:", err)
	}
	defer policies.Close()

	for policies.Next() {
		var ptype, v0, v1, v2 string
		if err := policies.Scan(&ptype, &v0, &v1, &v2); err != nil {
			log.Printf("Failed to scan policy: %v", err)
			continue
		}

		// 检查是否是角色策略
		if len(v0) > 5 && v0[:5] == "role:" {
			// 提取角色ID
			roleIDStr := v0[5:]
			var roleID uint
			if _, err := fmt.Sscanf(roleIDStr, "%d", &roleID); err == nil {
				// 获取角色UUID
				var role model.Role
				if err := db.First(&role, roleID).Error; err == nil {
					// 更新策略
					newV0 := fmt.Sprintf("role:%s", role.UUID)
					if err := db.Exec("UPDATE casbin_rule SET v0 = ? WHERE ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", 
						newV0, ptype, v0, v1, v2).Error; err != nil {
						log.Printf("Failed to update policy: %v", err)
					} else {
						fmt.Printf("Updated policy from %s to %s\n", v0, newV0)
					}
				}
			}
		}
	}

	fmt.Println("Migration completed successfully!")
}