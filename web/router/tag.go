package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/impl/con_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
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
