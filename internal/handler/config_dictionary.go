package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"Authos/internal/model"
	"Authos/internal/service"

	"github.com/labstack/echo/v4"
)

type ConfigDictionaryHandler struct {
	ConfigDictionaryService *service.ConfigDictionaryService
}

func NewConfigDictionaryHandler(configDictionaryService *service.ConfigDictionaryService) *ConfigDictionaryHandler {
	return &ConfigDictionaryHandler{
		ConfigDictionaryService: configDictionaryService,
	}
}

func (h *ConfigDictionaryHandler) ListConfigDictionaries(c echo.Context) error {
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	key := c.QueryParam("key")
	db := h.ConfigDictionaryService.DB.Where("app_id = ?", appID)
	if key != "" {
		db = db.Where("key LIKE ?", "%"+key+"%")
	}

	var items []*model.ConfigDictionary
	if err := db.Order("id asc").Find(&items).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "获取配置字典列表失败"})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *ConfigDictionaryHandler) GetConfigDictionary(c echo.Context) error {
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的字典ID"})
	}

	item, err := h.ConfigDictionaryService.GetConfigDictionary(uint(id), appID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "配置字典不存在"})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *ConfigDictionaryHandler) CreateConfigDictionary(c echo.Context) error {
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Desc  string `json:"desc"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "请求参数错误"})
	}

	item, err := h.ConfigDictionaryService.CreateConfigDictionary(appID, req.Key, req.Value, req.Desc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

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

	h.ConfigDictionaryService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "CREATE",
		Resource:   "CONFIG_DICTIONARY",
		ResourceID: fmt.Sprintf("%d", item.ID),
		Content:    fmt.Sprintf("创建配置字典: %s", item.Key),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusCreated, item)
}

func (h *ConfigDictionaryHandler) UpdateConfigDictionary(c echo.Context) error {
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的字典ID"})
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Desc  string `json:"desc"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "请求参数错误"})
	}

	item, err := h.ConfigDictionaryService.UpdateConfigDictionary(uint(id), appID, req.Key, req.Value, req.Desc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

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

	h.ConfigDictionaryService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "UPDATE",
		Resource:   "CONFIG_DICTIONARY",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("更新配置字典: %s", item.Key),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, item)
}

func (h *ConfigDictionaryHandler) DeleteConfigDictionary(c echo.Context) error {
	appID, err := getAppIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "获取应用ID失败"})
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "无效的字典ID"})
	}

	if err := h.ConfigDictionaryService.DeleteConfigDictionary(uint(id), appID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

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

	h.ConfigDictionaryService.DB.Create(&model.AuditLog{
		AppID:      appID,
		UserID:     userID,
		Username:   username,
		Action:     "DELETE",
		Resource:   "CONFIG_DICTIONARY",
		ResourceID: fmt.Sprintf("%d", id),
		Content:    fmt.Sprintf("删除配置字典ID: %d", id),
		IP:         c.RealIP(),
		Status:     1,
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "配置字典删除成功"})
}
