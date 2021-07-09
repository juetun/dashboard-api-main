/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-12
 * Time: 23:06
 */
package con_impl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons_admin"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ControllerCategory struct {
	base.ControllerBase
}

func NewControllerCategory() cons_admin.Console {
	controller := &ControllerCategory{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerCategory) Index(c *gin.Context) {

	srv := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))
	cates, err := srv.CateListBySort()
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Cate.Index", "err", err.Error())
		r.Response(c, 402000001, nil)
		return
	}
	r.Response(c, 0, cates)
	return
}

func (r *ControllerCategory) Create(c *gin.Context) {

}

func (r *ControllerCategory) Store(c *gin.Context) {
	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Logger.Errorln("message", "Cate.Store", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil)
		return
	}
	var cs wrappers.CateStore
	cs, ok := requestJson.(wrappers.CateStore)
	if !ok {
		r.Log.Logger.Errorln("message", "Cate.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))
	_, err := srv.CateStore(cs)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Cate.Store", "err", err.Error())
		r.Response(c, 402000010, nil, err.Error())
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerCategory) Edit(c *gin.Context) {
	cateIdStr := c.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Cate.Edit", "err", err.Error())
		r.Response(c, 400001002, nil)
		return
	}
	srv := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))
	cateData, err := srv.GetCateById(cateIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Cate.Edit", "err", err.Error())
		r.Response(c, 402000000, nil)
		return
	}
	r.Response(c, 0, cateData)
	return
}

func (r *ControllerCategory) Update(c *gin.Context) {
	cateIdStr := c.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Cate.Update", "err", err.Error())
		r.Response(c, 400001002, nil)
		return
	}
	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Logger.Errorln("message", "Cate.Update", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil)
		return
	}
	var cs wrappers.CateStore
	cs, ok := requestJson.(wrappers.CateStore)
	if !ok {
		r.Log.Logger.Errorln("message", "cate.Update", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))
	_, err = srv.CateUpdate(cateIdInt, cs)
	if err != nil {
		r.Log.Logger.Errorln("message", "cate.Update", "error", err.Error())
		r.Response(c, 402000009, nil,err.Error())
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerCategory) Destroy(c *gin.Context) {
	cateIdStr := c.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Cate.Destroy", "err", err.Error())
		r.Response(c, 400001002, nil)
		return
	}
	srv := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))
	_, err = srv.GetCateById(cateIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Cate.Destroy", "err", err.Error())
		r.Response(c, 402000000, nil)
		return
	}

	pd, err := srv.GetCateByParentId(cateIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Cate.Destroy", "err", err.Error())
		r.Response(c, 402000000, nil)
		return
	}
	if pd.Id > 0 {
		r.Log.Logger.Errorln("message", "console.Cate.Destroy", err, "It has children node")
		r.Response(c, 402000011, nil)
		return
	}

	srv.DelCateRel(cateIdInt)
	r.Response(c, 0, nil)
	return
}
