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
	var (
		args wrappers.ArgumentsExportCancel
		err  error
		res  wrappers.ResultExportCancel
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &args, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	srv := srv_impl.NewServiceExport(ctx)

	if res, err = srv.Cancel(&args); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) Init(c *gin.Context) {
	var (
		args = wrappers.ArgumentsExportInit{}
		err  error
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &args, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	srv := srv_impl.NewServiceExport(ctx)
	var res wrappers.ResultExportInit
	res, err = srv.Init(&args)
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) Progress(c *gin.Context) {
	var (
		args wrappers.ArgumentsExportProgress
		err  error
		res  wrappers.ResultExportProgress
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &args, base.ControllerGetParamTypeBindQuery); haveErr {
		return
	}
	if res, err = srv_impl.NewServiceExport(ctx).
		Progress(&args); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, res)
}

func (r ControllerExportData) List(c *gin.Context) {

	var (
		args wrappers.ArgumentsExportList
		err  error
		res  wrappers.ResultExportList
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &args, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	if res, err = srv_impl.NewServiceExport(ctx).
		List(&args); err != nil {
		r.Log.Logger.Errorln("message", "export.list", "err", err.Error())
		r.ResponseError(c, err)
		return
	}
	r.Response(c, 0, res)
}

func NewControllerExportData() admins.Export {
	controller := &ControllerExportData{}
	controller.ControllerBase.Init()
	return controller
}
