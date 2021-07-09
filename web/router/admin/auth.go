package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons_admin/con_impl"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		consoleHome := con_impl.NewControllerHome()
		h := r.Group(urlPrefix + "/console")
		h.GET("/home", consoleHome.Index)
	})
}
