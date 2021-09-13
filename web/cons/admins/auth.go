package admins

import "github.com/gin-gonic/gin"

type ConsoleAuth interface {
	Register(*gin.Context)
	AuthRegister(*gin.Context)

	// Login 用户登录信息接口
	Login(*gin.Context)

	// AuthLogin 验证登录接口
	AuthLogin(*gin.Context)
	Logout(*gin.Context)
	DelCache(*gin.Context)
}
