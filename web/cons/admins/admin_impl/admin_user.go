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
	)

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	if res, err = srv_impl.NewSrvPermitUserImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminUserUpdateWithColumn(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *ConAdminUser) AdminUser(c *gin.Context) {
	var arg wrappers.ArgAdminUser
	err := c.Bind(&arg)
	if err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	res, err := srv_impl.NewSrvPermitUserImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminUser(&arg)

	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ConAdminUser) AdminUserAdd(c *gin.Context) {
	var arg wrappers.ArgAdminUserAdd
	var err error

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		return
	}

	// 记录日志
	res, err := srv_impl.NewSrvPermitUserImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminUserAdd(&arg)

	if err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ConAdminUser) AdminUserDelete(c *gin.Context) {
	var arg wrappers.ArgAdminUserDelete
	var err error
	var res wrappers.ResultAdminUserDelete

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	// 记录日志
	if res, err = srv_impl.NewSrvPermitUserImpl(base.CreateContext(&r.ControllerBase, c)).
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
