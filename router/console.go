package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/impl/auth_impl"
)

func init() {
	HandleFunc = append(HandleFunc,
		func(r *gin.Engine, urlPrefix string) {

			consoleAuth := auth_impl.NewControllerAuth()

			c := r.Group(urlPrefix + "/console")
			c.DELETE("/logout", consoleAuth.Logout)
			c.DELETE("/cache",consoleAuth.DelCache)
		},
	)
}
