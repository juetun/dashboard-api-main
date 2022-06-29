// Package admin_impl
/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-27
 * Time: 00:14
 */
package admin_impl

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/common/response"
	cons_admin2 "github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type ControllerPost struct {
	base.ControllerBase
}

func NewControllerPost() cons_admin2.Console {
	controller := &ControllerPost{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerPost) Index(c *gin.Context) {
	var arg response.PageQuery

	var err error
	arg.PageNo, err = strconv.Atoi(c.DefaultQuery("pages", strconv.Itoa(response.DefaultPageNo)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	arg.PageSize, err = strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(response.DefaultPageSize)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	arg.DefaultPage()
	offset := arg.GetOffset()
	pager := response.NewPager(response.PagerBaseQuery(&arg))
	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	dba, postCount, err := srv.ConsolePostCount(pager.PageSize, offset, false)
	pager.List = &[]wrappers.ConsolePostList{}
	if postCount > 0 {
		pager.List, err = srv.ConsolePostIndex(dba, pager.PageSize, offset, false)
		if err != nil {
			r.Response(c, 500000000, nil)
			return
		}
	}
	data := make(map[string]interface{})
	r.Response(c, 0, data)
	return
}

func (r *ControllerPost) Create(c *gin.Context) {
	srv := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))
	cates, err := srv.CateListBySort()

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	srvTag := srv_impl.NewTagService(base.CreateContext(&r.ControllerBase, c))
	tags, err := srvTag.AllTags()
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = common.Conf.ImgUploadUrl
	r.Response(c, 0, data, "添加成功")
	return
}

func (r *ControllerPost) Store(c *gin.Context) {

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Logger.Errorln("message", "post.Store", "error", "get request_params from context fail")
		r.Response(c, 401000004, nil, "参数异常")
		return
	}
	var ps wrappers.PostStore
	ps, ok := requestJson.(wrappers.PostStore)
	if !ok {
		r.Log.Logger.Errorln("message", "post.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil, "参数异常")
		return
	}
	userId := r.GetUser(c).UserId
	srvPost := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	srvPost.PostStore(ps, userId)
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerPost) Edit(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000002, nil, err.Error())
		return
	}
	srvPost := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	post, err := srvPost.PostDetail(postIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Edit(116)", "err", err.Error())
		r.Response(c, 500000001, nil, err.Error())
		return
	}

	srvCate := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))
	srvTag := srv_impl.NewTagService(base.CreateContext(&r.ControllerBase, c))
	postTags, err := srvPost.PostIdTag(postIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Edit(125)", "err", err.Error())
		r.Response(c, 500000005, nil)
		return
	}
	var postCate *map[string]models.ZPostCate
	postCate, err = srvCate.GetPostCates(&[]int{postIdInt})
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Edit(12)", "err", err.Error())
		r.Response(c, 500000006, nil)
		return
	}
	data := make(map[string]interface{})
	posts := make(map[string]interface{})
	posts["post"] = post
	posts["postCate"] = 0
	if _, ok := (*postCate)[postIdStr]; ok {
		posts["postCate"] = (*postCate)[postIdStr].CateId
	}

	posts["postTag"] = postTags
	data["post"] = posts
	cates, err := srvCate.CateListBySort()
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000007, nil)
		return
	}
	tags, err := srvTag.AllTags()
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000008, nil)
		return
	}
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = common.Conf.ImgUploadUrl
	r.Response(c, 0, data, "操作成功")
	return
}

func (r *ControllerPost) Update(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Update", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Logger.Errorln("message", "post.Store", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil, "参数格式异常")
		return
	}
	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	var ps wrappers.PostStore
	ps, ok := requestJson.(wrappers.PostStore)
	if !ok {
		r.Log.Logger.Errorln("message", "post.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil, "参数格式异常")
		return
	}
	if err = srv.PostUpdate(postIdInt, ps); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerPost) Destroy(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	_, err = srv.PostDestroy(postIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerPost) TrashIndex(c *gin.Context) {
	var arg response.PageQuery

	var err error
	arg.PageNo, err = strconv.Atoi(c.DefaultQuery("pages", strconv.Itoa(response.DefaultPageNo)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	arg.PageSize, err = strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(response.DefaultPageSize)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	arg.DefaultPage()
	offset := arg.GetOffset()
	pager := response.NewPager(response.PagerBaseQuery(&arg))
	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	var dba *gorm.DB
	dba, postCount, err := srv.ConsolePostCount(pager.PageSize, offset, true)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.TrashIndex", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	if postCount > 0 {
		pager.List, err = srv.ConsolePostIndex(dba, pager.PageSize, offset, true)
		if err != nil {
			r.Log.Logger.Errorln("message", "console.TrashIndex", "err", err.Error())
			r.Response(c, 500000000, nil)
			return
		}
	}

	r.Response(c, 0, pager)
	return
}

func (r *ControllerPost) UnTrash(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Destroy", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	_, err = srv.PostUnTrash(postIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.UnTrash", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, nil, "操作成功")
	return
}

func (r *ControllerPost) ImgUpload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		r.Log.Logger.Infoln("message", "post.ImgUpload", "err", err.Error())
		r.Response(c, 401000004, nil)
		return
	}

	filename := filepath.Base(file.Filename)
	dst := common.Conf.ImgUploadDst + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		r.Log.Logger.Infoln("message", "post.ImgUpload", "error", err.Error())
		r.Response(c, 401000005, nil)
		return
	}

	srvQiniu := srv_impl.NewQiuNiuService(base.CreateContext(&r.ControllerBase, c))
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
