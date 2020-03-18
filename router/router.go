package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/common"
)

type HandleRouter func(c *gin.Engine, urlPrefix string)

var HandleFunc = make([]HandleRouter, 0)
var UrlPrefix = "/v1/app"

func RunLoadRouter(c *gin.Engine) (err error) {
	io := common.NewSystemOut().SetInfoType(common.LogLevelInfo)
	io.SystemOutPrintln("Start load route config....")
	defer func() {
		io.SystemOutPrintln("Load route finished")
	}()
	for _, router := range HandleFunc {
		router(c, UrlPrefix)
	}

	return
}
