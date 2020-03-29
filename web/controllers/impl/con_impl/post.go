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
	"github.com/jinzhu/gorm"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/controllers/inter"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services"
)

type ControllerPost struct {
	base.ControllerBase
}

func NewControllerPost() inter.Console {
	controller := &ControllerPost{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerPost) Index(c *gin.Context) {
	pager := base.NewPager()
	limit, offset := pager.InitPageBy(c, "GET")
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	dba, postCount, err := srv.ConsolePostCount(limit, offset, false)
	postList := &[]pojos.ConsolePostList{}
	if postCount > 0 {
		postList, err = srv.ConsolePostIndex(dba, limit, offset, false)
		if err != nil {
			r.Response(c, 500000000, nil)
			return
		}
	}
	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, pager.PageNo)
	r.Response(c, 0, data)
	return
}

func (r *ControllerPost) Create(c *gin.Context) {
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
	srvPost := services.NewConsolePostService(&base.Context{Log: r.Log})
	srvPost.PostStore(ps, userId.(int))
	r.Response(c, 0, nil)
	return
}

func (r *ControllerPost) Edit(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000002, nil, err.Error())
		return
	}
	srvPost := services.NewConsolePostService(&base.Context{Log: r.Log})
	post, err := srvPost.PostDetail(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Edit(116)", "err", err.Error())
		r.Response(c, 500000001, nil, err.Error())
		return
	}

	srvCate := services.NewCategoryService(&base.Context{Log: r.Log})
	srvTag := services.NewTagService(&base.Context{Log: r.Log})
	postTags, err := srvPost.PostIdTag(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Edit(125)", "err", err.Error())
		r.Response(c, 500000005, nil)
		return
	}
	var postCate *map[int]models.ZPostCate
	postCate, err = srvCate.GetPostCates(&[]int{postIdInt})
	if err != nil {
		r.Log.Errorln("message", "console.Edit(12)", "err", err.Error())
		r.Response(c, 500000006, nil)
		return
	}
	data := make(map[string]interface{})
	posts := make(map[string]interface{})
	posts["post"] = post
	posts["postCate"] = 0
	if _, ok := (*postCate)[postIdInt]; ok {
		posts["postCate"] = (*postCate)[postIdInt].CateId
	}

	posts["postTag"] = postTags
	data["post"] = posts
	cates, err := srvCate.CateListBySort()
	if err != nil {
		r.Log.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000007, nil)
		return
	}
	tags, err := srvTag.AllTags()
	if err != nil {
		r.Log.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000008, nil)
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
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
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
	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	_, err = srv.PostDestroy(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerPost) TrashIndex(c *gin.Context) {

	pager := base.NewPager()
	limit, offset := pager.InitPageBy(c, "GET")

	srv := services.NewConsolePostService(&base.Context{Log: r.Log})
	var dba *gorm.DB
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
	data["page"] = common.MyPaginate(postCount, limit, pager.PageNo)

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
