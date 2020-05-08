package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_start"
	"github.com/juetun/dashboard-api-main/web/controllers/impl/home_impl"
)

func init() {

	app_start.HandleFunc = append(app_start.HandleFunc, func(r *gin.Engine, urlPrefix string) {
		// consoleAuth := auth_impl.NewControllerAuth()
		// // 不需要登录状态权限
		// al := r.Group(urlPrefix + "/console")
		//
		// authLoginV := validate.NewValidate().NewAuthLoginV.MyValidate()
		// al.GET("/login", consoleAuth.Login) //
		// al.POST("/login", authLoginV, consoleAuth.AuthLogin)
		//
		// authRegisterV := validate.NewValidate().NewAuthRegister.MyValidate()
		// al.GET("/register", consoleAuth.Register)
		// al.POST("/register", authRegisterV, consoleAuth.AuthRegister)

		consoleHome := home_impl.NewControllerHome()
		h := r.Group(urlPrefix + "/console")
		{
			h.GET("/home", consoleHome.Index)
		}
	})
}
