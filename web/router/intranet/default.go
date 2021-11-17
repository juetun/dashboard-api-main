// Package intranet
// /**
package intranet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/intranets/intranet_impl"
)

func init() {
	app_start.HandleFuncIntranet = append(app_start.HandleFuncIntranet, func(r *gin.Engine, urlPrefix string) {
		con := intranet_impl.NewConPermitIntranet()

		path := r.Group(urlPrefix)
		path.GET("/get_import_permit", con.GetImportPermit) // 获取接口权限

		path.GET("/get_have_permit", con.GetUerImportPermit) // 判断用户是否有接口权限
	})
}
