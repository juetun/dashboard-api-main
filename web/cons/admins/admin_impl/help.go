package admin_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type ConsoleHelpImpl struct {
	base.ControllerBase
}

func (r *ConsoleHelpImpl) OperateLog(c *gin.Context) {
	var (
		arg wrapper_admin.ArgOperateLog
		res *wrapper_admin.ResultOperateLog
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvOperate(ctx).
		OperateLog(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *ConsoleHelpImpl) HelpTree(c *gin.Context) {

	var (
		arg wrapper_admin.ArgHelpTree
		res *wrapper_admin.ResultHelpTree
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvHelpRelate(ctx).
		HelpTree(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *ConsoleHelpImpl) TreeEditNode(c *gin.Context) {
	var (
		arg wrapper_admin.ArgTreeEditNode
		res *wrapper_admin.ResultTreeEditNode
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvHelpRelate(ctx).
		TreeEditNode(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *ConsoleHelpImpl) HelpList(c *gin.Context) {

	var (
		arg wrapper_admin.ArgHelpList
		res *wrapper_admin.ResultHelpList
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvHelp(ctx).
		HelpList(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *ConsoleHelpImpl) HelpDetail(c *gin.Context) {

	var (
		arg wrapper_admin.ArgHelpDetail
		res *wrapper_admin.ResultHelpDetail
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}
	if res, err = srv_impl.NewSrvHelp(ctx).
		HelpDetail(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *ConsoleHelpImpl) HelpEdit(c *gin.Context) {

	var (
		arg wrapper_admin.ArgHelpEdit
		res *wrapper_admin.ResultHelpEdit
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}
	if res, err = srv_impl.NewSrvHelp(ctx).
		HelpEdit(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func NewConsoleHelp() admins.ConsoleHelp {
	controller := &ConsoleHelpImpl{}
	controller.ControllerBase.Init()
	return controller
}
