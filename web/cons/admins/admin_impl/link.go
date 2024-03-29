// Package admin_impl
/**
 * Created by GoLand.
 * Date: 2019-05-06
 * Time: 23:33
 */
package admin_impl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ControllerLink struct {
	base.ControllerBase
}

func NewControllerLink() admins.Console {
	controller := &ControllerLink{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerLink) Index(c *gin.Context) {
	var arg response.PageQuery

	var err error

	if arg.PageNo, err = strconv.Atoi(c.DefaultQuery("pages", strconv.Itoa(response.DefaultPageNo))); err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	arg.PageSize, err = strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(response.DefaultPageSize)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	arg.DefaultPage()
	pager := response.NewPager(response.PagerBaseQuery(&arg))
	srv := srv_impl.NewLinkService(base.CreateContext(&r.ControllerBase, c))
	pager.List, pager.TotalCount, err = srv.LinkList(arg.GetOffset(), pager.PageSize)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Link.Index2", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}

	r.Response(c, 0, pager)
	return
}
func (r *ControllerLink) Create(c *gin.Context) {
}
func (r *ControllerLink) Store(c *gin.Context) {
	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Logger.Errorln("message", "link.Store", "error", "get request_params from context fail")
		r.Response(c, 401000004, nil)
		return
	}
	ls, ok := requestJson.(wrappers.LinkStore)
	if !ok {
		r.Log.Logger.Errorln("message", "link.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := srv_impl.NewLinkService(base.CreateContext(&r.ControllerBase, c))
	err := srv.LinkSore(ls)
	if err != nil {
		r.Log.Logger.Errorln("message", "link.Store", "error", err.Error())
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
		r.Log.Logger.Errorln("message", "console.Link.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := srv_impl.NewLinkService(base.CreateContext(&r.ControllerBase, c))
	link, err := srv.LinkDetail(linkIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Link.Edit", "err", err.Error())
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
		r.Log.Logger.Errorln("message", "console.Link.Update", "err", err.Error())
		r.Response(c, 500000000, "链接数据格式异常")
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Logger.Errorln("message", "Link.Update", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil, "您输入的数据格式错误")
		return
	}
	ls, ok := requestJson.(wrappers.LinkStore)
	if !ok {
		r.Log.Logger.Errorln("message", "Link.Update", "error", "request_params turn to error")
		r.Response(c, 400001001, nil, "您输入的数据格式错误")
		return
	}
	srv := srv_impl.NewLinkService(base.CreateContext(&r.ControllerBase, c))
	err = srv.LinkUpdate(ls, linkIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "Link.Update", "error", err.Error())
		r.Response(c, 406000007, nil, err.Error())
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}
func (r *ControllerLink) Destroy(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Link.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	srv := srv_impl.NewLinkService(base.CreateContext(&r.ControllerBase, c))
	err = srv.LinkDestroy(linkIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Link.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}
