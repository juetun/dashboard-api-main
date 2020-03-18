/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-27
 * Time: 00:14
 */
package con_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/study/app-content/common"
	"github.com/juetun/study/app-content/conf"
	"github.com/juetun/study/app-content/service"
	"github.com/juetun/study/app-dashboard/gin/api"
	"github.com/juetun/study/app-dashboard/lib/base"
	"github.com/juetun/study/app-dashboard/web/controllers"
	"github.com/juetun/study/app-dashboard/web/services"
	"net/http"
	"path/filepath"
	"strconv"
)

type ControllerPost struct {
	base.ControllerBase
}

func NewControllerPost() controllers.Console {
	controller := &ControllerPost{}
	controller.ControllerBase.Init()
	return controller
}

func (p *ControllerPost) Index(c *gin.Context) {
	pageNo := c.DefaultQuery("page", base.DefaultPageNo)
	pageSize := c.DefaultQuery("limit", base.DefaultPageSize)
	srv := services.NewConsolePostService()
	postList, err := srv.ConsolePostIndex(pageNo, pageSize, false)
	if err != nil {
		c.Response(http.StatusOK, 500000000, nil)
		return
	}

	postCount, err := service.ConsolePostCount(limit, offset, false)

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *ControllerPost) Create(c *gin.Context) {
	cates, err := service.CateListBySort()
	appG := api.Gin{C: c}
	if err != nil {
		p.Log.Error("message", "console.Create", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	tags, err := service.AllTags()
	if err != nil {
		p.Log.Error("message", "console.Create", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = conf.Cnf.ImgUploadUrl
	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *ControllerPost) Store(c *gin.Context) {
	appG := api.Gin{C: c}
	requestJson, exists := c.Get("json")
	if !exists {
		p.Log.Error("message", "post.Store", "error", "get request_params from context fail")
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	var ps common.PostStore
	ps, ok := requestJson.(common.PostStore)
	if !ok {
		p.Log.Error("message", "post.Store", "error", "request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}

	userId, exist := c.Get("userId")
	if !exist || userId.(int) == 0 {
		p.Log.Error("message", "post.Store", "error", "Can not get user")
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}

	service.PostStore(ps, userId.(int))
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *ControllerPost) Edit(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		p.Log.Error("message", "console.Edit", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	post, err := service.PostDetail(postIdInt)
	if err != nil {
		p.Log.Error("message", "console.Edit", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postTags, err := service.PostIdTag(postIdInt)
	if err != nil {
		p.Log.Error("message", "console.Edit", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postCate, err := service.PostCate(postIdInt)
	if err != nil {
		p.Log.Error("message", "console.Edit", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	posts := make(map[string]interface{})
	posts["post"] = post
	posts["postCate"] = postCate
	posts["postTag"] = postTags
	data["post"] = posts
	cates, err := service.CateListBySort()
	if err != nil {
		p.Log.Error("message", "console.Create", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	tags, err := service.AllTags()
	if err != nil {
		p.Log.Error("message", "console.Create", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = conf.Cnf.ImgUploadUrl
	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *ControllerPost) Update(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		p.Log.Error("message", "console.Update", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		p.Log.Error("message", "post.Store", "error", "get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	var ps common.PostStore
	ps, ok := requestJson.(common.PostStore)
	if !ok {
		p.Log.Error("message", "post.Store", "error", "request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	service.PostUpdate(postIdInt, ps)
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *ControllerPost) Destroy(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		p.Log.Error("message", "console.Destroy", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	_, err = service.PostDestroy(postIdInt)
	if err != nil {
		p.Log.Error("message", "console.Destroy", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *ControllerPost) TrashIndex(c *gin.Context) {
	appG := api.Gin{C: c}

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)
	postList, err := service.ConsolePostIndex(limit, offset, true)
	if err != nil {
		p.Log.Error("message", "console.TrashIndex", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		p.Log.Error("message", "console.TrashIndex", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postCount, err := service.ConsolePostCount(limit, offset, true)

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *ControllerPost) UnTrash(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		p.Log.Error("message", "console.Destroy", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	_, err = service.PostUnTrash(postIdInt)
	if err != nil {
		p.Log.Error("message", "console.UnTrash", "err", err.Error())
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *ControllerPost) ImgUpload(c *gin.Context) {
	appG := api.Gin{C: c}

	file, err := c.FormFile("file")
	if err != nil {
		p.Log.Info("message", "post.ImgUpload", "err", err.Error())
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}

	filename := filepath.Base(file.Filename)
	dst := conf.Cnf.ImgUploadDst + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		p.Log.Info("message", "post.ImgUpload", "error", err.Error())
		appG.Response(http.StatusOK, 401000005, nil)
		return
	}

	// Default upload both
	data := make(map[string]interface{})
	if conf.Cnf.ImgUploadBoth {
		go service.Qiniu(dst, filename)
		data["path"] = conf.Cnf.AppImgUrl + filename
	} else {
		if conf.Cnf.QiNiuUploadImg {
			go service.Qiniu(dst, filename)
			data["path"] = conf.Cnf.QiNiuHostName + filename
		} else {
			data["path"] = conf.Cnf.AppImgUrl + filename
		}
	}

	appG.Response(http.StatusOK, 0, data)
	return
}
