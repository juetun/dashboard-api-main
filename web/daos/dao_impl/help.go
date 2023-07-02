package dao_impl

import (
	"context"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/base/cache_act"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl/cache_act_local"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"time"
)

type DaoHelpImpl struct {
	base.ServiceDao
	ctx context.Context
}

func (r *DaoHelpImpl) getActErrorHandlerResultByArg(arg *wrapper_admin.ArgHelpList) (actRes *base.ActErrorHandlerResult) {
	var m *models.HelpDocument
	actRes = r.GetDefaultActErrorHandlerResult(m)
	actRes.Db = actRes.Db.Table(actRes.TableName).
		Scopes(base.ScopesDeletedAt())
	if arg.PKey != "" {
		actRes.Db = actRes.Db.Where("p_key = ?", arg.PKey)
	}
	return
}

func (r *DaoHelpImpl) GetCountByArg(arg *wrapper_admin.ArgHelpList, pager *response.Pager) (actRes *base.ActErrorHandlerResult, err error) {

	actRes = r.getActErrorHandlerResultByArg(arg)
	err = r.ActErrorHandler(func() (ar *base.ActErrorHandlerResult) {
		ar = actRes
		ar.Err = ar.Db.Count(&pager.TotalCount).Error
		return
	})
	return
}

func (r *DaoHelpImpl) GetListByArg(arg *wrapper_admin.ArgHelpList, actRes *base.ActErrorHandlerResult, pagerObject *response.Pager) (res []*models.HelpDocument, err error) {
	if actRes == nil {
		actRes = r.getActErrorHandlerResultByArg(arg)
	}
	actRes.Err = actRes.Db.
		Limit(pagerObject.PageSize).
		Offset(pagerObject.GetOffset()).
		Find(&res).
		Error
	return
}

func (r *DaoHelpImpl) GetByIds(arg *base.ArgGetByNumberIds) (res map[int64]*models.HelpDocument, err error) {

	res, err = cache_act_local.NewCacheHelpDocWithIdsAction(
		cache_act_local.CacheHelpDocWithIdsActionArg(arg),
		cache_act_local.CacheHelpDocWithIdsActionGetByIdsFromDb(r.getWithIdsByIdsFromDb),
	).LoadBaseOption(
		cache_act.CacheActionBaseGetCacheKey(r.getWithIdsCacheKey),
		cache_act.CacheActionBaseContext(r.Context),
		cache_act.CacheActionBaseCtx(r.ctx), ).
		Action()
	return
}

func (r *DaoHelpImpl) getExtByIdsFromDb(pks ...string) (res map[string]*models.HelpDocument, err error) {
	l := len(pks)
	res = make(map[string]*models.HelpDocument, l)
	if l == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"pks": pks,
			"err": err.Error(),
		}, "DaoHelpImplGetExtByIdsFromDb")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var data *models.HelpDocument
		var list []*models.HelpDocument
		actRes = r.GetDefaultActErrorHandlerResult(data)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("p_key IN (?)", pks).Find(&list).
			Error
		if actRes.Err != nil {
			return
		}
		for _, ext := range list {
			res[ext.PKey] = ext
		}
		return
	})

	return
}

func (r *DaoHelpImpl) getWithIdsByIdsFromDb(ids ...int64) (res map[int64]*models.HelpDocument, err error) {
	l := len(ids)
	res = make(map[int64]*models.HelpDocument, l)
	if l == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"ids": ids,
			"err": err.Error(),
		}, "DaoHelpImplGetWithIdsByIdsFromDb")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var data *models.HelpDocument
		var list []*models.HelpDocument
		actRes = r.GetDefaultActErrorHandlerResult(data)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("id IN (?)", ids).Find(&list).
			Error
		if actRes.Err != nil {
			return
		}
		for _, ext := range list {
			res[ext.Id] = ext
		}
		return
	})

	return
}

func (r *DaoHelpImpl) GetByPKey(arg *base.ArgGetByStringIds) (res map[string]*models.HelpDocument, err error) {

	res, err = cache_act_local.NewCacheHelpDocAction(
		cache_act_local.CacheHelpDocActionArg(arg),
		cache_act_local.CacheHelpDocActionGetByIdsFromDb(r.getExtByIdsFromDb),
	).LoadBaseOption(
		cache_act.CacheActionBaseGetCacheKey(r.getExtCacheKey),
		cache_act.CacheActionBaseContext(r.Context),
		cache_act.CacheActionBaseCtx(r.ctx), ).
		Action()
	return

}

func (r *DaoHelpImpl) getExtCacheKey(id interface{}, expireTimeRands ...bool) (res string, timeExpire time.Duration) {

	res = fmt.Sprintf(parameters.CacheKeyHelp.Key, id)
	timeExpire = parameters.CacheKeyHelp.Expire
	var expireTimeRand bool
	if len(expireTimeRands) > 0 {
		expireTimeRand = expireTimeRands[0]
	}
	if expireTimeRand {
		randNumber, _ := r.RealRandomNumber(60) //打乱缓存时长，防止缓存同一时间过期导致数据库访问压力变大
		timeExpire = timeExpire + time.Duration(randNumber)*time.Second
	}

	return
}

func (r *DaoHelpImpl) getWithIdsCacheKey(id interface{}, expireTimeRands ...bool) (res string, timeExpire time.Duration) {

	res = fmt.Sprintf(parameters.CacheKeyHelpWithId.Key, id)
	timeExpire = parameters.CacheKeyHelpWithId.Expire
	var expireTimeRand bool
	if len(expireTimeRands) > 0 {
		expireTimeRand = expireTimeRands[0]
	}
	if expireTimeRand {
		randNumber, _ := r.RealRandomNumber(60) //打乱缓存时长，防止缓存同一时间过期导致数据库访问压力变大
		timeExpire = timeExpire + time.Duration(randNumber)*time.Second
	}

	return
}

func (r *DaoHelpImpl) AddHelpDocument(data *models.HelpDocument) (err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err.Error(),
		}, "DaoHelpImplAddHelpDocument")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(data)
		actRes.Err = r.AddOneData(actRes.ParseAddOneDataParameter(base.AddOneDataParameterModel(data)))
		return
	})
	return
}

func NewDaoHelp(c ...*base.Context) daos.DaoHelp {
	p := &DaoHelpImpl{}
	p.SetContext(c...)
	if p.ctx == nil {
		p.ctx = context.TODO()
	}
	return p
}
