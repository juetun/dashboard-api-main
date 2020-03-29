package base

import (
	"net/http"
	"strings"

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


type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"message"`
}

func (r *ControllerBase) Response(c *gin.Context, code int, data interface{}, msg ...string) {
	c.JSON(http.StatusOK, Result{Code: code, Data: data, Msg: strings.Join(msg, ",")})
}