package main

import (
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

// 跨域中间件配置
func corsConfig() echo.MiddlewareFunc {
	return echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodPatch,
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"X-System-Token",
			"X-App-Token",
			"X-App-ID",
		},
		AllowCredentials: true,
	})
}

// 预留代码：使用 embed.FS 托管前端构建产物
// 假设前端构建产物在 web/dist 目录下
//
// var webFS embed.FS

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
	apiPermissionService := service.NewApiPermissionService(dbService.DB, casbinService, roleService)
	applicationService := service.NewApplicationService(dbService.DB)
	auditLogService := service.NewAuditLogService(dbService.DB)

	// 初始化 JWT 配置
	jwtConfig := utils.NewJWTConfig(jwtSecret, jwtExpireTime)

	// 初始化 HTTP 处理器
	authHandler := handler.NewAuthHandler(userService, applicationService, auditLogService, jwtConfig)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	menuHandler := handler.NewMenuHandler(menuService)
	apiPermissionHandler := handler.NewApiPermissionHandler(apiPermissionService)
	applicationHandler := handler.NewApplicationHandler(applicationService)
	auditLogHandler := handler.NewAuditLogHandler(auditLogService)
	authzHandler := handler.NewAuthzHandler(casbinService, menuService)

	// 初始化 JWT 中间件
	jwtMiddleware := customMiddleware.NewJWTMiddleware(jwtConfig)

	// 创建 Echo 实例
	e := echo.New()

	// 配置中间件
	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogHost:   true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			log.Printf("%s %s %s %d", v.Host, v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(echoMiddleware.Recover())
	e.Use(corsConfig())

	// 静态文件服务 - 托管前端构建产物
	// e.StaticFS("/", webFS)

	// 公共路由
	public := e.Group("/api/public")
	{
		// 认证相关
		public.POST("/login", authHandler.Login)
		public.POST("/system-login", authHandler.SystemLogin)
		public.POST("/app-login", authHandler.AppLogin)
		public.POST("/logout", authHandler.Logout)

	}

	// API 路由 - 需要 JWT 认证
	api := e.Group("/api/v1")
	api.Use(jwtMiddleware.Middleware())
	{
		// 应用相关
		api.GET("/applications", applicationHandler.ListApplications)
		api.POST("/applications", applicationHandler.CreateApplication)
		api.GET("/applications/by-code/:code", applicationHandler.GetApplicationByCode)
		api.GET("/applications/:id", applicationHandler.GetApplication)
		api.PUT("/applications/:id", applicationHandler.UpdateApplication)
		api.DELETE("/applications/:id", applicationHandler.DeleteApplication)

		// 权限检查
		api.POST("/check", authzHandler.CheckPermission)
		api.POST("/auth/check-permission", authzHandler.CheckPermissionByKey)

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

		// 仪表盘统计
		api.GET("/dashboard/stats", authHandler.GetDashboardStats)

		// 角色管理
		roles := api.Group("/roles")
		{
			roles.POST("", roleHandler.CreateRole)
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			roles.GET("/:id/menus", roleHandler.GetRoleMenus)
			roles.POST("/:id/menus", roleHandler.AssignMenus)
			roles.PUT("/:id/menus", roleHandler.UpdateRoleMenus)
			roles.POST("/:id/permissions", roleHandler.AssignPermissions)
			roles.PUT("/:id/permissions", roleHandler.UpdatePermissions)
		}

		// 接口权限管理
		apiPermissions := api.Group("/api-permissions")
		{
			apiPermissions.GET("", apiPermissionHandler.ListApiPermissions)
			apiPermissions.POST("", apiPermissionHandler.CreateApiPermission)
			apiPermissions.GET("/:id", apiPermissionHandler.GetApiPermission)
			apiPermissions.PUT("/:id", apiPermissionHandler.UpdateApiPermission)
			apiPermissions.DELETE("/:id", apiPermissionHandler.DeleteApiPermission)
			apiPermissions.GET("/roles/:roleUUID", apiPermissionHandler.GetApiPermissionsForRole)
			apiPermissions.POST("/roles/:roleUUID", apiPermissionHandler.AddApiPermissionToRole)
			apiPermissions.DELETE("/roles/:roleUUID", apiPermissionHandler.RemoveApiPermissionFromRole)
		}

		// 菜单管理
		menus := api.Group("/menus")
		{
			menus.POST("", menuHandler.CreateMenu)
			menus.GET("", menuHandler.ListMenus)
			menus.GET("/tree", menuHandler.GetMenuTree)
			menus.GET("/non-system-tree", menuHandler.GetNonSystemMenuTree)
			menus.GET("/:id", menuHandler.GetMenu)
			menus.PUT("/:id", menuHandler.UpdateMenu)
			menus.DELETE("/:id", menuHandler.DeleteMenu)
		}

		// 审计日志
		api.GET("/audit-logs", auditLogHandler.ListAuditLogs)
		api.GET("/system/audit-logs", auditLogHandler.ListSystemAuditLogs)
	}

	// 启动服务器
	serverAddr := ":8080"
	log.Printf("Server starting on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
