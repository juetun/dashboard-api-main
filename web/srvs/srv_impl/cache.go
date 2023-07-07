package srv_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type (
	SrvCacheImpl struct {
		base.ServiceBase
		dao daos.DaoCacheParam
	}
	orgHandler func(headerInfo *common.HeaderInfo, products []*models.CacheKeyData) (res interface{}, err error)
)

func (r *SrvCacheImpl) ReloadAppCacheConfig(arg *wrapper_admin.ArgReloadAppCacheConfig) (res *wrapper_admin.ResultReloadAppCacheConfig, err error) {
	res = &wrapper_admin.ResultReloadAppCacheConfig{}
	res.Result = true
	return
}

func (r *SrvCacheImpl) CacheParamList(arg *wrapper_admin.ArgCacheParamList) (res *wrapper_admin.ResultCacheParamList, err error) {
	res = &wrapper_admin.ResultCacheParamList{Pager: response.NewPager(response.PagerBaseQuery(&arg.PageQuery)),}
	switch arg.PageType {
	case response.DefaultPageTypeList:
		err = r.getDataAsList(arg, res, r.orgHandler)
	case response.DefaultPageTypeNext:
		err = r.getDataAsNext(arg, res, r.orgHandler)
	default:
		err = fmt.Errorf("不支持你选择的分页类型")
		return
	}
	return
}

func (r *SrvCacheImpl) getDataAsNext(arg *wrapper_admin.ArgCacheParamList, res *wrapper_admin.ResultCacheParamList, orgHandler orgHandler) (err error) {
	var list []*models.CacheKeyData
	if list, err = r.dao.GetListAsNext(arg, res); err != nil {
		return
	}
	if res.Pager.List, err = orgHandler(&arg.HeaderInfo, list); err != nil {
		return
	}
	return
}

func (r *SrvCacheImpl) getDataAsList(arg *wrapper_admin.ArgCacheParamList, res *wrapper_admin.ResultCacheParamList, orgHandler orgHandler) (err error) {

	var actReturn *base.ActErrorHandlerResult

	err = res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		actReturn, err = r.dao.GetCount(arg, res)
		return
	}, func(pagerObject *response.Pager) (err error) {

		var list []*models.CacheKeyData
		if list, err = r.dao.GetList(actReturn, arg); err != nil {
			return
		}

		if res.Pager.List, err = orgHandler(&arg.HeaderInfo, list); err != nil {
			return
		}

		return
	})
	return
}

func (r *SrvCacheImpl) orgHandler(headerInfo *common.HeaderInfo, products []*models.CacheKeyData) (res interface{}, err error) {
	resTmp := make([]*wrapper_admin.ResultCacheParamItem, 0, len(products))
	defer func() {
		res = resTmp
	}()
	var (
		dt *wrapper_admin.ResultCacheParamItem
	)

	for _, value := range products {
		dt = &wrapper_admin.ResultCacheParamItem{}
		dt.ParseWithCacheParamData(value)
		resTmp = append(resTmp, dt)
	}
	return
}

func (r *SrvCacheImpl) ClearCache(arg *wrapper_admin.ArgClearCache) (res *wrapper_admin.ResultClearCache, err error) {
	res = &wrapper_admin.ResultClearCache{}

	res.Result = true
	return
}

func NewSrvCache(context ...*base.Context) srvs.SrvCache {
	p := &SrvCacheImpl{}
	p.SetContext(context...)
	if p.dao == nil {
		p.dao = dao_impl.NewDaoCacheParam(p.Context)
	}
	return p
}
