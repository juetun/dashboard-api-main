package impl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/admin"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ControllerHome struct {
	base.ControllerBase
}

func NewControllerHome() admin.System {
	controller := &ControllerHome{}
	controller.ControllerBase.Init()
	return controller
}

type Themes struct {
	Key   int    `json:"key"`
	Label string `json:"label"`
}

func (r *ControllerHome) Index(c *gin.Context) {
	themes := []Themes{
		{
			Key:   1,
			Label: "默认模板",
		},
	}

	srv := srv_impl.NewSystemService(base.CreateContext(&r.ControllerBase, c))
	system, err := srv.GetSystemList()
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Home.Index", "err", err.Error())
		return
	}
	data := make(map[string]interface{})
	data["themes"] = themes
	data["system"] = system
	r.Log.Logger.Infoln("message", "console.Home.Index", "message", " Succeed to get system index ")
	r.Response(c, 0, data)
	return
}

func (r *ControllerHome) Update(c *gin.Context) {
	systemIdStr := c.Param("id")
	systemIdInt, err := strconv.Atoi(systemIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Update", "err", err.Error())
		r.Response(c, 500000000, nil, "数据ID格式不正确")
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Logger.Errorln("message", "system.Update", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil, "参数异常")
		return
	}
	// var ss common.ConsoleSystem
	ss, ok := requestJson.(wrappers.ConsoleSystem)
	if !ok {
		r.Log.Logger.Errorln("message", "system.Update", "error", "request_params turn to error")
		r.Response(c, 400001001, nil, "参数异常")
		return
	}
	srv := srv_impl.NewSystemService(base.CreateContext(&r.ControllerBase, c))
	err = srv.SystemUpdate(systemIdInt, ss)
	if err != nil {
		r.Log.Logger.Errorln("message", "system.Update", "error", err.Error())
		r.Response(c, 405000000, nil, err.Error())
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}
