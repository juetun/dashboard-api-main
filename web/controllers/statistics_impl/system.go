package statistics_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/study/app-content/common"
	"github.com/juetun/study/app-content/service"
	"github.com/juetun/study/app-dashboard/gin/api"
	"github.com/juetun/study/app-dashboard/lib/base"
	"github.com/juetun/study/app-dashboard/web/controllers"
	"net/http"
	"strconv"
)

type ControllerHome struct {
	base.ControllerBase
}

func NewControllerHome() controllers.Statistics {
	controller := &ControllerHome{}
	controller.ControllerBase.Init()
	return controller
}

func (s *ControllerHome) Index(c *gin.Context) {
	appG := api.Gin{C: c}
	themes := make(map[int]interface{})
	themes[1] = 1
	system, err := service.GetSystemList()
	if err != nil {
		zgh.ZLog().Error("message", "console.Home.Index", "err", err.Error())
		return
	}
	data := make(map[string]interface{})
	data["themes"] = themes
	data["system"] = system
	zgh.ZLog().Info("message", "console.Home.Index", "message", " Succeed to get system index ")
	appG.Response(http.StatusOK, 0, data)
	return
}

func (s *ControllerHome) Update(c *gin.Context) {
	systemIdStr := c.Param("id")
	systemIdInt, err := strconv.Atoi(systemIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		zgh.ZLog().Error("message", "console.Update", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		zgh.ZLog().Error("message", "system.Update", "error", "get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	//var ss common.ConsoleSystem
	ss, ok := requestJson.(common.ConsoleSystem)
	if !ok {
		zgh.ZLog().Error("message", "system.Update", "error", "request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.SystemUpdate(systemIdInt, ss)
	if err != nil {
		zgh.ZLog().Error("message", "system.Update", "error", err.Error())
		appG.Response(http.StatusOK, 405000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
