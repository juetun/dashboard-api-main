package admin_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type AdminConPermitImportImpl struct {
	base.ControllerBase
}

func (r *AdminConPermitImportImpl) UserPageImport(c *gin.Context) {
	var (
		res *wrapper_admin.ResultPageImport
		arg wrapper_admin.ArgPageImport
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(ctx).
		UserPageImport(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}


func (r *AdminConPermitImportImpl) MenuImportSet(c *gin.Context) {
	var (
		res *wrappers.ResultMenuImportSet
		arg wrappers.ArgMenuImportSet
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitGroupImpl(ctx).
		MenuImportSet(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func (r *AdminConPermitImportImpl) DeleteImport(c *gin.Context) {
	var (
		err error
		arg wrappers.ArgDeleteImport
		res *wrappers.ResultDeleteImport
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeShouldBindUri); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(ctx).
		DeleteImport(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *AdminConPermitImportImpl) UpdateImportValue(c *gin.Context) {
	var (
		arg wrappers.ArgUpdateImportValue
		err error
		res *wrappers.ResultUpdateImportValue
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(ctx).
		UpdateImportValue(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *AdminConPermitImportImpl) ImportList(c *gin.Context) {
	var (
		arg wrappers.ArgImportList
		err error
		res *wrappers.ResultImportList
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeShouldBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(ctx).
		ImportList(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *AdminConPermitImportImpl) EditImport(c *gin.Context) {
	var (
		arg wrappers.ArgEditImport
		err error
		res *wrappers.ResultEditImport
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg, base.ControllerGetParamTypeBind); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(ctx).
		EditImport(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}

func (r *AdminConPermitImportImpl) GetImport(c *gin.Context) {
	var (
		res *wrappers.ResultGetImport
		arg wrappers.ArgGetImport
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(ctx).
		GetImport(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}


// GetImportByMenuId 根据菜单号 获取页面的接口ID
//
func (r *AdminConPermitImportImpl) GetImportByMenuId(c *gin.Context) {
	var (
		arg wrappers.ArgGetImportByMenuId
		res = wrappers.ResultGetImportByMenuId{}
		err error
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(ctx).
		GetImportByMenuId(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func (r *AdminConPermitImportImpl) MenuImport(c *gin.Context) {
	var (
		arg wrapper_admin.ArgMenuImport
		err error
		res *wrapper_admin.ResultMenuImport
		ctx = base.CreateContext(&r.ControllerBase, c)
	)
	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(ctx).
		MenuImport(&arg); err != nil {
		r.ResponseError(c, err)
		return
	} else {
		r.Response(c, base.SuccessCode, res)
	}
}
func NewAdminConPermitImportImpl() admins.AdminConPermitImport {
	controller := &AdminConPermitImportImpl{}
	controller.ControllerBase.Init()
	return controller
}
