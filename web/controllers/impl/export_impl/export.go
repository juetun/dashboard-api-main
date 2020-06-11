/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:09 上午
 */
package export_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/controllers/inter"
	"github.com/juetun/dashboard-api-main/web/pojos"
	"github.com/juetun/dashboard-api-main/web/services"
)

type ControllerExportData struct {
	base.ControllerBase
}

func (r ControllerExportData) Cancel(c *gin.Context) {
	var args pojos.ArgumentsExportCancel
	var err error

	err = c.Bind(&args)
	if err != nil {
		r.Log.Errorln("message", "export.Cancel", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	args.User = r.GetUser(c)
	srv := services.NewServiceExport(&base.Context{Log: r.Log})
	var res pojos.ResultExportCancel
	res, err = srv.Cancel(&args)
	if err != nil {
		r.Log.Errorln("message", "export.Cancel", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) Init(c *gin.Context) {
	var args = pojos.ArgumentsExportInit{}
	var err error

	err = c.Bind(&args)
	if err != nil {
		r.Log.Errorln("message", "export.Init", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	args.User = r.GetUser(c)
	args.HttpHeader = c.Request.Header
	srv := services.NewServiceExport(&base.Context{Log: r.Log})
	var res pojos.ResultExportInit
	res, err = srv.Init(&args)
	if err != nil {
		r.Log.Errorln("message", "export.Init", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) Progress(c *gin.Context) {
	var args pojos.ArgumentsExportProgress
	var err error

	err = c.BindQuery(&args)
	if err != nil {
		r.Log.Errorln("message", "export.Progress", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	args.InitIds()
	args.User = r.GetUser(c)

	srv := services.NewServiceExport(&base.Context{Log: r.Log})
	var res pojos.ResultExportProgress
	res, err = srv.Progress(&args)
	if err != nil {
		r.Log.Errorln("message", "export.Progress", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) List(c *gin.Context) {
	var args pojos.ArgumentsExportList
	var err error

	err = c.Bind(&args)
	if err != nil {
		r.Log.Errorln("message", "export.list", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	args.User = r.GetUser(c)
	srv := services.NewServiceExport(&base.Context{Log: r.Log})
	var res pojos.ResultExportList
	res, err = srv.List(&args)
	if err != nil {
		r.Log.Errorln("message", "export.list", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func NewControllerExportData() inter.Export {
	controller := &ControllerExportData{}
	controller.ControllerBase.Init()
	return controller
}
