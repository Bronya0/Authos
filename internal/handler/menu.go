package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Authos/internal/model"
	"Authos/internal/service"
)

// MenuHandler 菜单处理器
type MenuHandler struct {
	MenuService *service.MenuService
}

// NewMenuHandler 创建菜单处理器实例
func NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{MenuService: menuService}
}

// CreateMenu 创建菜单
func (h *MenuHandler) CreateMenu(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	var menu model.Menu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	menu.AppID = appID // 设置应用ID

	if err := h.MenuService.CreateMenu(&menu); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create menu"})
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

	h.MenuService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "CREATE",
		Resource:   "MENU",
		ResourceID: fmt.Sprintf("%d", menu.ID),
		Content:    fmt.Sprintf("创建菜单: %s", menu.Name),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"menu":    menu,
		"message": "Menu created successfully",
	})
}

// UpdateMenu 更新菜单
func (h *MenuHandler) UpdateMenu(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid menu ID"})
	}

	var menu model.Menu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	menu.ID = uint(id)
	menu.AppID = appID

	if err := h.MenuService.UpdateMenu(&menu); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update menu"})
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

	h.MenuService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "UPDATE",
		Resource:   "MENU",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("更新菜单: %s", menu.Name),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"menu":    menu,
		"message": "Menu updated successfully",
	})
}

// DeleteMenu 删除菜单
func (h *MenuHandler) DeleteMenu(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid menu ID"})
	}

	if err := h.MenuService.DeleteMenu(uint(id), appID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete menu"})
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

	h.MenuService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "DELETE",
		Resource:   "MENU",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("删除菜单ID: %d", id),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Menu deleted successfully"})
}

// GetMenu 获取菜单
func (h *MenuHandler) GetMenu(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid menu ID"})
	}

	menu, err := h.MenuService.GetMenuByID(uint(id), appID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Menu not found"})
	}

	return c.JSON(http.StatusOK, menu)
}

// ListMenus 列出所有菜单（扁平结构）
func (h *MenuHandler) ListMenus(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	name := c.QueryParam("name")

	db := h.MenuService.DB.Where("app_id = ?", appID)
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	var menus []*model.Menu
	if err := db.Order("sort asc").Find(&menus).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get menus"})
	}

	return c.JSON(http.StatusOK, menus)
}

// GetMenuTree 获取菜单树
func (h *MenuHandler) GetMenuTree(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	menuTree, err := h.MenuService.GetMenuTreeByApp(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get menu tree"})
	}

	return c.JSON(http.StatusOK, menuTree)
}

// GetNonSystemMenuTree 获取非系统菜单树
func (h *MenuHandler) GetNonSystemMenuTree(c echo.Context) error {
	// 从 JWT token 中获取 appID
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	menuTree, err := h.MenuService.GetNonSystemMenuTreeByApp(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get non-system menu tree"})
	}

	return c.JSON(http.StatusOK, menuTree)
}
