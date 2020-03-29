package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/impl/con_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		tag := r.Group(urlPrefix + "/console/tag")
		consoleTag := con_impl.NewControllerTag()
		tagV := validate.NewValidate().NewTagV.MyValidate()
		tag.GET("/", consoleTag.Index)
		tag.POST("/", tagV, consoleTag.Store)
		tag.GET("/edit/:id", consoleTag.Edit)
		tag.PUT("/:id", tagV, consoleTag.Update)
		tag.DELETE("/:id", consoleTag.Destroy)
	})
}
