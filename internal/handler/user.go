package handler

import (
	"fmt"
	"log"
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
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// 数据验证
	if user.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username is required"})
	}
	if len(user.Username) > 50 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username cannot exceed 50 characters"})
	}
	if user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Password is required"})
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		// 记录详细错误信息
		log.Printf("Failed to create user: %v", err)
		
		// 根据错误类型返回不同的错误信息
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.JSON(http.StatusConflict, map[string]string{"message": fmt.Sprintf("Username '%s' already exists", user.Username)})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusConflict, map[string]string{"message": fmt.Sprintf("Username '%s' already exists", user.Username)})
		}
		
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Failed to create user: %v", err)})
	}

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

	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	user.ID = uint(id)

	if err := h.UserService.UpdateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user"})
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

	if err := h.UserService.DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// GetUser 获取用户
func (h *UserHandler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	user, err := h.UserService.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// ListUsers 列出所有用户
func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.UserService.ListUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get users"})
	}

	return c.JSON(http.StatusOK, users)
}
