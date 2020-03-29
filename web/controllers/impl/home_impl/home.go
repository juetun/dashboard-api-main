package home_impl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/web/controllers/inter"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services"
)

type ControllerHome struct {
	base.ControllerBase
}

func NewControllerHome() inter.System {
	controller := &ControllerHome{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerHome) Index(c *gin.Context) {
	themes := make(map[int]interface{})
	themes[1] = 1
	srv := services.NewSystemService(&base.Context{Log: r.Log})
	system, err := srv.GetSystemList()
	if err != nil {
		r.Log.Errorln("message", "console.Home.Index", "err", err.Error())
		return
	}
	data := make(map[string]interface{})
	data["themes"] = themes
	data["system"] = system
	r.Log.Infoln("message", "console.Home.Index", "message", " Succeed to get system index ")
	r.Response(c, 0, data)
	return
}

func (r *ControllerHome) Update(c *gin.Context) {
	systemIdStr := c.Param("id")
	systemIdInt, err := strconv.Atoi(systemIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Update", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Errorln("message", "system.Update", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil)
		return
	}
	// var ss common.ConsoleSystem
	ss, ok := requestJson.(pojos.ConsoleSystem)
	if !ok {
		r.Log.Errorln("message", "system.Update", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := services.NewSystemService(&base.Context{Log: r.Log})
	err = srv.SystemUpdate(systemIdInt, ss)
	if err != nil {
		r.Log.Errorln("message", "system.Update", "error", err.Error())
		r.Response(c, 405000000, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}
