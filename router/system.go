package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/impl/home_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		system := r.Group(urlPrefix + "/console")
		cSystem := home_impl.NewControllerHome()
		systemV := validate.NewValidate().NewSystemV.MyValidate()
		system.GET("/system", cSystem.Index)
		system.PUT("/system/:id", systemV, cSystem.Update)
	})
}
