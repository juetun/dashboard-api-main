package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons_admin/con_impl"
	"github.com/juetun/dashboard-api-main/web/validate"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		tag := r.Group(urlPrefix + "/console")
		consoleTag := con_impl.NewControllerTag()
		tagV := validate.NewValidate().NewTagV.MyValidate()
		tag.GET("/tag", consoleTag.Index)
		tag.POST("/tag", tagV, consoleTag.Store)
		tag.GET("/tag/edit/:id", consoleTag.Edit)
		tag.PUT("/tag/:id", tagV, consoleTag.Update)
		tag.DELETE("/tag/:id", consoleTag.Destroy)
	})
}
