// Package admin_impl
/**
* @Author:ChangJiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:09 上午
 */
package admin_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ControllerExportData struct {
	base.ControllerBase
}

func (r ControllerExportData) Cancel(c *gin.Context) {
	var args wrappers.ArgumentsExportCancel
	var err error
	var res wrappers.ResultExportCancel

	if err = c.Bind(&args); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = args.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	srv := srv_impl.NewServiceExport(base.CreateContext(&r.ControllerBase, c))

	if res, err = srv.Cancel(&args); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) Init(c *gin.Context) {
	var args = wrappers.ArgumentsExportInit{}
	var err error

	if err = c.Bind(&args); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = args.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	srv := srv_impl.NewServiceExport(base.CreateContext(&r.ControllerBase, c))
	var res wrappers.ResultExportInit
	res, err = srv.Init(&args)
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) Progress(c *gin.Context) {
	var args wrappers.ArgumentsExportProgress
	var err error

	if err = c.BindQuery(&args); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = args.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	srv := srv_impl.NewServiceExport(base.CreateContext(&r.ControllerBase, c))
	var res wrappers.ResultExportProgress
	res, err = srv.Progress(&args)
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) List(c *gin.Context) {
	var args wrappers.ArgumentsExportList
	var err error

	if err = c.Bind(&args); err != nil {
 		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = args.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	srv := srv_impl.NewServiceExport(base.CreateContext(&r.ControllerBase, c))
	var res wrappers.ResultExportList
	res, err = srv.List(&args)
	if err != nil {
		r.Log.Logger.Errorln("message", "export.list", "err", err.Error())
		r.Response(c, 500000000, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func NewControllerExportData() admins.Export {
	controller := &ControllerExportData{}
	controller.ControllerBase.Init()
	return controller
}
