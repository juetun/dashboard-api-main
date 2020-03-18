package con_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/study/app-content/common"
	"github.com/juetun/study/app-content/conf"
	"github.com/juetun/study/app-content/service"
	"github.com/juetun/app-dashboard/gin/api"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/web/controllers"
	"net/http"
	"strconv"
)

type ControllerTag struct {
	base.ControllerBase
}

func NewControllerTag() controllers.Console {
	controller := &ControllerTag{}
	controller.ControllerBase.Init()
	return controller
}

func (t *ControllerTag) Index(c *gin.Context) {
	appG := api.Gin{C: c}

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)
	count, tags, err := service.TagsIndex(limit, offset)
	if err != nil {
		zgh.ZLog().Error("message", "console.Tag.Index", "err", err.Error())
		appG.Response(http.StatusOK, 402000001, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		zgh.ZLog().Error("message", "console.Tag.Index", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = tags
	data["page"] = common.MyPaginate(count, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (t *ControllerTag) Create(c *gin.Context) {

}

func (t *ControllerTag) Store(c *gin.Context) {
	appG := api.Gin{C: c}
	requestJson, exists := c.Get("json")
	if !exists {
		zgh.ZLog().Error("message", "Tag.Store", "error", "get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	var ts common.TagStore
	ts, ok := requestJson.(common.TagStore)
	if !ok {
		zgh.ZLog().Error("message", "Tag.Store", "error", "request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err := service.TagStore(ts)
	if err != nil {
		zgh.ZLog().Error("message", "console.Cate.Store", "err", err.Error())
		appG.Response(http.StatusOK, 403000006, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (t *ControllerTag) Edit(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		zgh.ZLog().Error("message", "console.Tag.Edit", "err", err.Error())
		appG.Response(http.StatusOK, 400001002, nil)
		return
	}
	tagData, err := service.GetTagById(tagIdInt)
	if err != nil {
		zgh.ZLog().Error("message", "console.Tag.Edit", "err", err.Error())
		appG.Response(http.StatusOK, 403000008, nil)
		return
	}
	appG.Response(http.StatusOK, 0, tagData)
	return
}

func (t *ControllerTag) Update(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		zgh.ZLog().Error("message", "console.Tag.Update", "err", err.Error())
		appG.Response(http.StatusOK, 400001002, nil)
		return
	}
	requestJson, exists := c.Get("json")
	if !exists {
		zgh.ZLog().Error("message", "Tag.Update", "error", "get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	ts, ok := requestJson.(common.TagStore)
	if !ok {
		zgh.ZLog().Error("message", "Tag.Update", "error", "request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.TagUpdate(tagIdInt, ts)
	if err != nil {
		zgh.ZLog().Error("message", "Tag.Update", "error", err.Error())
		appG.Response(http.StatusOK, 403000007, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (t *ControllerTag) Destroy(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		zgh.ZLog().Error("message", "console.Tag.Destroy", "err", err.Error())
		appG.Response(http.StatusOK, 400001002, nil)
		return
	}

	_, err = service.GetTagById(tagIdInt)
	if err != nil {
		zgh.ZLog().Error("message", "console.Tag.Destroy", "err", err.Error())
		appG.Response(http.StatusOK, 403000008, nil)
		return
	}
	service.DelTagRel(tagIdInt)
	appG.Response(http.StatusOK, 0, nil)
	return
}
