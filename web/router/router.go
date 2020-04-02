package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/common"
)

type HandleRouter func(c *gin.Engine, urlPrefix string)

var HandleFunc = make([]HandleRouter, 0)

func RunLoadRouter(c *gin.Engine) (err error) {
	appConfig := common.GetAppConfig()
	var UrlPrefix = appConfig.AppName + "/" + appConfig.AppApiVersion
	io := common.NewSystemOut().SetInfoType(common.LogLevelInfo)
	io.SystemOutPrintf("Start load route config.... %s", UrlPrefix)
	defer func() {
		io.SystemOutPrintln("Load route finished")
	}()
	for _, router := range HandleFunc {
		router(c, UrlPrefix)
	}

	return
}
