package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons_outernet/con_impl"
)

func init() {

	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet, func(r *gin.Engine, urlPrefix string) {
		consoleHome := con_impl.NewControllerHome()
		h := r.Group(urlPrefix + "/console")
		h.GET("/home", consoleHome.Index)
	})
}
