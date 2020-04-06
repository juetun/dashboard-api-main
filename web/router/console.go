package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/impl/auth_impl"
	"github.com/juetun/base-wrapper/lib/app_start"
)

func init() {
	app_start.HandleFunc = append(app_start.HandleFunc,
		func(r *gin.Engine, urlPrefix string) {

			consoleAuth := auth_impl.NewControllerAuth()

			c := r.Group(urlPrefix + "/console")
			c.DELETE("/logout", consoleAuth.Logout)
			c.DELETE("/cache",consoleAuth.DelCache)
		},
	)
}
