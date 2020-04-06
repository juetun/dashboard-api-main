package con_impl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/utils"
	"github.com/juetun/app-dashboard/web/controllers/inter"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services"
)

type ControllerTag struct {
	base.ControllerBase
}

func NewControllerTag() inter.Console {
	controller := &ControllerTag{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerTag) Index(c *gin.Context) {

	pager := base.NewPager()
	limit, offset := pager.InitPageBy(c, "GET")

	srv := services.NewTagService(&base.Context{Log: r.Log})
	count, tags, err := srv.TagsIndex(limit, offset)
	if err != nil {
		r.Log.Errorln("message", "console.Tag.Index", "err", err.Error())
		r.Response(c, 402000001, nil, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["list"] = tags
	data["page"] = utils.MyPaginate(count, limit, pager.PageNo)

	r.Response(c, 0, data)
	return
}

func (r *ControllerTag) Create(c *gin.Context) {

}

func (r *ControllerTag) Store(c *gin.Context) {
	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Errorln("message", "Tag.Store", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil, "请求参数异常")
		return
	}
	var ts pojos.TagStore
	ts, ok := requestJson.(pojos.TagStore)
	if !ok {
		r.Log.Errorln("message", "Tag.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil, "请求参数异常")
		return
	}
	srv := services.NewTagService(&base.Context{Log: r.Log})
	err := srv.TagStore(ts)
	if err != nil {
		r.Log.Errorln("message", "console.Cate.Store", "err", err.Error())
		r.Response(c, 403000006, nil, err.Error())
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerTag) Edit(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Tag.Edit", "err", err.Error())
		r.Response(c, 400001002, nil)
		return
	}
	srv := services.NewTagService(&base.Context{Log: r.Log})
	tagData, err := srv.GetTagById(tagIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Tag.Edit", "err", err.Error())
		r.Response(c, 403000008, nil, err.Error())
		return
	}
	r.Response(c, 0, tagData, "操作成功")
	return
}

func (r *ControllerTag) Update(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Tag.Update", "err", err.Error())
		r.Response(c, 400001002, nil, err.Error())
		return
	}
	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Errorln("message", "Tag.Update", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil, err.Error())
		return
	}
	ts, ok := requestJson.(pojos.TagStore)
	if !ok {
		r.Log.Errorln("message", "Tag.Update", "error", "request_params turn to error")
		r.Response(c, 400001001, nil, "您输入的参数格式不正确")
		return
	}
	srv := services.NewTagService(&base.Context{Log: r.Log})
	err = srv.TagUpdate(tagIdInt, ts)
	if err != nil {
		r.Log.Errorln("message", "Tag.Update", "error", err.Error())
		r.Response(c, 403000007, nil, err.Error())
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerTag) Destroy(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Tag.Destroy", "err", err.Error())
		r.Response(c, 400001002, nil, err.Error())
		return
	}
	srv := services.NewTagService(&base.Context{Log: r.Log})
	_, err = srv.GetTagById(tagIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Tag.Destroy", "err", err.Error())
		r.Response(c, 403000008, nil, err.Error())
		return
	}
	srv.DelTagRel(tagIdInt)
	r.Response(c, 0, nil, "操作成功")
	return
}
