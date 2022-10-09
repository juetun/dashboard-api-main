// Package admin
/**
* @Author:ChangJiang
* @Description:
* @File:service
* @Version: 1.0.0_expanded
* @Date 2021/5/22 5:10 下午
 */
package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/admins/admin_impl"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		cons := admin_impl.NewConServiceImpl()
		rt := r.Group(urlPrefix + "/service")
		rt.POST("/list", cons.List)        // 服务列表
		rt.GET("/detail/:id", cons.Detail) // 服务列表
		rt.POST("/edit", cons.Edit)        // 服务编辑
	})
}
