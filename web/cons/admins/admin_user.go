package admins

import (
	"github.com/gin-gonic/gin"
)

type ConAdminUser interface {
	// AdminUser 用户
	AdminUser(c *gin.Context)

	// AdminUserEdit 用户添加
	AdminUserEdit(c *gin.Context)

	// AdminUserDelete 用户删除
	AdminUserDelete(c *gin.Context)

	// AdminUserUpdateWithColumn 按字段修改值
	AdminUserUpdateWithColumn(context *gin.Context)
}
