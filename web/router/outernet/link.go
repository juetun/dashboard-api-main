package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons_outernet/con_impl"
	"github.com/juetun/dashboard-api-main/web/validate"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet, func(r *gin.Engine, urlPrefix string) {
		link := r.Group(urlPrefix + "/console")
		consoleLink := con_impl.NewControllerLink()
		linkV := validate.NewValidate().NewLinkV.MyValidate()
		link.GET("/link", consoleLink.Index)
		link.POST("/link", linkV, consoleLink.Store)
		link.GET("/link/edit/:id", consoleLink.Edit)
		link.PUT("/link/:id", linkV, consoleLink.Update)
		link.DELETE("/link/:id", consoleLink.Destroy)
	}, )
}
