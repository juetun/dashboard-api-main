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

	MenuSave(c *gin.Context)

	// Menu 权限菜单
	Menu(c *gin.Context)

	// GetAppConfig 获取应用配置
	GetAppConfig(c *gin.Context)
}
