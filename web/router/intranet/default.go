// Package intranet
// /**
package intranet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
)

func init() {
	app_start.HandleFuncIntranet = append(app_start.HandleFuncIntranet, func(r *gin.Engine, urlPrefix string) {
	})
}
