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
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/web/controllers"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services"
)

type ControllerCategory struct {
	base.ControllerBase
}

func NewControllerCategory() controllers.Console {
	controller := &ControllerCategory{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerCategory) Index(c *gin.Context) {

	srv := services.NewCategoryService()
	cates, err := srv.CateListBySort()
	if err != nil {
		r.Log.Errorln("message", "console.Cate.Index", "err", err.Error())
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
		r.Log.Errorln("message", "Cate.Store", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil)
		return
	}
	var cs pojos.CateStore
	cs, ok := requestJson.(pojos.CateStore)
	if !ok {
		r.Log.Errorln("message", "Cate.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := services.NewCategoryService()
	_, err := srv.CateStore(cs)
	if err != nil {
		r.Log.Errorln("message", "console.Cate.Store", "err", err.Error())
		r.Response(c, 402000010, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}

func (r *ControllerCategory) Edit(c *gin.Context) {
	cateIdStr := c.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Cate.Edit", "err", err.Error())
		r.Response(c, 400001002, nil)
		return
	}
	srv := services.NewCategoryService()
	cateData, err := srv.GetCateById(cateIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Cate.Edit", "err", err.Error())
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
		r.Log.Errorln("message", "console.Cate.Update", "err", err.Error())
		r.Response(c, 400001002, nil)
		return
	}
	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Errorln("message", "Cate.Update", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil)
		return
	}
	var cs pojos.CateStore
	cs, ok := requestJson.(pojos.CateStore)
	if !ok {
		r.Log.Errorln("message", "cate.Update", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := services.NewCategoryService()
	_, err = srv.CateUpdate(cateIdInt, cs)
	if err != nil {
		r.Log.Errorln("message", "cate.Update", "error", err.Error())
		r.Response(c, 402000009, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}

func (r *ControllerCategory) Destroy(c *gin.Context) {
	cateIdStr := c.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Cate.Destroy", "err", err.Error())
		r.Response(c, 400001002, nil)
		return
	}
	srv := services.NewCategoryService()
	_, err = srv.GetCateById(cateIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Cate.Destroy", "err", err.Error())
		r.Response(c, 402000000, nil)
		return
	}

	pd, err := srv.GetCateByParentId(cateIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Cate.Destroy", "err", err.Error())
		r.Response(c, 402000000, nil)
		return
	}
	if pd.Id > 0 {
		r.Log.Errorln("message", "console.Cate.Destroy", err, "It has children node")
		r.Response(c, 402000011, nil)
		return
	}

	srv.DelCateRel(cateIdInt)
	r.Response(c, 0, nil)
	return
}
