package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/middlewares"
	"github.com/juetun/app-dashboard/web/controllers/auth_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		consoleAuth := auth_impl.NewControllerAuth()
		// 不需要登录状态权限
		al := r.Group(urlPrefix + "/console/login")
		{
			authLoginV := validate.NewValidate().NewAuthLoginV.MyValidate()
			al.GET("/", middlewares.Permission("console.login.index"), consoleAuth.Login)
			al.POST("/", middlewares.Permission("console.login.store"), authLoginV, consoleAuth.AuthLogin)
		}
		ar := r.Group(urlPrefix + "/console/register")
		{
			authRegisterV := validate.NewValidate().NewAuthRegister.MyValidate()
			ar.GET("/", middlewares.Permission("console.register.index"), consoleAuth.Register)
			ar.POST("/", middlewares.Permission("console.register.store"), authRegisterV, consoleAuth.AuthRegister)
		}
	})
}
