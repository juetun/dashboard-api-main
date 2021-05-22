package con_impl

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/cons_outernet"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ControllerTrash struct {
	base.ControllerBase
}

func NewControllerTrash() cons_outernet.Trash {
	controller := &ControllerTrash{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerTrash) Index(c *gin.Context) {

	pager := response.NewPager()
	var err error
	pager.PageNo, err = strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(response.DefaultPageNo)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	pager.PageSize, err = strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(response.DefaultPageSize)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	pager.DefaultPage()
	offset := pager.GetOffset()
	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	dba, postCount, err := srv.ConsolePostCount(pager.PageSize, offset, false)

	var postList = &[]wrappers.ConsolePostList{}
	if postCount > 0 {
		postList, err = srv.ConsolePostIndex(dba, pager.PageSize, offset, false)
		if err != nil {
			r.Log.Logger.Errorln("message", "console.Index", "err", err.Error())
			r.Response(c, 500000000, nil)
			return
		}
	}

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = utils.MyPaginate(postCount, pager.PageSize, pager.PageNo)

	r.Response(c, 0, data)
	return
}
func (r *ControllerTrash) Create(c *gin.Context) {
	srv := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))
	cates, err := srv.CateListBySort()

	if err != nil {
		r.Log.Logger.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srvTag := srv_impl.NewTagService(base.CreateContext(&r.ControllerBase, c))
	tags, err := srvTag.AllTags()
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Create", "err", err.Error())
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
		r.Log.Logger.Errorln("message", "post.Store", "error", "get request_params from context fail")
		r.Response(c, 401000004, nil)
		return
	}
	var ps wrappers.PostStore
	ps, ok := requestJson.(wrappers.PostStore)
	if !ok {
		r.Log.Logger.Errorln("message", "post.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	userId := r.GetUser(c).UserId
	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	srv.PostStore(ps, userId)
	r.Response(c, 0, nil)
	return
}
func (r *ControllerTrash) Edit(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	srvTag := srv_impl.NewTagService(base.CreateContext(&r.ControllerBase, c))
	srvCate := srv_impl.NewCategoryService(base.CreateContext(&r.ControllerBase, c))

	post, err := srv.PostDetail(postIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	postTags, err := srv.PostIdTag(postIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Edit", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	postCate, err := srvCate.PostCate(postIdInt)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Edit", "err", err.Error())
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
		r.Log.Logger.Errorln("message", "console.Create", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	tags, err := srvTag.AllTags()
	if err != nil {
		r.Log.Logger.Errorln("message", "console.Create", "err", err.Error())
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
		r.Log.Logger.Errorln("message", "console.Update", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {
		r.Log.Logger.Errorln("message", "post.Store", "error", "get request_params from context fail")
		r.Response(c, 400001003, nil)
		return
	}
	var ps wrappers.PostStore
	ps, ok := requestJson.(wrappers.PostStore)
	if !ok {
		r.Log.Logger.Errorln("message", "post.Store", "error", "request_params turn to error")
		r.Response(c, 400001001, nil)
		return
	}
	srv := srv_impl.NewConsolePostService()
	srv.PostUpdate(postIdInt, ps)
	r.Response(c, 0, nil)
	return
}
func (r *ControllerTrash) Destroy(c *gin.Context) {
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
	r.Response(c, 0, nil)
	return
}
func (r *ControllerTrash) TrashIndex(c *gin.Context) {

	pager := response.NewPager()
	var err error
	pager.PageNo, err = strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(response.DefaultPageNo)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	pager.PageSize, err = strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(response.DefaultPageSize)))
	if err != nil {
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	pager.DefaultPage()
	offset := pager.GetOffset()

	srv := srv_impl.NewConsolePostService(base.CreateContext(&r.ControllerBase, c))
	dba, postCount, err := srv.ConsolePostCount(pager.PageSize, offset, true)
	if err != nil {
		r.Log.Logger.Errorln("message", "console.TrashIndex", "err", err.Error())
		r.Response(c, 500000000, nil)
		return
	}
	var postList = &[]wrappers.ConsolePostList{}
	if postCount > 0 {
		postList, err = srv.ConsolePostIndex(dba, pager.PageSize, offset, true)
		if err != nil {
			r.Log.Logger.Errorln("message", "console.TrashIndex", "err", err.Error())
			r.Response(c, 500000000, nil)
			return
		}
	}

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = utils.MyPaginate(postCount, pager.PageSize, pager.PageNo)
	r.Response(c, 0, data)
	return
}
func (r *ControllerTrash) UnTrash(c *gin.Context) {
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
func (r *ControllerTrash) ImgUpload(c *gin.Context) {

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
