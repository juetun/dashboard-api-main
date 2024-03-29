package dao_impl

import (
	"context"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/base/cache_act"
	"github.com/juetun/base-wrapper/lib/common/redis_pkg"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl/cache_act_local"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"time"
)

type DaoHelpRelateImpl struct {
	base.ServiceDao
	ctx context.Context
}

func (r *DaoHelpRelateImpl) GetByIdFromDb(ids ...int64) (res []*models.HelpDocumentRelate, err error) {
	var l = len(ids)
	res = make([]*models.HelpDocumentRelate, 0, l)
	if l == 0 {
		return
	}
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.HelpDocumentRelate
		actRes = r.GetDefaultActErrorHandlerResult(m)
		if actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("id IN (?)", ids).
			Find(&res).Error; actRes.Err != nil {
			return
		}
		return
	})

	return
}

func (r *DaoHelpRelateImpl) getByDocKeyFromDb(keys ...string) (res map[string]*models.HelpDocumentRelate, err error) {
	var l = len(keys)
	res = make(map[string]*models.HelpDocumentRelate, l)
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"err": err.Error(),
		}, "DaoHelpRelateImplGetByDocKeyFromDb")
	}()
	if l == 0 {
		return
	}
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.HelpDocumentRelate
		var data = make([]*models.HelpDocumentRelate, 0, l)
		actRes = r.GetDefaultActErrorHandlerResult(m)
		if actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("doc_key IN (?)", keys).
			Find(&data).Error; actRes.Err != nil {
			return
		}
		for _, item := range data {
			res [item.DocKey] = item
		}
		return
	})

	return
}

func (r *DaoHelpRelateImpl) getByDocKeyCacheKey(id interface{}, expireTimeRands ...bool) (res string, timeExpire time.Duration, err error) {
	var CacheHelpDocRelateByKeyUpdating *redis_pkg.CacheProperty
	if CacheHelpDocRelateByKeyUpdating, err = parameters.GetCacheParamConfig("CacheHelpDocRelateByKeyUpdating"); err != nil {
		return
	}
	res = fmt.Sprintf(CacheHelpDocRelateByKeyUpdating.Key, id)
	timeExpire = CacheHelpDocRelateByKeyUpdating.Expire
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

func (r *DaoHelpRelateImpl) GetByDocKeys(arg *base.ArgGetByStringIds) (res map[string]*models.HelpDocumentRelate, err error) {
	res, err = cache_act_local.NewCacheHelpRelateAction(
		cache_act_local.CacheHelpRelateActionArg(arg),
		cache_act_local.CacheHelpRelateActionGetByIdsFromDb(r.getByDocKeyFromDb),
	).LoadBaseOption(
		cache_act.CacheActionBaseGetCacheKey(r.getByDocKeyCacheKey),
		cache_act.CacheActionBaseContext(r.Context),
		cache_act.CacheActionBaseCtx(r.ctx), ).
		Action()

	return
}

func (r *DaoHelpRelateImpl) UpdateById(id int64, data map[string]interface{}) (err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"err": err.Error(),
		}, "DaoHelpRelateImplUpdateById")
	}()

	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.HelpDocumentRelate
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("id = ?", id).
			Updates(data).Error
		return
	})
	if err != nil {
		return
	}
	return
}

func (r *DaoHelpRelateImpl) GetByTopHelp() (res []*models.HelpDocumentRelate, err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"err": err.Error(),
		}, "DaoHelpRelateImplGetByTopHelp")
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.HelpDocumentRelate
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("parent_id = ?", models.HelpDocumentRelateTopPid).
			Find(&res).Error
		return
	})
	return
}

func (r *DaoHelpRelateImpl) GetAllHelpRelate(arg *base.ArgGetByStringIds) (res models.HelpDocumentRelateCaches, err error) {
	res = models.HelpDocumentRelateCaches{}
	keyName := "GetAllHelpRelate"
	arg.Ids = []string{keyName}
	var (
		resData map[string]models.HelpDocumentRelateCaches
	)
	if resData, err = cache_act_local.NewCacheAllHelpRelateAction(
		cache_act_local.CacheAllHelpRelateActionArg(arg),
		cache_act_local.CacheAllHelpRelateActionGetByIdsFromDb(r.getAllHelpRelateFromDb),
	).LoadBaseOption(
		cache_act.CacheActionBaseGetCacheKey(r.getAllHelpRelateCacheKey),
		cache_act.CacheActionBaseContext(r.Context),
		cache_act.CacheActionBaseCtx(r.ctx), ).
		Action(); err != nil {
		return
	}
	if tmp, ok := resData[keyName]; ok {
		res = tmp
	}
	return
}

func (r *DaoHelpRelateImpl) getAllHelpRelateCacheKey(id interface{}, expireTimeRands ...bool) (res string, timeExpire time.Duration, err error) {
	var CacheKeyAllHelpRelate *redis_pkg.CacheProperty
	if CacheKeyAllHelpRelate, err = parameters.GetCacheParamConfig("CacheKeyAllHelpRelate"); err != nil {
		return
	}
	res = fmt.Sprintf(CacheKeyAllHelpRelate.Key, id)
	timeExpire = CacheKeyAllHelpRelate.Expire
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

func (r *DaoHelpRelateImpl) getAllHelpRelateFromDb(ids ...string) (res map[string]models.HelpDocumentRelateCaches, err error) {
	keyName := "GetAllHelpRelate"
	res = make(map[string]models.HelpDocumentRelateCaches, len(ids))
	var list []*models.HelpDocumentRelateCache
	var actHandler = func() (actRes *base.ActErrorHandlerResult) {
		var m *models.HelpDocumentRelate
		actRes = r.GetDefaultActErrorHandlerResult(m)
		if actRes.Err = actRes.Db.Select("id,biz_code,display,parent_id,label,is_leaf_node,doc_key").
			Table(actRes.TableName).
			Where("display = ?", models.DisplayYes).
			Scopes(base.ScopesDeletedAt()).
			Find(&list).
			Error; actRes.Err != nil {
			return
		}
		res[keyName] = list
		return
	}
	if err = r.ActErrorHandler(actHandler); err != nil {
		return
	}
	return
}

func (r *DaoHelpRelateImpl) GetByTopId(bizCode string, topIds ...int64) (res map[int64][]*wrapper_admin.ResultHelpTreeItem, err error) {
	if len(topIds) == 0 || bizCode == "" {
		return
	}
	var list []*models.HelpDocumentRelate
	var m *models.HelpDocumentRelate
	var actResParam = r.GetDefaultActErrorHandlerResult(m)

	var actHandler = func() (actRes *base.ActErrorHandlerResult) {
		actRes = actResParam
		actRes.Err = actRes.Db.
			Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("id IN (?) AND biz_code = ?", topIds, bizCode).
			Find(&list).
			Error
		if actRes.Err != nil {
			return
		}
		return
	}
	if err = r.ActErrorHandler(actHandler); err != nil {
		return
	}
	if res, err = r.orgGroupRelate(bizCode, actResParam, topIds...); err != nil {
		return
	}
	return
}

func (r *DaoHelpRelateImpl) getChildByTopIds(bizCode string, actRes *base.ActErrorHandlerResult, topIds ...int64) (dataMap map[int64][]*wrapper_admin.ResultHelpTreeItem, childTopIds []int64, err error) {
	if len(topIds) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"topIds": topIds,
			"err":    err.Error(),
		}, "DaoHelpRelateImplGetChildByTopIds")
	}()
	var list []*models.HelpDocumentRelate
	if err = actRes.Db.Table(actRes.TableName).
		Where("parent_id IN(?) AND biz_code = ?", topIds, bizCode).
		Scopes(base.ScopesDeletedAt()).
		Find(&list).
		Error; err != nil {
		return
	}
	dataMap, childTopIds = r.groupRelate(list)
	return
}

func (r *DaoHelpRelateImpl) orgGroupRelate(bizCode string, actRes *base.ActErrorHandlerResult, topIds ...int64) (res map[int64][]*wrapper_admin.ResultHelpTreeItem, err error) {
	var (
		childTopIds []int64
	)
	if res, childTopIds, err = r.getChildByTopIds(bizCode, actRes, topIds...); err != nil {
		return
	}
	if len(childTopIds) > 0 {
		var resChildren map[int64][]*wrapper_admin.ResultHelpTreeItem
		//获得儿子的数据
		if resChildren, err = r.orgGroupRelate(bizCode, actRes, childTopIds...); err != nil {
			return
		}
		for _, item := range res {
			for _, it := range item {
				if dt, ok := resChildren[it.Id]; ok {
					it.Child = dt
				}
			}
		}
		return
	}

	return
}

//组织数据使用
func (r *DaoHelpRelateImpl) groupRelate(list []*models.HelpDocumentRelate) (dataMap map[int64][]*wrapper_admin.ResultHelpTreeItem, childTopIds []int64) {
	dataMap = map[int64][]*wrapper_admin.ResultHelpTreeItem{}
	var (
		l     = len(list)
		dtRes *wrapper_admin.ResultHelpTreeItem
	)
	for _, item := range list {
		if _, ok := dataMap[item.ParentId]; !ok {
			dataMap[item.ParentId] = make([]*wrapper_admin.ResultHelpTreeItem, 0, l)
		}
		dtRes = &wrapper_admin.ResultHelpTreeItem{}
		dtRes.HelpDocumentRelate = *item
		//如果不是叶子节点
		if item.IsLeafNode == models.HelpDocumentRelateIsLeafNodeNo && item.Id != 0 {
			childTopIds = append(childTopIds, item.Id)
		}
		dataMap[item.ParentId] = append(dataMap[item.ParentId], dtRes)
	}
	return
}

func (r *DaoHelpRelateImpl) AddOneHelpRelate(relate *models.HelpDocumentRelate) (err error) {
	if relate == nil {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"relate": relate,
			"err":    err.Error(),
		}, "DaoHelpRelateImplAddOneHelpRelate")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(relate)
		actRes.Err = r.AddOneData(actRes.ParseAddOneDataParameter(base.AddOneDataParameterModel(relate)))
		return
	})
	return
}

func NewDaoHelpRelate(c ...*base.Context) daos.DaoHelpRelate {
	p := &DaoHelpRelateImpl{}
	p.SetContext(c...)
	if p.ctx == nil {
		p.ctx = context.TODO()
	}
	return p
}
