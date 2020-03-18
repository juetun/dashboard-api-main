/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-06
 * Time: 23:33
 */
package con_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/dashboard/gin/api"
	"github.com/juetun/study/app-content/common"
	"github.com/juetun/study/app-content/conf"
	"github.com/juetun/study/app-content/service"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/web/controllers"
	"net/http"
	"strconv"
)

type ControllerLink struct {
	base.ControllerBase
}

func NewControllerLink() controllers.Console {
	controller := &ControllerLink{}
	controller.ControllerBase.Init()
	return controller
}

func (l *ControllerLink) Index(c *gin.Context) {
	appG := api.Gin{C: c}

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultLimit)
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		zgh.ZLog().Error("message", "console.Post.Index", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	limit, offset := common.Offset(queryPage, queryLimit)

	links, cnt, err := service.LinkList(offset, limit)
	if err != nil {
		zgh.ZLog().Error("message", "console.Link.Index", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = links
	data["page"] = common.MyPaginate(cnt, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}
func (l *ControllerLink) Create(c *gin.Context) {
}
func (l *ControllerLink) Store(c *gin.Context) {
	appG := api.Gin{C: c}
	requestJson, exists := c.Get("json")
	if !exists {
		zgh.ZLog().Error("message", "link.Store", "error", "get request_params from context fail")
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	ls, ok := requestJson.(common.LinkStore)
	if !ok {
		zgh.ZLog().Error("message", "link.Store", "error", "request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}

	err := service.LinkSore(ls)
	if err != nil {
		zgh.ZLog().Error("message", "link.Store", "error", err.Error())
		appG.Response(http.StatusOK, 406000005, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
func (l *ControllerLink) Edit(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		zgh.ZLog().Error("message", "console.Link.Edit", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	link, err := service.LinkDetail(linkIdInt)
	if err != nil {
		zgh.ZLog().Error("message", "console.Link.Edit", "err", err.Error())
		appG.Response(http.StatusOK, 406000006, nil)
		return
	}
	appG.Response(http.StatusOK, 0, link)
	return
}
func (l *ControllerLink) Update(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		zgh.ZLog().Error("message", "console.Link.Update", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		zgh.ZLog().Error("message", "Link.Update", "error", "get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	ls, ok := requestJson.(common.LinkStore)
	if !ok {
		zgh.ZLog().Error("message", "Link.Update", "error", "request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.LinkUpdate(ls, linkIdInt)
	if err != nil {
		zgh.ZLog().Error("message", "Link.Update", "error", err.Error())
		appG.Response(http.StatusOK, 406000007, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
func (l *ControllerLink) Destroy(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		zgh.ZLog().Error("message", "console.Link.Destroy", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	err = service.LinkDestroy(linkIdInt)
	if err != nil {
		zgh.ZLog().Error("message", "console.Link.Destroy", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
