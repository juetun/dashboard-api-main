package cache_act_local

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/base/cache_act"
	"github.com/juetun/dashboard-api-main/web/models"

	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type (
	CacheMenuImportAction struct {
		cache_act.CacheActionBase
		arg            *base.ArgGetByNumberIds
		GetByIdsFromDb GetMenuImportByIdsFromDb
	}
	GetMenuImportByIdsFromDb    func(id ...int64) (resData map[int64]models.AdminMenuImportCache, err error)
	CacheMenuImportActionOption func(cacheFreightAction *CacheMenuImportAction)
)

func CacheMenuImportActionArg(arg *base.ArgGetByNumberIds) CacheMenuImportActionOption {
	return func(cacheFreightAction *CacheMenuImportAction) {
		cacheFreightAction.arg = arg
		return
	}
}

func CacheMenuImportActionGetByIdsFromDb(getByIdsFromDb GetMenuImportByIdsFromDb) CacheMenuImportActionOption {
	return func(cacheFreightAction *CacheMenuImportAction) {
		cacheFreightAction.GetByIdsFromDb = getByIdsFromDb
		return
	}
}

func NewMenuImportAction(options ...CacheMenuImportActionOption) (res *CacheMenuImportAction) {
	res = &CacheMenuImportAction{CacheActionBase: cache_act.NewCacheActionBase()}
	for _, handler := range options {
		handler(res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	return
}

func (r *CacheMenuImportAction) LoadBaseOption(options ...cache_act.CacheActionBaseOption) *CacheMenuImportAction {
	r.LoadBase(options...)
	return r
}

func (r *CacheMenuImportAction) Action() (res map[int64]models.AdminMenuImportCache, err error) {
	res = map[int64]models.AdminMenuImportCache{}
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
				if err = r.SetToCache(id, &value); err != nil {
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

func (r *CacheMenuImportAction) getFromCache(id interface{}) (data models.AdminMenuImportCache, err error) {
	data = models.AdminMenuImportCache{}
	defer func() {
		if err != nil && err != redis.Nil {
			r.Context.Info(map[string]interface{}{
				"productId": id,
				"err":       err.Error(),
			}, "CacheMenuImportActionGetFromCache")
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

	if errString := cmd.Scan(&data).Error(); errString != "" {
		err = fmt.Errorf(errString)
		return
	}
	return
}

func (r *CacheMenuImportAction) getByIdsFromCache(ids ...int64) (res map[int64]models.AdminMenuImportCache, noCacheIds []int64, err error) {
	var e error

	res = map[int64]models.AdminMenuImportCache{}

	//收集缓存中没有的数据ID，便于后边查询使用
	noCacheIds = make([]int64, 0, len(ids))

	for _, id := range ids {
		if res[id], e = r.getFromCache(id); e != nil {
			if e != redis.Nil {
				err = e
				return
			}
			noCacheIds = append(noCacheIds, id)
			res[id] = models.AdminMenuImportCache{}
		}
	}
	return
}

func (r *CacheMenuImportAction) getByIdsFromAll(ids ...int64) (res map[int64]models.AdminMenuImportCache, err error) {
	var pIds []int64
	if res, pIds, err = r.getByIdsFromCache(ids...); err != nil {
		return
	}
	if len(pIds) == 0 {
		return
	}

	var dt map[int64]models.AdminMenuImportCache
	if dt, err = r.GetByIdsFromDb(pIds...); err != nil {
		return
	}

	for _, pid := range pIds {
		if dta, ok := dt[pid]; ok {
			res[pid] = dta
			_ = r.SetToCache(pid, &dta)
			continue
		}
		_ = r.SetToCache(pid, &[]models.AdminMenuImport{})
	}

	return
}
