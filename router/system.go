package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/middlewares"
	"github.com/juetun/app-dashboard/web/controllers/statistics_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		system := r.Group(urlPrefix + "/console/system")
		cSystem := statistics_impl.NewControllerHome()
		systemV := validate.NewValidate().NewSystemV.MyValidate()
		system.GET("/", middlewares.Permission("console.systemiddlewares.index"), cSystem.Index)
		system.PUT("/:id", middlewares.Permission("console.systemiddlewares.update"), systemV, cSystem.Update)
	})
}
