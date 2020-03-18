package base

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/app_log"
)

type ControllerBase struct {
	Log *app_log.AppLog
}

func (r *ControllerBase) Init() *ControllerBase {
	r.Log = app_log.GetLog()

	return r
}

func (r *ControllerBase) Response(c *gin.Context, code int, data interface{}) {

}
