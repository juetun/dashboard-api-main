package dao_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/common/redis_pkg"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoGatewayPermitImpl struct {
	base.ServiceDao
}

func (r *DaoGatewayPermitImpl) GetCacheKeyImportWithAppKey(appName string) (key string, duration time.Duration, err error) {
	var CacheKeyImportWithAppKey *redis_pkg.CacheProperty
	if CacheKeyImportWithAppKey, err = parameters.GetCacheParamConfig("CacheKeyImportWithAppKey"); err != nil {
		return
	}
	key = fmt.Sprintf(CacheKeyImportWithAppKey.Key, appName)
	duration = CacheKeyImportWithAppKey.Expire
	return
}

func (r *DaoGatewayPermitImpl) SetCacheKeyImportWithApp(appName string, data interface{}) (err error) {
	var k string
	var duration time.Duration
	if k, duration, err = r.GetCacheKeyImportWithAppKey(appName); err != nil {
		return
	}
	if err = r.Context.CacheClient.
		Set(r.Context.GinContext.Request.Context(), k, data, duration).
		Err(); err != nil {
		r.Context.Error(map[string]interface{}{
			"appName":   appName,
			"err":       err.Error(),
			"cacheName": r.Context.CacheName,
		}, "DaoGatewayPermitImplSetCacheKeyImportWithApp")
	}
	return
}

func (r *DaoGatewayPermitImpl) GetCacheKeyImportWithApp(appName string, data interface{}) (isNil bool, err error) {
	var k string
	if k, _, err = r.GetCacheKeyImportWithAppKey(appName); err != nil {
		return
	}
	var e error
	isNil = true
	if e = r.Context.CacheClient.
		Get(r.Context.GinContext.Request.Context(), k).
		Scan(data); e != nil {
		if e.Error() == redis.Nil.Error() {
			return
		}
		r.Context.Error(map[string]interface{}{
			"appName":   appName,
			"err":       e.Error(),
			"cacheName": r.Context.CacheName,
		}, "DaoGatewayPermitImplGetCacheKeyImportWithApp")
		err = base.NewErrorRuntime(e, base.ErrorRedisCode)
		return
	}
	isNil = false
	return
}
func (r *DaoGatewayPermitImpl) GetImportListByAppName(appName string, refreshCache ...bool) (res map[int64]models.AdminImport, err error) {
	res = make(map[int64]models.AdminImport, 50)
	var isNil bool
	var f bool
	if len(refreshCache) > 0 {
		f = refreshCache[0]
	}
	if !f { // 如果不是强制刷新缓存,则优先从缓存读取
		if isNil, err = r.GetCacheKeyImportWithApp(appName, &res); err != nil {
			return
		}
		if !isNil {
			return
		}
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{"appName": appName, "err": err.Error(),}, "DaoGatewayPermitImplGetImportListByAppName")
	}()
	var dt []models.AdminImport
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(m.TableName()).
			Scopes(base.ScopesDeletedAt()).
			Where("app_name = ?", appName).
			Find(&dt).
			Error
		return
	})
	if err != nil {
		return
	}
	for _, adminImport := range dt {
		res[adminImport.Id] = adminImport
	}
	// 更新缓存中的数据
	_ = r.SetCacheKeyImportWithApp(appName, &res)
	return
}

func NewDaoGatewayPermit(ctx ...*base.Context) daos.DaoGatewayPermit {
	p := &DaoGatewayPermitImpl{}
	p.SetContext(ctx...)
	return p
}
