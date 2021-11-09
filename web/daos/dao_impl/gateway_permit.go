package dao_impl

import (
	"fmt"
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

func (r *DaoGatewayPermitImpl) GetCacheKeyImportWithAppKey(appName string) (key string, duration time.Duration) {
	key = fmt.Sprintf(parameters.CacheKeyImportWithAppKey, appName)
	duration = parameters.CacheKeyImportWithAppKeyTime
	return
}

func (r *DaoGatewayPermitImpl) SetCacheKeyImportWithApp(appName string, data interface{}) (err error) {
	k, duration := r.GetCacheKeyImportWithAppKey(appName)
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
	k, _ := r.GetCacheKeyImportWithAppKey(appName)
	var e error
	if e = r.Context.CacheClient.
		Get(r.Context.GinContext.Request.Context(), k).
		Scan(data); e != nil {
		if e == redis.Nil {
			isNil = true
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
	return
}
func (r *DaoGatewayPermitImpl) GetImportListByAppName(appName string, refreshCache bool) (res map[int]models.AdminImport, err error) {
	res = make(map[int]models.AdminImport, 50)
	var isNil bool
	if !refreshCache { // 如果不是强制刷新缓存,则优先从缓存读取
		if isNil, err = r.GetCacheKeyImportWithApp(appName, &res); err != nil {
			return
		}
		if !isNil {
			return
		}
	}
	var m models.AdminImport
	var dt []models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where("app_name = ?", appName).
		Find(&dt).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"appName": appName,
		}, "DaoGatewayPermitImplGetImportListByAppName")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
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
