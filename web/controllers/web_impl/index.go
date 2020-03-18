/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-20
 * Time: 23:36
 */
package web_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/study/app-content/common"
	"github.com/juetun/study/app-content/conf"
	"github.com/juetun/study/app-content/entity"
	"github.com/juetun/study/app-content/service"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/web/controllers"
	"html/template"
	"net/http"
	"sort"
	"time"
)

type ControllerWeb struct {
	//ApiController
	base.ControllerBase
}

func NewControllerWeb() controllers.Web {
	controller := &ControllerWeb{}
	controller.ControllerBase.Init()
	return controller
}
func (w *ControllerWeb) Index(c *gin.Context) {
	w.C = c
	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)

	h, err := service.CommonData()
	if err != nil {
		zgh.ZLog().Error("message", "Index.Index", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(queryPage, queryLimit, "default", "")
	if err != nil {
		zgh.ZLog().Error("message", "Index.Index", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["title"] = h["system"].(*entity.ZSystems).Title
	w.Response(http.StatusOK, 0, h)
	return
}

func (w *ControllerWeb) IndexTag(c *gin.Context) {
	w.C = c
	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)
	name := c.Param("name")
	h, err := service.CommonData()
	if err != nil {
		zgh.ZLog().Error("message", "Index.Index", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(queryPage, queryLimit, "tag", name)
	if err != nil {
		zgh.ZLog().Error("message", "Index.Index", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["tagName"] = name
	h["tem"] = "tagList"
	h["title"] = template.HTML(name + " --  tags &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*entity.ZSystems).Title)

	c.HTML(http.StatusOK, "master.tmpl", h)
	return
}

func (w *ControllerWeb) IndexCate(c *gin.Context) {
	w.C = c
	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)
	name := c.Param("name")

	h, err := service.CommonData()
	if err != nil {
		zgh.ZLog().Error("message", "Index.IndexCate", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(queryPage, queryLimit, "cate", name)
	if err != nil {
		zgh.ZLog().Error("message", "Index.IndexCate", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["cateName"] = name
	h["tem"] = "cateList"
	h["title"] = template.HTML(name + " --  category &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*entity.ZSystems).Title)

	w.Response(http.StatusOK, 0, h)
	return

}

func (w *ControllerWeb) Detail(c *gin.Context) {
	w.C = c
	postIdStr := c.Param("id")

	h, err := service.CommonData()
	if err != nil {
		zgh.ZLog().Error("message", "Index.Detail", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	postDetail, err := service.IndexPostDetail(postIdStr)
	if err != nil {
		zgh.ZLog().Error("message", "Index.Detail", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	go service.PostViewAdd(postIdStr)

	github := common.IndexGithubParam{
		GithubName:         conf.Cnf.GithubName,
		GithubRepo:         conf.Cnf.GithubRepo,
		GithubClientId:     conf.Cnf.GithubClientId,
		GithubClientSecret: conf.Cnf.GithubClientSecret,
		GithubLabels:       conf.Cnf.GithubLabels,
	}

	h["post"] = postDetail
	h["github"] = github
	h["tem"] = "detail"
	h["title"] = template.HTML(postDetail.Post.Title + " &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*entity.ZSystems).Title)

	w.Response(http.StatusOK, 0, h)
	return
}

func (w *ControllerWeb) Archives(c *gin.Context) {
	w.C = c
	h, err := service.CommonData()
	if err != nil {
		zgh.ZLog().Error("message", "Index.Archives", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
		return
	}

	res, err := service.PostArchives()
	if err != nil {
		zgh.ZLog().Error("message", "Index.Archives", "err", err.Error())
		w.Response(http.StatusOK, 404, h)
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
	h["title"] = template.HTML("归档 &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*entity.ZSystems).Title)

	w.Response(http.StatusOK, 0, h)
	return
}

func (w *ControllerWeb) NoFound(c *gin.Context) {
	w.C = c
	w.Response(http.StatusOK, 404, gin.H{
		"themeJs":  "/static/home/assets/js",
		"themeCss": "/static/home/assets/css",
	})
	return
}
