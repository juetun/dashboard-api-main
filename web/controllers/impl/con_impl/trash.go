package con_impl

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/controllers/inter"
	"github.com/juetun/dashboard-api-main/web/pojos"
	"github.com/juetun/dashboard-api-main/web/services"
)

type ControllerTrash struct {
	base.ControllerBase
}

func NewControllerTrash() inter.Trash {
	controller := &ControllerTrash{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerTrash) Index(c *gin.Context) {

	pager := base.NewPager()
	limit, offset := pager.InitPageBy(c, "GET")
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	dba, postCount, err := srv.ConsolePostCount(limit, offset, false)

	var postList = &[]pojos.ConsolePostList{}
	if postCount > 0 {
		postList, err = srv.ConsolePostIndex(dba, limit, offset, false)
		if err != nil {
			r.Log.Errorln("message", "console.Index", "err", err.Error())
			r.Response(c, 500000000, nil)
			return
		}
	}

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = utils.MyPaginate(postCount, limit, pager.PageNo)

	r.Response(c, 0, data)
	return
}
func (r *ControllerTrash) Create(c *gin.Context) {
	srv := services.NewCategoryService(&base.Context{Log: r.Log})
	cates, err := srv.CateListBySort()

	if err != nil {
		r.Log.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srvTag := services.NewTagService(&base.Context{Log: r.Log})
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
func (r *ControllerTrash) Store(c *gin.Context) {
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
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	srv.PostStore(ps, userId.(int))
	r.Response(c, 0, nil)
	return
}
func (r *ControllerTrash) Edit(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		r.Log.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	srvTag := services.NewTagService(&base.Context{Log: r.Log})
	srvCate := services.NewCategoryService(&base.Context{Log: r.Log})

	post, err := srv.PostDetail(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	postTags, err := srv.PostIdTag(postIdInt)
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
func (r *ControllerTrash) Update(c *gin.Context) {
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
	var ps pojos.PostStore
	ps, ok := requestJson.(pojos.PostStore)
	if !ok {
		r.Log.Errorln("message", "post.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := services.NewConsolePostService()
	srv.PostUpdate(postIdInt, ps)
	r.Response(c, 0, nil)
	return
}
func (r *ControllerTrash) Destroy(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	_, err = srv.PostDestroy(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, nil)
	return
}
func (r *ControllerTrash) TrashIndex(c *gin.Context) {

	pager := base.NewPager()
	limit, offset := pager.InitPageBy(c, "GET")
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	dba, postCount, err := srv.ConsolePostCount(limit, offset, true)
	if err != nil {
		r.Log.Errorln("message", "console.TrashIndex", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	var postList = &[]pojos.ConsolePostList{}
	if postCount > 0 {
		postList, err = srv.ConsolePostIndex(dba, limit, offset, true)
		if err != nil {
			r.Log.Errorln("message", "console.TrashIndex", "err", err.Error())
			r.Response(c, 500000000, nil)
			return
		}
	}

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = utils.MyPaginate(postCount, limit, pager.PageNo)
	r.Response(c, 0, data)
	return
}
func (r *ControllerTrash) UnTrash(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	_, err = srv.PostUnTrash(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.UnTrash", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}
func (r *ControllerTrash) ImgUpload(c *gin.Context) {

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
	srvQiniu := services.NewQiuNiuService(&base.Context{Log: r.Log})
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
