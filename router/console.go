package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/study/app-dashboard/gin/api"
	"github.com/juetun/study/app-dashboard/lib/middlewares"
	"github.com/juetun/study/app-dashboard/web/controllers/auth_impl"
	"net/http"
)

func init() {
	HandleFunc = append(HandleFunc,
		func(r *gin.Engine, urlPrefix string) {
			c := r.Group(urlPrefix + "/console")
			consoleAuth := auth_impl.NewControllerAuth()
			c.DELETE("/logout", middlewares.Permission("console.auth.logout"), consoleAuth.Logout)
			c.DELETE("/cache", middlewares.Permission("console.auth.cache"), consoleAuth.DelCache)
		},
	)
}
func recoverHandler(c *gin.Context, err interface{}) {
	apiG := api.Gin{C: c}
	apiG.Response(http.StatusOK, 400000000, []string{})
	return
}
