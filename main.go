package main

import (
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"Authos/internal/handler"
	customMiddleware "Authos/internal/middleware"
	"Authos/internal/service"
	"Authos/pkg/utils"
)

// 预留代码：使用 embed.FS 托管前端构建产物
// 假设前端构建产物在 web/dist 目录下
//
//go:embed web/dist
var webFS embed.FS

func main() {
	// 配置信息
	jwtSecret := "authos-secret-key" // 实际部署时应使用环境变量
	jwtExpireTime := 24 * time.Hour

	// 初始化数据库服务
	dbService, err := service.NewDBService()
	if err != nil {
		log.Fatalf("Failed to initialize database service: %v", err)
	}

	// 初始化 Casbin 服务
	casbinService, err := service.NewCasbinService(dbService.DB)
	if err != nil {
		log.Fatalf("Failed to initialize casbin service: %v", err)
	}

	// 初始化各种服务
	userService := service.NewUserService(dbService.DB)
	roleService := service.NewRoleService(dbService.DB, casbinService)
	menuService := service.NewMenuService(dbService.DB)
	externalAppService := service.NewExternalAppService(dbService.DB)
	permissionService := service.NewPermissionService(dbService.DB, casbinService, roleService)

	// 初始化外部应用相关服务
	externalAppMenuService := service.NewExternalAppMenuService(dbService.DB)
	externalAppRoleService := service.NewExternalAppRoleService(dbService.DB)
	externalAppUserService := service.NewExternalAppUserService(dbService.DB)
	externalAppPermissionService := service.NewExternalAppPermissionService(dbService.DB)

	// 初始化 JWT 配置
	jwtConfig := utils.NewJWTConfig(jwtSecret, jwtExpireTime)

	// 初始化 HTTP 处理器
	authHandler := handler.NewAuthHandler(userService, jwtConfig)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	menuHandler := handler.NewMenuHandler(menuService)
	permissionHandler := handler.NewPermissionHandler(permissionService, roleService)
	authzHandler := handler.NewAuthzHandler(casbinService, menuService)
	externalAppHandler := handler.NewExternalAppHandler(externalAppService)
	externalAPIHandler := handler.NewExternalAPIHandler(externalAppService, menuService, roleService)

	// 初始化外部应用相关处理器
	externalAppMenuHandler := handler.NewExternalAppMenuHandler(externalAppMenuService)
	externalAppRoleHandler := handler.NewExternalAppRoleHandler(externalAppRoleService)
	externalAppUserHandler := handler.NewExternalAppUserHandler(externalAppUserService)
	externalAppPermissionHandler := handler.NewExternalAppPermissionHandler(externalAppPermissionService)

	// 初始化 JWT 中间件
	jwtMiddleware := customMiddleware.NewJWTMiddleware(jwtConfig)

	// 初始化外部应用认证中间件
	externalAppAuthMiddleware := customMiddleware.NewExternalAppAuthMiddleware(externalAppService)

	// 创建 Echo 实例
	e := echo.New()

	// 配置中间件
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// 静态文件服务 - 托管前端构建产物
	e.StaticFS("/", webFS)

	// 公共路由
	public := e.Group("")
	{
		// 认证相关
		public.POST("/login", authHandler.Login)
		public.POST("/logout", authHandler.Logout)

		// 外部应用注册和令牌获取
		public.POST("/external/register", externalAPIHandler.RegisterApp)
		public.POST("/external/token", externalAPIHandler.GetAppToken)
	}

	// API 路由 - 需要 JWT 认证
	api := e.Group("/api/v1")
	api.Use(jwtMiddleware.Middleware())
	{
		// 权限检查
		api.POST("/check", authzHandler.CheckPermission)

		// 用户导航菜单
		api.GET("/user/nav", authzHandler.GetUserNav)

		// 用户管理
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// 角色管理
		roles := api.Group("/roles")
		{
			roles.POST("", roleHandler.CreateRole)
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			roles.POST("/:id/menus", roleHandler.AssignMenus)
			roles.POST("/:id/permissions", roleHandler.AssignPermissions)
			roles.PUT("/:id/permissions", roleHandler.UpdatePermissions)
		}

		// 权限管理
		permissions := api.Group("/permissions")
		{
			permissions.GET("", permissionHandler.ListPermissions)
			permissions.POST("", permissionHandler.CreatePermission)
			permissions.DELETE("", permissionHandler.DeletePermission)
			permissions.GET("/roles", permissionHandler.GetPermissionRoles)
		}

		// 菜单管理
		menus := api.Group("/menus")
		{
			menus.POST("", menuHandler.CreateMenu)
			menus.GET("", menuHandler.ListMenus)
			menus.GET("/tree", menuHandler.GetMenuTree)
			menus.GET("/:id", menuHandler.GetMenu)
			menus.PUT("/:id", menuHandler.UpdateMenu)
			menus.DELETE("/:id", menuHandler.DeleteMenu)
		}

		// 外部应用管理
		externalApps := api.Group("/external/apps")
		{
			externalApps.POST("", externalAppHandler.CreateExternalApp)
			externalApps.GET("", externalAppHandler.ListExternalApps)
			externalApps.GET("/:id", externalAppHandler.GetExternalApp)
			externalApps.PUT("/:id", externalAppHandler.UpdateExternalApp)
			externalApps.DELETE("/:id", externalAppHandler.DeleteExternalApp)
			externalApps.POST("/:id/token", externalAppHandler.GenerateAppToken)

			// 应用权限管理
			externalApps.POST("/:id/permissions", externalAppHandler.CreateAppPermission)
			externalApps.GET("/:id/permissions", externalAppHandler.GetAppPermissions)
			externalApps.PUT("/permissions/:id", externalAppHandler.UpdateAppPermission)
			externalApps.DELETE("/permissions/:id", externalAppHandler.DeleteAppPermission)
		}
	}

	// 外部API路由 - 需要外部应用认证
	externalAPI := e.Group("/external/api")
	externalAPI.Use(externalAppAuthMiddleware.Middleware())
	{
		// 菜单管理
		externalAPI.GET("/menus", externalAppMenuHandler.GetMenus)
		externalAPI.POST("/menus", externalAppMenuHandler.CreateMenu)
		externalAPI.PUT("/menus/:id", externalAppMenuHandler.UpdateMenu)
		externalAPI.DELETE("/menus/:id", externalAppMenuHandler.DeleteMenu)

		// 角色管理
		externalAPI.GET("/roles", externalAppRoleHandler.GetRoles)
		externalAPI.POST("/roles", externalAppRoleHandler.CreateRole)
		externalAPI.PUT("/roles/:id", externalAppRoleHandler.UpdateRole)
		externalAPI.DELETE("/roles/:id", externalAppRoleHandler.DeleteRole)
		externalAPI.POST("/roles/:id/menus", externalAppRoleHandler.AssignRoleMenus)

		// 用户管理
		externalAPI.GET("/users", externalAppUserHandler.GetUsers)
		externalAPI.GET("/users/:id", externalAppUserHandler.GetUser)
		externalAPI.POST("/users", externalAppUserHandler.CreateUser)
		externalAPI.PUT("/users/:id", externalAppUserHandler.UpdateUser)
		externalAPI.DELETE("/users/:id", externalAppUserHandler.DeleteUser)
		externalAPI.POST("/users/:id/roles", externalAppUserHandler.AssignUserRoles)

		// 权限管理
		externalAPI.GET("/permissions", externalAppPermissionHandler.GetPermissions)
		externalAPI.POST("/permissions", externalAppPermissionHandler.CreatePermission)
		externalAPI.PUT("/permissions/:id", externalAppPermissionHandler.UpdatePermission)
		externalAPI.DELETE("/permissions/:id", externalAppPermissionHandler.DeletePermission)
		externalAPI.POST("/roles/:id/permissions", externalAppPermissionHandler.AssignRolePermissions)

		// 权限检查
		externalAPI.POST("/check-permission", externalAppPermissionHandler.CheckUserPermission)
		externalAPI.POST("/check-menu-permission", externalAppPermissionHandler.CheckUserMenuPermission)
	}

	// 启动服务器
	serverAddr := ":8080"
	log.Printf("Server starting on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
