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
	CacheHelpDocAction struct {
		cache_act.CacheActionBase
		arg            *base.ArgGetByStringIds
		GetByIdsFromDb GetHelpDocByIdsFromDb
	}
	GetHelpDocByIdsFromDb    func(id ...string) (resData map[string]*models.HelpDocument, err error)
	CacheHelpDocActionOption func(cacheFreightAction *CacheHelpDocAction)
)

func CacheHelpDocActionArg(arg *base.ArgGetByStringIds) CacheHelpDocActionOption {
	return func(cacheFreightAction *CacheHelpDocAction) {
		cacheFreightAction.arg = arg
		return
	}
}

func CacheHelpDocActionGetByIdsFromDb(getByIdsFromDb GetHelpDocByIdsFromDb) CacheHelpDocActionOption {
	return func(cacheFreightAction *CacheHelpDocAction) {
		cacheFreightAction.GetByIdsFromDb = getByIdsFromDb
		return
	}
}

func NewCacheHelpDocAction(options ...CacheHelpDocActionOption) (res *CacheHelpDocAction) {
	res = &CacheHelpDocAction{CacheActionBase: cache_act.NewCacheActionBase()}
	for _, handler := range options {
		handler(res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	return
}

func (r *CacheHelpDocAction) LoadBaseOption(options ...cache_act.CacheActionBaseOption) *CacheHelpDocAction {
	r.LoadBase(options...)
	return r
}

func (r *CacheHelpDocAction) Action() (res map[string]*models.HelpDocument, err error) {
	res = map[string]*models.HelpDocument{}
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

func (r *CacheHelpDocAction) getFromCache(id interface{}) (data *models.HelpDocument, err error) {
	data = &models.HelpDocument{}
	defer func() {
		if err != nil && err != redis.Nil {
			r.Context.Info(map[string]interface{}{
				"productId": id,
				"err":       err.Error(),
			}, "CacheActionGetFromCache")
			err = base.NewErrorRuntime(err, base.ErrorRedisCode)
			return
		}
	}()
	key, _ := r.GetCacheKey(id)
	cmd := r.Context.CacheClient.Get(r.Ctx, key)
	if err = cmd.Err(); err != nil {
		return
	}

	if err = cmd.Scan(data); err != nil {
		return
	}
	return
}

func (r *CacheHelpDocAction) getByIdsFromCache(ids ...string) (res map[string]*models.HelpDocument, noCacheIds []string, err error) {
	var e error

	res = map[string]*models.HelpDocument{}

	//收集缓存中没有的数据ID，便于后边查询使用
	noCacheIds = make([]string, 0, len(ids))

	for _, id := range ids {
		if res[id], e = r.getFromCache(id); e != nil {
			if e != redis.Nil {
				err = e
				return
			}
			noCacheIds = append(noCacheIds, id)
			res[id] = &models.HelpDocument{}
		}
	}
	return
}

func (r *CacheHelpDocAction) getByIdsFromAll(ids ...string) (res map[string]*models.HelpDocument, err error) {
	var pIds []string
	res, pIds, err = r.getByIdsFromCache(ids...)
	if len(pIds) == 0 {
		return
	}

	var dt map[string]*models.HelpDocument
	if dt, err = r.GetByIdsFromDb(pIds...); err != nil {
		return
	}

	for _, pid := range pIds {
		if dta, ok := dt[pid]; ok {
			res[pid] = dta
			_ = r.SetToCache(pid, dta)
			continue
		}
		_ = r.SetToCache(pid, &models.HelpDocument{})
	}

	return
}
