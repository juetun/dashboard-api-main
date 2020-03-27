package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/middlewares"
	"github.com/juetun/app-dashboard/web/controllers/auth_impl"
)

func init() {
	HandleFunc = append(HandleFunc,
		func(r *gin.Engine, urlPrefix string) {

			consoleAuth := auth_impl.NewControllerAuth()

			c := r.Group(urlPrefix + "/console")
			c.DELETE("/logout", middlewares.Permission("console.auth.logout"), consoleAuth.Logout)
			c.DELETE("/cache", middlewares.Permission("console.auth.cache"), consoleAuth.DelCache)
		},
	)
}

// func recoverHandler(c *gin.Context, err interface{}) {
// 	apiG := api.Gin{C: c}
// 	apiG.Response(http.StatusOK, 400000000, []string{})
// 	return
// }
