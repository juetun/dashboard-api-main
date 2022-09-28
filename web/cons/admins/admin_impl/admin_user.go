package admin_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type ConAdminUser struct {
	base.ControllerBase
}

func (r *ConAdminUser) AdminUserUpdateWithColumn(c *gin.Context) {
	var (
		err error
		arg wrapper_admin.ArgAdminUserUpdateWithColumn
		res *wrapper_admin.ResultAdminUserUpdateWithColumn
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvPermitUserImpl(ctx).
		AdminUserUpdateWithColumn(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *ConAdminUser) AdminUser(c *gin.Context) {
	var (
		err error
		arg wrappers.ArgAdminUser
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	res, err := srv_impl.NewSrvPermitUserImpl(ctx).
		AdminUser(&arg)

	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ConAdminUser) AdminUserEdit(c *gin.Context) {

	var (
		res wrappers.ResultAdminUserAdd
		arg wrappers.ArgAdminUserAdd
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}
	// 记录日志
	res, err = srv_impl.NewSrvPermitUserImpl(ctx).
		AdminUserEdit(&arg)

	if err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ConAdminUser) AdminUserDelete(c *gin.Context) {
	var (
		arg wrappers.ArgAdminUserDelete
		err error
		res wrappers.ResultAdminUserDelete
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}
	// 记录日志
	if res, err = srv_impl.NewSrvPermitUserImpl(ctx).
		AdminUserDelete(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func NewConAdminUser() (res admins.ConAdminUser) {
	p := &ConAdminUser{}
	p.ControllerBase.Init()
	return p
}
