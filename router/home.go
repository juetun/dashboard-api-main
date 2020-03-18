package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/middlewares"
	"github.com/juetun/app-dashboard/web/controllers/statistics_impl"
)

func init()  {
	HandleFunc=append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		consoleHome := statistics_impl.NewControllerHome()
		h := r.Group(urlPrefix + "/console/home")
		{
			h.GET("/", middlewares.Permission("console.home.index"), consoleHome.Index)
		}
	})
}
