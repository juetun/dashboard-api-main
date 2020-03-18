package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/middlewares"
	console "github.com/juetun/app-dashboard/web/controllers/con_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		link := r.Group(urlPrefix + "/console/link")

		consoleLink := console.NewControllerLink()
		linkV := validate.NewValidate().NewLinkV.MyValidate()
		link.GET("/", middlewares.Permission("console.link.index"), consoleLink.Index)
		link.POST("/", middlewares.Permission("console.link.store"), linkV, consoleLink.Store)
		link.GET("/edit/:id", middlewares.Permission("console.link.edit"), consoleLink.Edit)
		link.PUT("/:id", middlewares.Permission("console.link.update"), linkV, consoleLink.Update)
		link.DELETE("/:id", middlewares.Permission("console.link.destroy"), consoleLink.Destroy)
	}, )
}
