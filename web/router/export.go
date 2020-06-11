/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:07 上午
 */
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_start"
	"github.com/juetun/dashboard-api-main/web/controllers/impl/export_impl"
)

func init() {
	app_start.HandleFunc = append(app_start.HandleFunc, func(r *gin.Engine, urlPrefix string) {
		rou := r.Group(urlPrefix + "/export")
		impl := export_impl.NewControllerExportData()

		rou.GET("/list", impl.List)
		rou.POST("/init", impl.Init)
		rou.GET("/progress", impl.Progress)
		rou.GET("/cancel", impl.Cancel)
	})
}
