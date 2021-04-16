/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:05 上午
 */
package cons

import (
	"github.com/gin-gonic/gin"
)

type Permit interface {
	// 用户
	AdminUser(c *gin.Context)
	// 用户添加
	AdminUserAdd(c *gin.Context)
	// 用户删除
	AdminUserDelete(c *gin.Context)
	// 用户所属用户权限组解除
	AdminUserGroupRelease(c *gin.Context)

	// 用户所属权限组添加
	AdminUserGroupAdd(c *gin.Context)

	// 用户组
	AdminGroup(c *gin.Context)

	// 用户组编辑
	AdminGroupEdit(c *gin.Context)

	// 删除用户组
	AdminGroupDelete(c *gin.Context)

	// 菜单
	AdminMenu(c *gin.Context)

	// 获取菜单信息
	GetMenu(c *gin.Context)

	// 菜单搜索
	AdminMenuSearch(c *gin.Context)

	MenuAdd(c *gin.Context)

	MenuDelete(c *gin.Context)

	GetImport(c *gin.Context)

	MenuSave(c *gin.Context)

	// 权限菜单
	Menu(c *gin.Context)

	// 请求是否有权限
	Flag(c *gin.Context)
}
