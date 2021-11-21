// Package admins
/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:05 上午
 */
package admins

import (
	"github.com/gin-gonic/gin"
)

type Permit interface {
	// AdminUser 用户
	AdminUser(c *gin.Context)

	// AdminUserAdd 用户添加
	AdminUserAdd(c *gin.Context)

	// AdminUserDelete 用户删除
	AdminUserDelete(c *gin.Context)

	// AdminUserGroupRelease 用户所属用户权限组解除
	AdminUserGroupRelease(c *gin.Context)

	// AdminUserGroupAdd 用户所属权限组添加
	AdminUserGroupAdd(c *gin.Context)

	// AdminGroup 用户组
	AdminGroup(c *gin.Context)

	// AdminSetPermit 用户组设置权限
	AdminSetPermit(c *gin.Context)

	// AdminGroupEdit 用户组编辑
	AdminGroupEdit(c *gin.Context)

	// AdminGroupDelete 删除用户组
	AdminGroupDelete(c *gin.Context)

	// AdminMenu 菜单
	AdminMenu(c *gin.Context)

	AdminMenuWithCheck(c *gin.Context)

	// GetMenu 获取菜单信息
	GetMenu(c *gin.Context)

	// AdminMenuSearch 菜单搜索
	AdminMenuSearch(c *gin.Context)

	MenuAdd(c *gin.Context)

	MenuDelete(c *gin.Context)

	GetImport(c *gin.Context)

	MenuImport(c *gin.Context)

	MenuImportSet(c *gin.Context)

	ImportList(c *gin.Context)

	UpdateImportValue(c *gin.Context)

	EditImport(c *gin.Context)

	DeleteImport(c *gin.Context)

	MenuSave(c *gin.Context)

	GetImportByMenuId(c *gin.Context)

	// Menu 权限菜单
	Menu(c *gin.Context)

	// GetAppConfig 获取应用配置
	GetAppConfig(c *gin.Context)
}
