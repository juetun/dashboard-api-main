// Package outernet
/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:07 上午
 */
package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	con_impl2 "github.com/juetun/dashboard-api-main/web/cons/admin/impl"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		rou := r.Group(urlPrefix + "/export")
		impl := con_impl2.NewControllerExportData()

		rou.GET("/list", impl.List)
		rou.POST("/init", impl.Init)
		rou.GET("/progress", impl.Progress)
		rou.GET("/cancel", impl.Cancel)
	})
}
