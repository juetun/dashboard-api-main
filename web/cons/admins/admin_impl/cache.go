package admin_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type ConAdminCacheImpl struct {
	base.ControllerBase
}

func (r *ConAdminCacheImpl) ReloadAppCacheConfig(c *gin.Context) {

	var (
		err  error
		arg  wrapper_admin.ArgReloadAppCacheConfig
		data = &wrapper_admin.ResultReloadAppCacheConfig{}
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	if data, err = srv_impl.NewSrvCache(ctx).
		ReloadAppCacheConfig(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, data)

	return
}

func (r *ConAdminCacheImpl) CacheParamList(c *gin.Context) {
	var (
		err  error
		arg  wrapper_admin.ArgCacheParamList
		data = &wrapper_admin.ResultCacheParamList{}
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	if data, err = srv_impl.NewSrvCache(ctx).
		CacheParamList(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, data)

	return
}

func (r *ConAdminCacheImpl) ClearCache(c *gin.Context) {

	var (
		err  error
		arg  wrapper_admin.ArgClearCache
		data = &wrapper_admin.ResultClearCache{}
		ctx  = base.CreateContext(&r.ControllerBase, c)
	)

	if haveErr := r.ParametersAccept(ctx, &arg); haveErr {
		return
	}

	if data, err = srv_impl.NewSrvCache(ctx).
		ClearCache(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, data)
	return
}

func NewConAdminCache() admins.ConAdminCache {
	p := &ConAdminCacheImpl{}
	p.ControllerBase.Init()
	return p
}
