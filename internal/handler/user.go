package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"Authos/internal/model"
	"Authos/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c echo.Context) error {
	// 定义创建用户请求结构体
	type CreateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Status   int    `json:"status"`
		RoleIDs  []uint `json:"roleIds"`
	}

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 数据验证
	if req.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username is required"})
	}
	if len(req.Username) > 50 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username cannot exceed 50 characters"})
	}
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Password is required"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	// 创建用户对象
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Status:   req.Status,
		RoleIDs:  req.RoleIDs,
		AppID:    appID,
	}

	if err := h.UserService.CreateUser(user); err != nil {
		// 记录详细错误信息
		service.Log.Errorf("Failed to create user: %v, username=%s, appID=%d", err, user.Username, user.AppID)

		// 根据错误类型返回不同的错误信息
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.JSON(http.StatusConflict, map[string]string{"message": fmt.Sprintf("Username '%s' already exists", user.Username)})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusConflict, map[string]string{"message": fmt.Sprintf("Username '%s' already exists", user.Username)})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Failed to create user: %v", err)})
	}

	// 记录审计日志
	userIDInterface := c.Get("userID")
	usernameInterface := c.Get("username")
	var userID uint
	var username string

	if userIDInterface != nil {
		if u, ok := userIDInterface.(uint); ok {
			userID = u
		} else if f, ok := userIDInterface.(float64); ok {
			userID = uint(f)
		}
	}

	if usernameInterface != nil {
		if s, ok := usernameInterface.(string); ok {
			username = s
		}
	}

	h.UserService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "CREATE",
		Resource:   "USER",
		ResourceID: fmt.Sprintf("%d", user.ID),
		Content:    fmt.Sprintf("创建用户: %s", user.Username),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user":    user,
		"message": "User created successfully",
	})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	// 定义更新请求结构
	type UpdateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password,omitempty"`
		Status   int    `json:"status"`
		RoleIDs  []uint `json:"roleIds"`
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 创建用户对象并设置ID和RoleIDs
	user := &model.User{
		Username: req.Username,
		Status:   req.Status,
		RoleIDs:  req.RoleIDs,
	}
	// 设置ID（ID字段来自嵌入式gorm.Model）
	user.ID = uint(id)

	// 更新用户信息
	if err := h.UserService.UpdateUser(user); err != nil {
		service.Log.Errorf("Failed to update user: %v, userID=%d, username=%s", err, user.ID, user.Username)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user"})
	}

	// 记录审计日志
	userIDInterface := c.Get("userID")
	usernameInterface := c.Get("username")
	var userID uint
	var username string

	if userIDInterface != nil {
		if u, ok := userIDInterface.(uint); ok {
			userID = u
		} else if f, ok := userIDInterface.(float64); ok {
			userID = uint(f)
		}
	}

	if usernameInterface != nil {
		if s, ok := usernameInterface.(string); ok {
			username = s
		}
	}

	h.UserService.DB.Create(&model.AuditLog{
		AppID:      user.AppID, // 这里如果是系统管理员修改可能需要从context拿appID
		UserID:     userID,
		Username:   username,
		Action:     "UPDATE",
		Resource:   "USER",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("更新用户: %s", user.Username),
		IP:         c.RealIP(),
		Status:     1,
	})

	// 如果提供了密码，则更新密码
	if req.Password != "" {
		if err := h.UserService.UpdateUserPassword(uint(id), req.Password); err != nil {
			service.Log.Errorf("Failed to update user password: %v, userID=%d", err, id)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update password"})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":    user,
		"message": "User updated successfully",
	})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	if err := h.UserService.DeleteUser(uint(id), appID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete user"})
	}

	// 记录审计日志
	userIDInterface := c.Get("userID")
	usernameInterface := c.Get("username")
	var userID uint
	var username string

	if userIDInterface != nil {
		if u, ok := userIDInterface.(uint); ok {
			userID = u
		} else if f, ok := userIDInterface.(float64); ok {
			userID = uint(f)
		}
	}

	if usernameInterface != nil {
		if s, ok := usernameInterface.(string); ok {
			username = s
		}
	}

	h.UserService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "DELETE",
		Resource:   "USER",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("删除用户ID: %d", id),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// GetUser 获取用户
func (h *UserHandler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	user, err := h.UserService.GetUserByID(uint(id), appID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// ListUsers 列出所有用户
func (h *UserHandler) ListUsers(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	username := c.QueryParam("username")
	status := c.QueryParam("status")

	var users []*model.User
	db := h.UserService.DB.Preload("Roles").Where("app_id = ?", appID)
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Order("id desc").Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get users"})
	}

	// 填充 RoleIDs
	for _, user := range users {
		user.RoleIDs = make([]uint, 0, len(user.Roles))
		for _, role := range user.Roles {
			user.RoleIDs = append(user.RoleIDs, role.ID)
		}
	}

	return c.JSON(http.StatusOK, users)
}
