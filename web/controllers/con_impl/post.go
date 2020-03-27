/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-27
 * Time: 00:14
 */
package con_impl

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/controllers"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services"
)

type ControllerPost struct {
	base.ControllerBase
}

func NewControllerPost() controllers.Console {
	controller := &ControllerPost{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerPost) Index(c *gin.Context) {
	queryPage := c.DefaultQuery("page", base.DefaultPageNo)
	queryLimit := c.DefaultQuery("limit", base.DefaultPageSize)

	limit, offset := common.Offset(queryPage, queryLimit)
	queryPageInt, err := strconv.Atoi(queryPage)

	srv := services.NewConsolePostService()
	postList, err := srv.ConsolePostIndex(limit, offset, false)
	if err != nil {
		r.Response(c, 500000000, nil)
		return
	}

	postCount, err := srv.ConsolePostCount(limit, offset, false)
	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, queryPageInt)
	r.Response(c, 0, data)
	return
}

func (r *ControllerPost) Create(c *gin.Context) {
	srv := services.NewCategoryService()
	cates, err := srv.CateListBySort()

	if err != nil {
		r.Log.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srvTag := services.NewTagService()
	tags, err := srvTag.AllTags()
	if err != nil {
		r.Log.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = common.Conf.ImgUploadUrl
	r.Response(c, 0, data)
	return
}

func (r *ControllerPost) Store(c *gin.Context) {

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Errorln("message", "post.Store", "error", "get request_params from context fail")
		r.Response(c, 401000004, nil)
		return
	}
	var ps pojos.PostStore
	ps, ok := requestJson.(pojos.PostStore)
	if !ok {
		r.Log.Errorln("message", "post.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}

	userId, exist := c.Get("userId")
	if !exist || userId.(int) == 0 {
		r.Log.Errorln("message", "post.Store", "error", "Can not get user")
		r.Response(c, 400001004, nil)
		return
	}
	srvPost := services.NewConsolePostService()
	srvPost.PostStore(ps, userId.(int))
	r.Response(c, 0, nil)
	return
}

func (r *ControllerPost) Edit(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srvPost := services.NewConsolePostService()
	srvCate := services.NewCategoryService()
	srvTag := services.NewTagService()

	post, err := srvPost.PostDetail(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	postTags, err := srvPost.PostIdTag(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	postCate, err := srvCate.PostCate(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	posts := make(map[string]interface{})
	posts["post"] = post
	posts["postCate"] = postCate
	posts["postTag"] = postTags
	data["post"] = posts
	cates, err := srvCate.CateListBySort()
	if err != nil {
		r.Log.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	tags, err := srvTag.AllTags()
	if err != nil {
		r.Log.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = common.Conf.ImgUploadUrl
	r.Response(c, 0, data)
	return
}

func (r *ControllerPost) Update(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Update", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Errorln("message", "post.Store", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil)
		return
	}
	srv := services.NewConsolePostService()
	var ps pojos.PostStore
	ps, ok := requestJson.(pojos.PostStore)
	if !ok {
		r.Log.Errorln("message", "post.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv.PostUpdate(postIdInt, ps)
	r.Response(c, 0, nil)
	return
}

func (r *ControllerPost) Destroy(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := services.NewConsolePostService()
	_, err = srv.PostDestroy(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}

func (r *ControllerPost) TrashIndex(c *gin.Context) {

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", common.Conf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)

	srv := services.NewConsolePostService()
	postList, err := srv.ConsolePostIndex(limit, offset, true)
	if err != nil {
		r.Log.Errorln("message", "console.TrashIndex", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		r.Log.Errorln("message", "console.TrashIndex", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	postCount, err := srv.ConsolePostCount(limit, offset, true)

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, queryPageInt)

	r.Response(c, 0, data)
	return
}

func (r *ControllerPost) UnTrash(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := services.NewConsolePostService()
	_, err = srv.PostUnTrash(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.UnTrash", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}

func (r *ControllerPost) ImgUpload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		r.Log.Infoln("message", "post.ImgUpload", "err", err.Error())
		r.Response(c, 401000004, nil)
		return
	}

	filename := filepath.Base(file.Filename)
	dst := common.Conf.ImgUploadDst + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		r.Log.Infoln("message", "post.ImgUpload", "error", err.Error())
		r.Response(c, 401000005, nil)
		return
	}

	srvQiniu := services.NewQiuNiuService()
	// Default upload both
	data := make(map[string]interface{})
	if common.Conf.ImgUploadBoth {
		go srvQiniu.Qiniu(dst, filename)
		data["path"] = common.Conf.AppImgUrl + filename
		r.Response(c, 0, data)
		return
	}
	if common.Conf.QiNiuUploadImg {
		go srvQiniu.Qiniu(dst, filename)
		data["path"] = common.Conf.QiNiuHostName + filename
	} else {
		data["path"] = common.Conf.AppImgUrl + filename
	}
	r.Response(c, 0, data)
	return
}
