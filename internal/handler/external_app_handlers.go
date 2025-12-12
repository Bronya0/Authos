package handler

import (
	"net/http"
	"strconv"

	"Authos/internal/model"
	"Authos/internal/service"

	"github.com/labstack/echo/v4"
)

// ExternalAppMenuHandler 外部应用菜单处理器
type ExternalAppMenuHandler struct {
	menuService *service.ExternalAppMenuService
}

// NewExternalAppMenuHandler 创建外部应用菜单处理器实例
func NewExternalAppMenuHandler(menuService *service.ExternalAppMenuService) *ExternalAppMenuHandler {
	return &ExternalAppMenuHandler{
		menuService: menuService,
	}
}

// GetMenus 获取菜单列表
func (h *ExternalAppMenuHandler) GetMenus(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	menus, err := h.menuService.GetMenuTree(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取菜单列表失败"})
	}

	return c.JSON(http.StatusOK, menus)
}

// CreateMenu 创建菜单
func (h *ExternalAppMenuHandler) CreateMenu(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	var menu model.ExternalAppMenu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	menu.AppID = appID

	if err := h.menuService.CreateMenu(&menu); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "创建菜单失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "菜单创建成功",
		"data":    menu,
	})
}

// UpdateMenu 更新菜单
func (h *ExternalAppMenuHandler) UpdateMenu(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的菜单ID"})
	}

	var menu model.ExternalAppMenu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	// 设置菜单ID和应用ID
	menu.ID = uint(id)
	menu.AppID = appID

	if err := h.menuService.UpdateMenu(&menu); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "更新菜单失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "菜单更新成功"})
}

// DeleteMenu 删除菜单
func (h *ExternalAppMenuHandler) DeleteMenu(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的菜单ID"})
	}

	// 检查菜单是否属于当前应用
	menu, err := h.menuService.GetMenuByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "菜单不存在"})
	}

	if menu.AppID != appID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "无权限操作此菜单"})
	}

	if err := h.menuService.DeleteMenu(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "删除菜单失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "菜单删除成功"})
}

// ExternalAppRoleHandler 外部应用角色处理器
type ExternalAppRoleHandler struct {
	roleService *service.ExternalAppRoleService
}

// NewExternalAppRoleHandler 创建外部应用角色处理器实例
func NewExternalAppRoleHandler(roleService *service.ExternalAppRoleService) *ExternalAppRoleHandler {
	return &ExternalAppRoleHandler{
		roleService: roleService,
	}
}

// GetRoles 获取角色列表
func (h *ExternalAppRoleHandler) GetRoles(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	roles, err := h.roleService.GetRolesByAppID(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取角色列表失败"})
	}

	return c.JSON(http.StatusOK, roles)
}

// CreateRole 创建角色
func (h *ExternalAppRoleHandler) CreateRole(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	var role model.ExternalAppRole
	if err := c.Bind(&role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	role.AppID = appID

	if err := h.roleService.CreateRole(&role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "创建角色失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "角色创建成功",
		"data":    role,
	})
}

// UpdateRole 更新角色
func (h *ExternalAppRoleHandler) UpdateRole(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的角色ID"})
	}

	var role model.ExternalAppRole
	if err := c.Bind(&role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	// 设置角色ID和应用ID
	role.ID = uint(id)
	role.AppID = appID

	if err := h.roleService.UpdateRole(&role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "更新角色失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "角色更新成功"})
}

// DeleteRole 删除角色
func (h *ExternalAppRoleHandler) DeleteRole(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的角色ID"})
	}

	// 检查角色是否属于当前应用
	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "角色不存在"})
	}

	if role.AppID != appID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "无权限操作此角色"})
	}

	if err := h.roleService.DeleteRole(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "删除角色失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "角色删除成功"})
}

// AssignRoleMenus 分配角色菜单
func (h *ExternalAppRoleHandler) AssignRoleMenus(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的角色ID"})
	}

	// 检查角色是否属于当前应用
	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "角色不存在"})
	}

	if role.AppID != appID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "无权限操作此角色"})
	}

	var req struct {
		MenuIDs []uint `json:"menu_ids"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.roleService.AssignMenus(uint(id), req.MenuIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "分配菜单权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "菜单权限分配成功"})
}

// ExternalAppUserHandler 外部应用用户处理器
type ExternalAppUserHandler struct {
	userService *service.ExternalAppUserService
}

// NewExternalAppUserHandler 创建外部应用用户处理器实例
func NewExternalAppUserHandler(userService *service.ExternalAppUserService) *ExternalAppUserHandler {
	return &ExternalAppUserHandler{
		userService: userService,
	}
}

// GetUsers 获取用户列表
func (h *ExternalAppUserHandler) GetUsers(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	users, err := h.userService.GetUsersByAppID(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取用户列表失败"})
	}

	return c.JSON(http.StatusOK, users)
}

// GetUser 获取单个用户
func (h *ExternalAppUserHandler) GetUser(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的用户ID"})
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "用户不存在"})
	}

	if user.AppID != appID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "无权限访问此用户"})
	}

	return c.JSON(http.StatusOK, user)
}

// CreateUser 创建用户
func (h *ExternalAppUserHandler) CreateUser(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	var req struct {
		model.ExternalAppUser
		RoleIDs []uint `json:"role_ids"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	user := req.ExternalAppUser
	user.AppID = appID

	if err := h.userService.CreateUserWithRoles(&user, req.RoleIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "创建用户失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "用户创建成功",
		"data":    user,
	})
}

// UpdateUser 更新用户
func (h *ExternalAppUserHandler) UpdateUser(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的用户ID"})
	}

	var req struct {
		model.ExternalAppUser
		RoleIDs []uint `json:"role_ids"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	user := req.ExternalAppUser
	// 设置用户ID和应用ID
	user.ID = uint(id)
	user.AppID = appID

	if err := h.userService.UpdateUserWithRoles(&user, req.RoleIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "更新用户失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "用户更新成功"})
}

// DeleteUser 删除用户
func (h *ExternalAppUserHandler) DeleteUser(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的用户ID"})
	}

	// 检查用户是否属于当前应用
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "用户不存在"})
	}

	if user.AppID != appID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "无权限操作此用户"})
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "删除用户失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "用户删除成功"})
}

// AssignUserRoles 分配用户角色
func (h *ExternalAppUserHandler) AssignUserRoles(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的用户ID"})
	}

	// 检查用户是否属于当前应用
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "用户不存在"})
	}

	if user.AppID != appID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "无权限操作此用户"})
	}

	var req struct {
		RoleIDs []uint `json:"role_ids"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.userService.AssignRoles(uint(id), req.RoleIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "分配角色失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "角色分配成功"})
}

// ExternalAppPermissionHandler 外部应用权限处理器
type ExternalAppPermissionHandler struct {
	permissionService *service.ExternalAppPermissionService
}

// NewExternalAppPermissionHandler 创建外部应用权限处理器实例
func NewExternalAppPermissionHandler(permissionService *service.ExternalAppPermissionService) *ExternalAppPermissionHandler {
	return &ExternalAppPermissionHandler{
		permissionService: permissionService,
	}
}

// GetPermissions 获取权限列表
func (h *ExternalAppPermissionHandler) GetPermissions(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	permissions, err := h.permissionService.GetPermissionsByAppID(appID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取权限列表失败"})
	}

	return c.JSON(http.StatusOK, permissions)
}

// CreatePermission 创建权限
func (h *ExternalAppPermissionHandler) CreatePermission(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	var permission model.ExternalAppPermission
	if err := c.Bind(&permission); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	permission.AppID = appID

	if err := h.permissionService.CreatePermission(&permission); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "创建权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "权限创建成功",
		"data":    permission,
	})
}

// UpdatePermission 更新权限
func (h *ExternalAppPermissionHandler) UpdatePermission(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的权限ID"})
	}

	var permission model.ExternalAppPermission
	if err := c.Bind(&permission); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	// 设置权限ID和应用ID
	permission.ID = uint(id)
	permission.AppID = appID

	if err := h.permissionService.UpdatePermission(&permission); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "更新权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "权限更新成功"})
}

// DeletePermission 删除权限
func (h *ExternalAppPermissionHandler) DeletePermission(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的权限ID"})
	}

	// 检查权限是否属于当前应用
	permission, err := h.permissionService.GetPermissionByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "权限不存在"})
	}

	if permission.AppID != appID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "无权限操作此权限"})
	}

	if err := h.permissionService.DeletePermission(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "删除权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "权限删除成功"})
}

// AssignRolePermissions 分配角色权限
func (h *ExternalAppPermissionHandler) AssignRolePermissions(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的角色ID"})
	}

	// 检查角色是否属于当前应用
	roleService := service.NewExternalAppRoleService(h.permissionService.DB)
	role, err := roleService.GetRoleByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "角色不存在"})
	}

	if role.AppID != appID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "无权限操作此角色"})
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	if err := h.permissionService.AssignPermissions(uint(id), req.PermissionIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "分配权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "权限分配成功"})
}

// CheckUserPermission 检查用户权限
func (h *ExternalAppPermissionHandler) CheckUserPermission(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	var req struct {
		Username string `json:"username"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	// 获取用户
	userService := service.NewExternalAppUserService(h.permissionService.DB)
	user, err := userService.GetUserByUsername(appID, req.Username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "用户不存在"})
	}

	// 检查权限
	hasPermission, err := h.permissionService.CheckUserPermission(user.ID, req.Resource, req.Action)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "检查权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"has_permission": hasPermission,
	})
}

// CheckUserMenuPermission 检查用户菜单权限
func (h *ExternalAppPermissionHandler) CheckUserMenuPermission(c echo.Context) error {
	app := c.Get("app").(*model.ExternalApp)
	appID := app.ID

	var req struct {
		Username string `json:"username"`
		MenuID   uint   `json:"menu_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的请求数据"})
	}

	// 获取用户
	userService := service.NewExternalAppUserService(h.permissionService.DB)
	user, err := userService.GetUserByUsername(appID, req.Username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "用户不存在"})
	}

	// 检查菜单权限
	hasPermission, err := h.permissionService.CheckUserMenuPermission(user.ID, req.MenuID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "检查菜单权限失败"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"has_permission": hasPermission,
	})
}
