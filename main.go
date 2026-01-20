package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
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

//go:embed web-vue3/dist
var embeddedDist embed.FS

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
	configDictionaryService := service.NewConfigDictionaryService(dbService.DB)

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
	authzHandler := handler.NewAuthzHandler(casbinService, menuService, applicationService, apiPermissionService, jwtConfig)
	configDictionaryHandler := handler.NewConfigDictionaryHandler(configDictionaryService)

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

	// 静态文件服务 - 托管前端构建产物（内嵌到二进制）
	// 注意：需要先在 web-vue3 目录执行构建命令生成 dist 目录
	if distFS, err := fs.Sub(embeddedDist, "web-vue3/dist"); err != nil {
		// 前端构建资源不存在时仅提供 API 服务
		log.Printf("frontend dist not embedded, only API will be available: %v", err)
	} else {
		fileServer := http.FileServer(http.FS(distFS))

		// 通用静态资源路由，排除 /api/ 前缀，兼容 SPA 前端路由
		e.GET("/*", func(c echo.Context) error {
			path := c.Request().URL.Path
			// API 路由直接交给后面的分组处理
			if strings.HasPrefix(path, "/api/") {
				return echo.ErrNotFound
			}

			trimmed := strings.TrimPrefix(path, "/")
			if trimmed == "" {
				trimmed = "index.html"
			}

			// 如果不存在对应静态文件，则回退到 index.html（支持前端 history 路由）
			if _, err := distFS.Open(trimmed); err != nil {
				c.Request().URL.Path = "/"
			}

			fileServer.ServeHTTP(c.Response(), c.Request())
			return nil
		})
	}

	// 公共路由
	public := e.Group("/api/public")
	{
		// 认证相关
		public.POST("/login", authHandler.Login)
		public.POST("/system-login", authHandler.SystemLogin)
		public.POST("/app-login", authHandler.AppLogin)
		public.POST("/proxy-login", authHandler.ProxyLogin)                  // 新增：后端代理登录
		public.POST("/check-access", authzHandler.CheckPermissionWithSecret) // 新增：统一鉴权
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

		configDictionaries := api.Group("/config-dictionaries")
		{
			configDictionaries.POST("", configDictionaryHandler.CreateConfigDictionary)
			configDictionaries.GET("", configDictionaryHandler.ListConfigDictionaries)
			configDictionaries.GET("/:id", configDictionaryHandler.GetConfigDictionary)
			configDictionaries.PUT("/:id", configDictionaryHandler.UpdateConfigDictionary)
			configDictionaries.DELETE("/:id", configDictionaryHandler.DeleteConfigDictionary)
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
