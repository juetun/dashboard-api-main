package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/con_impl"
)

func init() {

	app_start.HandleFunc = append(app_start.HandleFunc, func(r *gin.Engine, urlPrefix string) {
		consoleHome := con_impl.NewControllerHome()
		h := r.Group(urlPrefix + "/console")
		h.GET("/home", consoleHome.Index)
	})
}
