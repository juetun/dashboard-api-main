package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	con_impl2 "github.com/juetun/dashboard-api-main/web/cons/admin/impl"
	"github.com/juetun/dashboard-api-main/web/validate"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		system := r.Group(urlPrefix + "/console")
		cSystem := con_impl2.NewControllerHome()
		systemV := validate.NewValidate().NewSystemV.MyValidate()
		system.GET("/system", cSystem.Index)
		system.PUT("/system/:id", systemV, cSystem.Update)
	})
}
