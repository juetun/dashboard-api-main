package cache_act_local

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/base/cache_act"
	"github.com/juetun/dashboard-api-main/web/models"
)

type (
	CacheHelpRelateAction struct {
		cache_act.CacheActionBase
		arg            *base.ArgGetByStringIds
		GetByIdsFromDb GetHelpRelateByIdsFromDb
	}
	GetHelpRelateByIdsFromDb    func(id ...string) (resData map[string]*models.HelpDocumentRelate, err error)
	CacheHelpRelateActionOption func(cacheFreightAction *CacheHelpRelateAction)
)

func CacheHelpRelateActionArg(arg *base.ArgGetByStringIds) CacheHelpRelateActionOption {
	return func(cacheFreightAction *CacheHelpRelateAction) {
		cacheFreightAction.arg = arg
		return
	}
}

func CacheHelpRelateActionGetByIdsFromDb(getByIdsFromDb GetHelpRelateByIdsFromDb) CacheHelpRelateActionOption {
	return func(cacheFreightAction *CacheHelpRelateAction) {
		cacheFreightAction.GetByIdsFromDb = getByIdsFromDb
		return
	}
}

func NewCacheHelpRelateAction(options ...CacheHelpRelateActionOption) (res *CacheHelpRelateAction) {
	res = &CacheHelpRelateAction{CacheActionBase: cache_act.NewCacheActionBase()}
	for _, handler := range options {
		handler(res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	return
}

func (r *CacheHelpRelateAction) LoadBaseOption(options ...cache_act.CacheActionBaseOption) *CacheHelpRelateAction {
	r.LoadBase(options...)
	return r
}

func (r *CacheHelpRelateAction) Action() (res map[string]*models.HelpDocumentRelate, err error) {
	res = map[string]*models.HelpDocumentRelate{}
	if len(r.arg.Ids) == 0 {
		return
	}
	// 初始化 type默认值
	if err = r.arg.Default(); err != nil {
		return
	}
	switch r.arg.GetType {
	case base.GetDataTypeFromDb: // 从数据库获取数据
		if res, err = r.GetByIdsFromDb(r.arg.Ids...); err != nil {
			return
		}
		switch r.arg.RefreshCache {
		case base.RefreshCacheYes:
			for id, value := range res {
				if err = r.SetToCache(id, value); err != nil {
					return
				}
			}
		}
	case base.GetDataTypeFromCache: // 从缓存获取数据
		res, _, err = r.getByIdsFromCache(r.arg.Ids...)
	case base.GetDataTypeFromAll: // 优先从缓存获取，如果没有数据，则从数据库获取
		res, err = r.getByIdsFromAll(r.arg.Ids...)
	default:
		err = fmt.Errorf("当前不支持你选择的获取数据类型(%s)", r.arg.GetType)
	}
	return
}

func (r *CacheHelpRelateAction) getFromCache(id interface{}) (data *models.HelpDocumentRelate, err error) {
	defer func() {
		if err != nil && err != redis.Nil {
			r.Context.Info(map[string]interface{}{
				"productId": id,
				"err":       err.Error(),
			}, "CacheHelpRelateActionGetFromCache")
			err = base.NewErrorRuntime(err, base.ErrorRedisCode)
			return
		}
	}()
	var key string
	if key, _, err = r.GetCacheKey(id); err != nil {
		return
	}
	cmd := r.Context.CacheClient.Get(r.Ctx, key)
	if err = cmd.Err(); err != nil {
		return
	}

	if err = cmd.Scan(data); err != nil {
		return
	}
	return
}

func (r *CacheHelpRelateAction) getByIdsFromCache(ids ...string) (res map[string]*models.HelpDocumentRelate, noCacheIds []string, err error) {
	var e error

	res = map[string]*models.HelpDocumentRelate{}

	//收集缓存中没有的数据ID，便于后边查询使用
	noCacheIds = make([]string, 0, len(ids))

	for _, id := range ids {
		if res[id], e = r.getFromCache(id); e != nil {
			if e != redis.Nil {
				err = e
				return
			}
			noCacheIds = append(noCacheIds, id)
			res[id] = &models.HelpDocumentRelate{}
		}
	}
	return
}

func (r *CacheHelpRelateAction) getByIdsFromAll(ids ...string) (res map[string]*models.HelpDocumentRelate, err error) {
	var pIds []string
	if res, pIds, err = r.getByIdsFromCache(ids...); err != nil {
		return
	}
	if len(pIds) == 0 {
		return
	}

	var dt map[string]*models.HelpDocumentRelate
	if dt, err = r.GetByIdsFromDb(pIds...); err != nil {
		return
	}

	for _, pid := range pIds {
		if dta, ok := dt[pid]; ok {
			res[pid] = dta
			_ = r.SetToCache(pid, dta)
			continue
		}
		_ = r.SetToCache(pid, &models.HelpDocumentRelate{})
	}

	return
}
