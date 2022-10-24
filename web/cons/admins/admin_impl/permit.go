// Package admin_impl
package admin_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ControllerPermit struct {
	base.ControllerBase
}


func (r *ControllerPermit) GetMenu(c *gin.Context) {
	var (
		err error
		arg wrappers.ArgGetMenu
		res wrappers.ResultGetMenu
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(ctx).
		GetMenu(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}
func (r *ControllerPermit) AdminMenuSearch(c *gin.Context) {
	var (
		arg wrappers.ArgAdminMenu
		err error
		res wrappers.ResAdminMenuSearch
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(ctx).
		AdminMenuSearch(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func (r *ControllerPermit) AdminUserGroupRelease(c *gin.Context) {
	var (
		arg wrappers.ArgAdminUserGroupRelease
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	srv := srv_impl.NewSrvPermitGroupImpl(ctx)
	res, err := srv.AdminUserGroupRelease(&arg)

	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminUserGroupAdd(c *gin.Context) {
	var (
		arg wrappers.ArgAdminUserGroupAdd
		err error
		res wrappers.ResultAdminUserGroupAdd
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvPermitGroupImpl(ctx).
		AdminUserGroupAdd(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminGroupDelete(c *gin.Context) {
	var (
		arg wrappers.ArgAdminGroupDelete
		err error
		res wrappers.ResultAdminGroupDelete
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}
	ctx.Info(map[string]interface{}{"arg": arg})
	if res, err = srv_impl.NewSrvPermitGroupImpl(ctx).
		AdminGroupDelete(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminGroupEdit(c *gin.Context) {
	var (
		arg wrappers.ArgAdminGroupEdit
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	res, err := srv_impl.NewSrvPermitGroupImpl(ctx).
		AdminGroupEdit(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) MenuDelete(c *gin.Context) {

	var (
		arg wrappers.ArgMenuDelete
		err error
		res *wrappers.ResultMenuDelete
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(ctx).
		MenuDelete(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func (r *ControllerPermit) MenuSave(c *gin.Context) {
	var (
		arg wrappers.ArgMenuSave
		err error
		res *wrappers.ResultMenuSave
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		MenuSave(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) MenuAdd(c *gin.Context) {
	var (
		arg wrappers.ArgMenuAdd
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	// 记录日志
	res, err := srv_impl.NewPermitServiceImpl(ctx).
		MenuAdd(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminMenuWithCheck(c *gin.Context) {
	var (
		arg wrappers.ArgAdminMenuWithCheck
		res *wrappers.ResultMenuWithCheck
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvPermitMenu(ctx).
		AdminMenuWithCheck(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminMenu(c *gin.Context) {
	var (
		arg wrappers.ArgAdminMenu
		res *wrappers.ResultAdminMenu
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志

	if res, err = srv_impl.NewSrvPermitMenu(ctx).
		AdminMenu(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

// AdminSetPermit 设置权限
func (r *ControllerPermit) AdminSetPermit(c *gin.Context) {
	var (
		arg wrappers.ArgAdminSetPermit
		err error
		res *wrappers.ResultAdminSetPermit
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminSetPermit(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminGroup(c *gin.Context) {
	var (
		arg wrappers.ArgAdminGroup
		err error
		res *wrappers.ResultAdminGroup
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitGroupImpl(ctx).
		AdminGroup(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

// Menu 获取权限菜单
func (r *ControllerPermit) Menu(c *gin.Context) {
	var (
		err error
		arg wrappers.ArgPermitMenu
		res *wrappers.ResultPermitMenuReturn
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if r.ParametersAccept(ctx, &arg) {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitMenu(ctx).
		Menu(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) GetAppConfig(c *gin.Context) {

	var (
		arg wrappers.ArgGetAppConfig
		res *wrappers.ResultGetAppConfig
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	if res, err = srv_impl.NewSrvPermitAppImpl(ctx).
		GetAppConfig(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func NewControllerPermit() admins.Permit {
	controller := &ControllerPermit{}
	controller.ControllerBase.Init()
	return controller
}
