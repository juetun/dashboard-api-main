package router

import (
	"github.com/gin-gonic/gin"
	console "github.com/juetun/app-dashboard/web/controllers/impl/con_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		link := r.Group(urlPrefix + "/console/link")

		consoleLink := console.NewControllerLink()
		linkV := validate.NewValidate().NewLinkV.MyValidate()
		link.GET("/", consoleLink.Index)
		link.POST("/", linkV, consoleLink.Store)
		link.GET("/edit/:id", consoleLink.Edit)
		link.PUT("/:id", linkV, consoleLink.Update)
		link.DELETE("/:id", consoleLink.Destroy)
	}, )
}
