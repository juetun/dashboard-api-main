package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/middlewares"
	"github.com/juetun/app-dashboard/web/controllers/con_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		tag := r.Group(urlPrefix + "/console/tag")
		consoleTag := con_impl.NewControllerTag()
		tagV := validate.NewValidate().NewTagV.MyValidate()
		tag.GET("/", middlewares.Permission("console.tag.index"), consoleTag.Index)
		tag.POST("/", middlewares.Permission("console.tag.store"), tagV, consoleTag.Store)
		tag.GET("/edit/:id", middlewares.Permission("console.tag.edit"), consoleTag.Edit)
		tag.PUT("/:id", middlewares.Permission("console.tag.update"), tagV, consoleTag.Update)
		tag.DELETE("/:id", middlewares.Permission("console.tag.destroy"), consoleTag.Destroy)
	})
}
