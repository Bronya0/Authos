package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"Authos/internal/model"
	"Authos/internal/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ExternalAPIHandler 外部API处理器
type ExternalAPIHandler struct {
	externalAppService *service.ExternalAppService
	menuService        *service.MenuService
	roleService        *service.RoleService
}

// NewExternalAPIHandler 创建外部API处理器实例
func NewExternalAPIHandler(
	externalAppService *service.ExternalAppService,
	menuService *service.MenuService,
	roleService *service.RoleService,
) *ExternalAPIHandler {
	return &ExternalAPIHandler{
		externalAppService: externalAppService,
		menuService:        menuService,
		roleService:        roleService,
	}
}

// RegisterApp 注册外部应用
func (h *ExternalAPIHandler) RegisterApp(c echo.Context) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "应用名称不能为空"})
	}

	app := &model.ExternalApp{
		Name:        req.Name,
		Description: req.Description,
		Status:      1, // 默认启用
	}

	if err := h.externalAppService.CreateExternalApp(app); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "注册应用失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"app_id":     app.ID,
		"app_key":    app.AppKey,
		"app_secret": app.AppSecret,
		"message":    "应用注册成功",
	})
}

// GetAppToken 获取应用令牌
func (h *ExternalAPIHandler) GetAppToken(c echo.Context) error {
	var req struct {
		AppKey    string `json:"app_key"`
		AppSecret string `json:"app_secret"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if req.AppKey == "" || req.AppSecret == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "应用密钥不能为空"})
	}

	// 验证应用密钥
	app, err := h.externalAppService.GetExternalAppByKey(req.AppKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "应用不存在"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "验证应用失败"})
	}

	if app.AppSecret != req.AppSecret {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "应用密钥错误"})
	}

	if app.Status != 1 {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "应用已禁用"})
	}

	// 生成令牌
	token, err := h.externalAppService.CreateAppToken(app.ID, c.Get("expires_at").(interface{}).(time.Time))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "生成令牌失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_at": c.Get("expires_at"),
	})
}

// GetMenus 获取菜单列表
func (h *ExternalAPIHandler) GetMenus(c echo.Context) error {
	// 检查应用是否有菜单读取权限
	app := c.Get("app").(*model.ExternalApp)
	if !h.hasPermission(app, "menu", "read") {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "没有菜单读取权限"})
	}

	menus, err := h.menuService.ListMenus()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取菜单列表失败"})
	}

	return c.JSON(http.StatusOK, menus)
}

// CreateMenu 创建菜单
func (h *ExternalAPIHandler) CreateMenu(c echo.Context) error {
	// 检查应用是否有菜单创建权限
	app := c.Get("app").(*model.ExternalApp)
	if !h.hasPermission(app, "menu", "create") {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "没有菜单创建权限"})
	}

	var menu model.Menu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.menuService.CreateMenu(&menu); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "创建菜单失败"})
	}

	return c.JSON(http.StatusOK, menu)
}

// UpdateMenu 更新菜单
func (h *ExternalAPIHandler) UpdateMenu(c echo.Context) error {
	// 检查应用是否有菜单更新权限
	app := c.Get("app").(*model.ExternalApp)
	if !h.hasPermission(app, "menu", "update") {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "没有菜单更新权限"})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的菜单ID"})
	}

	var menu model.Menu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	// 设置菜单ID
	menu.ID = uint(id)

	if err := h.menuService.UpdateMenu(&menu); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "更新菜单失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "菜单更新成功"})
}

// DeleteMenu 删除菜单
func (h *ExternalAPIHandler) DeleteMenu(c echo.Context) error {
	// 检查应用是否有菜单删除权限
	app := c.Get("app").(*model.ExternalApp)
	if !h.hasPermission(app, "menu", "delete") {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "没有菜单删除权限"})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的菜单ID"})
	}

	if err := h.menuService.DeleteMenu(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "删除菜单失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "菜单删除成功"})
}

// AssignMenuPermission 分配菜单权限
func (h *ExternalAPIHandler) AssignMenuPermission(c echo.Context) error {
	// 检查应用是否有权限分配权限
	app := c.Get("app").(*model.ExternalApp)
	if !h.hasPermission(app, "role", "assign-menu") {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "没有菜单权限分配权限"})
	}

	var req struct {
		RoleID  uint   `json:"role_id"`
		MenuIDs []uint `json:"menu_ids"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.roleService.AssignMenus(req.RoleID, req.MenuIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "分配菜单权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "菜单权限分配成功"})
}

// AssignAPIPermission 分配API权限
func (h *ExternalAPIHandler) AssignAPIPermission(c echo.Context) error {
	// 检查应用是否有权限分配权限
	app := c.Get("app").(*model.ExternalApp)
	if !h.hasPermission(app, "role", "assign-permission") {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "没有API权限分配权限"})
	}

	var req struct {
		RoleID      uint     `json:"role_id"`
		Permissions []string `json:"permissions"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	// 将字符串权限转换为所需的格式
	var permissions []map[string]string
	for _, perm := range req.Permissions {
		// 假设权限格式为 "resource:action"
		parts := strings.SplitN(perm, ":", 2)
		if len(parts) == 2 {
			permissions = append(permissions, map[string]string{
				"obj": parts[0],
				"act": parts[1],
			})
		}
	}

	if err := h.roleService.AssignPermissions(req.RoleID, permissions); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "分配API权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "API权限分配成功"})
}

// hasPermission 检查应用是否有指定权限
func (h *ExternalAPIHandler) hasPermission(app *model.ExternalApp, resource, action string) bool {
	permissions, err := h.externalAppService.GetAppPermissions(app.ID)
	if err != nil {
		return false
	}

	for _, permission := range permissions {
		if permission.Resource == resource && permission.Action == action {
			return true
		}
	}

	return false
}
