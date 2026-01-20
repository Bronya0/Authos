package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func init() {
	// Windows下必须显式注册 mime 类型，否则 js 文件可能被识别为 text/plain 导致加载失败
	if err := mime.AddExtensionType(".js", "application/javascript"); err != nil {
		log.Printf("Error registering mime type .js: %v", err)
	}
	if err := mime.AddExtensionType(".css", "text/css"); err != nil {
		log.Printf("Error registering mime type .css: %v", err)
	}
	if err := mime.AddExtensionType(".svg", "image/svg+xml"); err != nil {
		log.Printf("Error registering mime type .svg: %v", err)
	}
}

// registerStatic 使用本地 dist 目录提供静态页面服务
func registerStatic(e *echo.Echo) {
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("failed to get current working directory: %v", err)
		return
	}

	// 假设 dist 目录在当前运行目录下的 web/dist
	distPath := filepath.Join(cwd, "web", "dist")

	// 检查 dist 目录是否存在
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		log.Printf("frontend dist directory not found at: %s", distPath)
		return
	}

	// 统一处理 /authos/* 请求，支持静态文件和 SPA 路由 Fallback
	e.GET("/authos/*", func(c echo.Context) error {
		reqPath := c.Request().URL.Path
		// 去除 /authos 前缀
		relPath := strings.TrimPrefix(reqPath, "/authos")

		// 处理可能的路径前缀问题（如 /assets/assets/）
		if strings.HasPrefix(relPath, "/assets/assets/") {
			relPath = strings.TrimPrefix(relPath, "/assets/")
		}

		// 清理路径
		cleanPath := filepath.Clean(relPath)
		fullPath := filepath.Join(distPath, cleanPath)

		// 1. 尝试查找实际存在的文件
		stat, err := os.Stat(fullPath)
		if err == nil && !stat.IsDir() {
			return c.File(fullPath)
		}

		// 2. 如果文件不存在，判断是否应该是静态资源（带后缀或是 assets 目录）
		// 如果是静态资源但找不到，返回 404
		if strings.HasPrefix(relPath, "/assets/") || strings.Contains(filepath.Base(relPath), ".") {
			return echo.ErrNotFound
		}

		// 3. 其他情况（通常是前端路由路径），返回 index.html
		return c.File(filepath.Join(distPath, "index.html"))
	})

	// 处理 /authos 根路径重定向
	e.GET("/authos", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/authos/")
	})
}
