/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 5:10 下午
 */
package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons_outernet/con_impl"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet, func(r *gin.Engine, urlPrefix string) {
		cons := con_impl.NewConServiceImpl()
		rt := r.Group(urlPrefix + "/service")
		rt.POST("/list", cons.List)        // 服务列表
		rt.GET("/detail/:id", cons.Detail) // 服务列表
		rt.POST("/edit", cons.Edit)        // 服务编辑
	})
}
