package srv_impl

import (
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/library/common/app_param"
)

type SrvCacheOperateImpl struct {
	base.ServiceBase
}

func (r *SrvCacheOperateImpl) ClearCacheByKeyPrefix(arg *app_param.ArgClearCacheByKeyPrefix) (res *app_param.ResultClearCacheByKeyPrefix, err error) {
	res = &app_param.ResultClearCacheByKeyPrefix{}
	res.Resutlt = true
	return
}

func (r *SrvCacheOperateImpl) GetCacheParamConfig(arg *app_param.ArgGetCacheParamConfig) (res app_param.ResultGetCacheParamConfig, err error) {
	res = parameters.CacheParamConfig
	for _, item := range res {
		if item.MicroApp == "" {
			item.MicroApp = app_obj.App.AppName
		}
	}
	return
}

func NewSrvCacheOperate(ctx ...*base.Context) srvs.SrvCacheOperate {
	p := &SrvCacheOperateImpl{}
	p.SetContext(ctx...)
	return p
}
