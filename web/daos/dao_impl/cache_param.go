package dao_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type DaoCacheParamImpl struct {
	base.ServiceDao
}

func (r *DaoCacheParamImpl) GetListAsNext(arg *wrapper_admin.ArgCacheParamList, resData *wrapper_admin.ResultCacheParamList) (res []*models.CacheKeyData, err error) {

	res = []*models.CacheKeyData{}
	var actHandler = func() (actRes *base.ActErrorHandlerResult) {
		var returnNil bool
		actRes, returnNil = r.getFetchListDbWithArg(arg)
		if returnNil {
			return
		}
		if arg.RequestId != "" {
			actRes.Db = actRes.Db.Where("`id` < ?", arg.RequestId)
		}
		actRes.Err = actRes.Db.Limit(arg.PageSize).
			Order("`id` DESC").
			Find(&res).
			Error
		if actRes.Err != nil {
			return
		}

		var minId int64
		for k, item := range res {
			if k == 0 || item.Id < minId {
				minId = item.Id
			}
		}
		resData.Pager.RequestId = fmt.Sprintf("%d", minId)
		if len(res) == arg.PageSize {
			resData.Pager.IsNext = true
		}
		return
	}
	err = r.ActErrorHandler(actHandler)
	return
}

func (r *DaoCacheParamImpl) GetCount(arg *wrapper_admin.ArgCacheParamList, res *wrapper_admin.ResultCacheParamList) (actResObject *base.ActErrorHandlerResult, err error) {
	var returnNil bool
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes, returnNil = r.getFetchListDbWithArg(arg)
		actResObject = actRes
		if returnNil {
			return
		}
		var e error
		if e = actRes.Db.Count(&res.Pager.TotalCount).
			Error; e != nil {
			var m *models.CacheKeyData
			actRes = r.GetDefaultActErrorHandlerResult(m)
			actRes.Err = e
			return
		}
		return
	})
	return
}

func (r *DaoCacheParamImpl) getFetchListDbWithArg(arg *wrapper_admin.ArgCacheParamList) (actRes *base.ActErrorHandlerResult, returnNil bool) {
	var m *models.CacheKeyData
	actRes = r.GetDefaultActErrorHandlerResult(m)
	actRes.Db = actRes.Db.Table(actRes.TableName).Scopes(base.ScopesDeletedAt())

	if arg.MicroApp != "" {
		actRes.Db = actRes.Db.Where("micro_app = ?", arg.MicroApp)
	}
	return
}

func (r *DaoCacheParamImpl) GetList(actResObject *base.ActErrorHandlerResult, arg *wrapper_admin.ArgCacheParamList) (res []*models.CacheKeyData, err error) {
	var returnNil bool

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "DaoCacheParamGetList")
	}()

	if actResObject == nil {
		actResObject, returnNil = r.getFetchListDbWithArg(arg)
		if returnNil {
			return
		}
	}
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = actResObject
		actRes.Err = actResObject.Db.
			Offset(arg.GetOffset()).
			Limit(arg.PageSize).
			Find(&res).Error
		return
	})

	return
}

func NewDaoCacheParam(ctx ...*base.Context) (res daos.DaoCacheParam) {
	p := &DaoCacheParamImpl{}
	p.SetContext(ctx...)
	return p
}
