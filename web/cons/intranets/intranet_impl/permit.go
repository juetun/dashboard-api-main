package intranet_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/intranets"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type ConPermitIntranetImpl struct {
	base.ControllerBase
}

func (r *ConPermitIntranetImpl) ValidateUserHavePermit(c *gin.Context) {
	var (
		err  error
		args wrapper_intranet.ArgValidateUserHavePermit
		res  *wrapper_intranet.ResultValidateUserHavePermit
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &args); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvPermitUserImpl(ctx).
		ValidateUserHavePermit(&args); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

// GetUerImportPermit 判断接口是否有权限访问
func (r *ConPermitIntranetImpl) GetUerImportPermit(c *gin.Context) {
	var (
		err  error
		args wrapper_intranet.ArgGetUerImportPermit
		res  *wrapper_intranet.ResultGetUerImportPermit
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &args); haveErr {
		return
	}
	srv := srv_impl.NewSrvGatewayImport(ctx)

	if res, err = srv.GetUerImportPermit(&args); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

// GetImportPermit 获取接口权限
func (r *ConPermitIntranetImpl) GetImportPermit(c *gin.Context) {
	var (
		args wrapper_intranet.ArgGetImportPermit
		err  error
		res  wrapper_intranet.ResultGetImportPermit
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &args, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	srv := srv_impl.NewSrvGatewayImport(ctx)

	if res, err = srv.GetImportPermit(&args); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func NewConPermitIntranet() (res intranets.ConPermitIntranet) {
	p := &ConPermitIntranetImpl{}
	p.ControllerBase.Init()
	return p
}
