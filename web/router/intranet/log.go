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
		con := intranet_impl.NewConLogIntranet()

		path := r.Group(urlPrefix)

		path.POST("/add_log", con.AddLog)

	})
}
