/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:05 上午
 */
package inter

import (
	"github.com/gin-gonic/gin"
)

type Permit interface {
	// 用户
	AdminUser(c *gin.Context)

	// 用户组
	AdminGroup(c *gin.Context)

	// 菜单
	AdminMenu(c *gin.Context)

	MenuAdd(c *gin.Context)

	MenuDelete(c *gin.Context)

	MenuSave(c *gin.Context)
	// 权限菜单
	Menu(c *gin.Context)

	// 请求是否有权限
	Flag(c *gin.Context)
}
