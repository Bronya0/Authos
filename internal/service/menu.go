package service

import (
	"gorm.io/gorm"

	"Authos/internal/model"
)

// MenuService 菜单服务
type MenuService struct {
	DB *gorm.DB
}

// NewMenuService 创建菜单服务实例
func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{DB: db}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	return s.DB.Create(menu).Error
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(menu *model.Menu) error {
	return s.DB.Updates(menu).Error
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	// 开始事务
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 删除菜单
		if err := tx.Delete(&model.Menu{}, id).Error; err != nil {
			return err
		}

		// 删除子菜单
		if err := tx.Where("parent_id = ?", id).Delete(&model.Menu{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	if err := s.DB.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// ListMenus 列出所有菜单（扁平结构）
func (s *MenuService) ListMenus() ([]*model.Menu, error) {
	var menus []*model.Menu
	if err := s.DB.Order("sort asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetMenuTree 获取菜单树
func (s *MenuService) GetMenuTree() ([]*model.Menu, error) {
	// 获取所有菜单
	menus, err := s.ListMenus()
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return buildMenuTree(menus, 0), nil
}

// GetUserMenuTree 根据用户ID获取用户有权访问的菜单树
func (s *MenuService) GetUserMenuTree(userID uint) ([]*model.Menu, error) {
	// 获取用户关联的角色
	var user model.User
	if err := s.DB.Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 提取用户有权访问的菜单ID
	menuIDMap := make(map[uint]bool)
	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			menuIDMap[menu.ID] = true
		}
	}

	// 获取所有菜单
	allMenus, err := s.ListMenus()
	if err != nil {
		return nil, err
	}

	// 过滤用户有权访问的菜单
	var userMenus []*model.Menu
	for _, menu := range allMenus {
		if menuIDMap[menu.ID] {
			userMenus = append(userMenus, menu)
		}
	}

	// 构建菜单树
	return buildMenuTree(userMenus, 0), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []*model.Menu, parentID uint) []*model.Menu {
	var tree []*model.Menu

	for _, menu := range menus {
		if menu.ParentID == parentID {
			children := buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				menu.Children = children
			}
			tree = append(tree, menu)
		}
	}

	return tree
}
