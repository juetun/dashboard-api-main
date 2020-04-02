package router

import (
	"github.com/gin-gonic/gin"
	console "github.com/juetun/app-dashboard/web/controllers/impl/con_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		link := r.Group(urlPrefix + "/console")
		consoleLink := console.NewControllerLink()
		linkV := validate.NewValidate().NewLinkV.MyValidate()
		link.GET("/link", consoleLink.Index)
		link.POST("/link", linkV, consoleLink.Store)
		link.GET("/link/edit/:id", consoleLink.Edit)
		link.PUT("/link/:id", linkV, consoleLink.Update)
		link.DELETE("/link/:id", consoleLink.Destroy)
	}, )
}
