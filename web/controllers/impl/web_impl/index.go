/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-20
 * Time: 23:36
 */
package web_impl

import (
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/controllers/inter"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services"
)

type ControllerWeb struct {
	// ApiController
	base.ControllerBase
}

func NewControllerWeb() inter.Web {
	controller := &ControllerWeb{}
	controller.ControllerBase.Init()
	return controller
}
func (r *ControllerWeb) Index(c *gin.Context) {
	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", common.Conf.DefaultIndexLimit)

	srv := services.NewIndexService(&base.Context{Log: r.Log})
	srcTag := services.NewTagService(&base.Context{Log: r.Log})
	h, err := srcTag.CommonData()
	if err != nil {
		r.Log.Errorln("message", "Index.Index", "err", err.Error())
		r.Response(c, 404, h)
		return
	}

	postData, err := srv.IndexPost(queryPage, queryLimit, "default", "")
	if err != nil {
		r.Log.Errorln("message", "Index.Index",
			"err", err.Error(),
		)
		r.Response(c, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["title"] = h["system"].(*models.ZBaseSys).Title
	r.Response(c, 0, h)
	return
}

func (r *ControllerWeb) IndexTag(c *gin.Context) {
	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", common.Conf.DefaultIndexLimit)
	name := c.Param("name")
	srcTag := services.NewTagService(&base.Context{Log: r.Log})
	h, err := srcTag.CommonData()
	if err != nil {
		r.Log.Errorln("message", "Index.Index", "err", err.Error())
		r.Response(c, 404, h)
		return
	}
	src := services.NewIndexService(&base.Context{Log: r.Log})
	postData, err := src.IndexPost(queryPage, queryLimit, "tag", name)
	if err != nil {
		r.Log.Errorln("message", "Index.Index", "err", err.Error())
		r.Response(c, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["tagName"] = name
	h["tem"] = "tagList"
	h["title"] = template.HTML(name + " --  tags &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*models.ZBaseSys).Title)
	c.HTML(http.StatusOK, "master.tmpl", h)
	return
}

func (r *ControllerWeb) IndexCate(c *gin.Context) {
	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", common.Conf.DefaultIndexLimit)
	name := c.Param("name")

	srv := services.NewTagService(&base.Context{Log: r.Log})
	h, err := srv.CommonData()
	if err != nil {
		r.Log.Errorln("message", "Index.IndexCate", "err", err.Error())
		r.Response(c, 404, h)
		return
	}

	srvIndex := services.NewIndexService(&base.Context{Log: r.Log})
	postData, err := srvIndex.IndexPost(queryPage, queryLimit, "cate", name)
	if err != nil {
		r.Log.Errorln("message", "Index.IndexCate", "err", err.Error())
		r.Response(c, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["cateName"] = name
	h["tem"] = "cateList"
	h["title"] = template.HTML(name + " --  category &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*models.ZBaseSys).Title)

	r.Response(c, 0, h)
	return

}

func (r *ControllerWeb) Detail(c *gin.Context) {
	postIdStr := c.Param("id")

	srv := services.NewTagService(&base.Context{Log: r.Log})
	h, err := srv.CommonData()
	if err != nil {
		r.Log.Errorln("message", "Index.Detail", "err", err.Error())
		r.Response(c, 404, h)
		return
	}

	srvIndex := services.NewIndexService(&base.Context{Log: r.Log})
	postDetail, err := srvIndex.IndexPostDetail(postIdStr)
	if err != nil {
		r.Log.Errorln("message", "Index.Detail", "err", err.Error())
		r.Response(c, 404, h)
		return
	}

	go srvIndex.PostViewAdd(postIdStr)

	github := pojos.IndexGithubParam{
		// GithubName:         common.Conf.GithubName,
		// GithubRepo:         common.Conf.GithubRepo,
		// GithubClientId:     common.Conf.GithubClientId,
		// GithubClientSecret: common.Conf.GithubClientSecret,
		// GithubLabels:       common.Conf.GithubLabels,
	}

	h["post"] = postDetail
	h["github"] = github
	h["tem"] = "detail"
	h["title"] = template.HTML(postDetail.Post.Title + " &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*models.ZBaseSys).Title)

	r.Response(c, 0, h)
	return
}

func (r *ControllerWeb) Archives(c *gin.Context) {

	srv := services.NewTagService(&base.Context{Log: r.Log})
	h, err := srv.CommonData()
	if err != nil {
		r.Log.Errorln("message", "Index.Archives", "err", err.Error())
		r.Response(c, 404, h)
		return
	}
	srvPost := services.NewConsolePostService(&base.Context{Log: r.Log})
	res, err := srvPost.PostArchives()
	if err != nil {
		r.Log.Errorln("message", "Index.Archives", "err", err.Error())
		r.Response(c, 404, h)
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")

	var dateIndexs []int
	for k, _ := range res {
		tt, _ := time.ParseInLocation("2006-01-02 15:04:05", k+"-01 00:00:00", loc)
		dateIndexs = append(dateIndexs, int(tt.Unix()))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(dateIndexs)))

	var newData []interface{}
	for _, j := range dateIndexs {
		dds := make(map[string]interface{})
		tm := time.Unix(int64(j), 0)
		dateIndex := tm.Format("2006-01")
		dds["dates"] = dateIndex
		dds["lists"] = res[dateIndex]
		newData = append(newData, dds)
	}

	h["tem"] = "archives"
	h["archives"] = newData
	h["title"] = template.HTML("归档 &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*models.ZBaseSys).Title)

	r.Response(c, 0, h)
	return
}

func (r *ControllerWeb) NoFound(c *gin.Context) {
	r.Response(c, 404, gin.H{
		"themeJs":  "/static/home/assets/js",
		"themeCss": "/static/home/assets/css",
	})
	return
}
