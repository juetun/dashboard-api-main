/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-06
 * Time: 23:33
 */
package con_impl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/controllers/inter"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services"
)

type ControllerLink struct {
	base.ControllerBase
}

func NewControllerLink() inter.Console {
	controller := &ControllerLink{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerLink) Index(c *gin.Context) {

	pager := base.NewPager()
	limit, offset := pager.InitPageBy(c, "GET")

	srv := services.NewLinkService(&base.Context{Log: r.Log})
	links, cnt, err := srv.LinkList(offset, limit)
	if err != nil {
		r.Log.Errorln("message", "console.Link.Index", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = links
	data["page"] = common.MyPaginate(cnt, limit, pager.PageNo)

	r.Response(c, 0, data)
	return
}
func (r *ControllerLink) Create(c *gin.Context) {
}
func (r *ControllerLink) Store(c *gin.Context) {
	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Errorln("message", "link.Store", "error", "get request_params from context fail")
		r.Response(c, 401000004, nil)
		return
	}
	ls, ok := requestJson.(pojos.LinkStore)
	if !ok {
		r.Log.Errorln("message", "link.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := services.NewLinkService(&base.Context{Log: r.Log})
	err := srv.LinkSore(ls)
	if err != nil {
		r.Log.Errorln("message", "link.Store", "error", err.Error())
		r.Response(c, 406000005, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}
func (r *ControllerLink) Edit(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Link.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := services.NewLinkService(&base.Context{Log: r.Log})
	link, err := srv.LinkDetail(linkIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Link.Edit", "err", err.Error())
		r.Response(c, 406000006, nil)
		return
	}
	r.Response(c, 0, link)
	return
}
func (r *ControllerLink) Update(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Link.Update", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Errorln("message", "Link.Update", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil)
		return
	}
	ls, ok := requestJson.(pojos.LinkStore)
	if !ok {
		r.Log.Errorln("message", "Link.Update", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := services.NewLinkService(&base.Context{Log: r.Log})
	err = srv.LinkUpdate(ls, linkIdInt)
	if err != nil {
		r.Log.Errorln("message", "Link.Update", "error", err.Error())
		r.Response(c, 406000007, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}
func (r *ControllerLink) Destroy(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Link.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := services.NewLinkService(&base.Context{Log: r.Log})
	err = srv.LinkDestroy(linkIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Link.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}
