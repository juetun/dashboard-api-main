// Package srv_impl
/**
* Created by GoLand.
* Date: 2019-05-07
* Time: 22:12
 */
package srv_impl

import (
	"encoding/json"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"

	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"

	"time"
)

type SystemService struct {
	base.ServiceBase
}

func NewSystemService(context ...*base.Context) (p *SystemService) {
	p = &SystemService{}
	p.SetContext(context...)
	return
}
func (r *SystemService) GetSystemList() (system *models.ZBaseSys, err error) {
	system = new(models.ZBaseSys)
	var list []*models.ZBaseSys
	if list, err = dao_impl.NewDaoSystem(r.Context).
		GetSystemList(); err != nil {
		return
	}
	if len(list) > 0 {
		system = list[0]
	}
	var e error
	e = r.Context.Db.
		Table(system.TableName()).
		Find(system).
		Error
	if e != nil && e != gorm.ErrRecordNotFound {
		err = e
		r.Context.Error(
			map[string]interface{}{
				"message": "service.GetSystemList",
				"err":     e.Error(),
			},
		)

		return
	}
	return
}

func (r *SystemService) SystemUpdate(sId int, ss wrappers.ConsoleSystem) (err error) {
	systemUpdate := models.ZBaseSys{
		Title:        ss.Title,
		Keywords:     ss.Keywords,
		Description:  ss.Description,
		RecordNumber: ss.RecordNumber,
		Theme:        ss.Theme,
	}
	if sId == 0 {
		err = r.Context.Db.Table((&models.ZBaseSys{}).TableName()).Create(&systemUpdate).Error
		return
	}
	err = r.Context.Db.Table((&models.ZBaseSys{}).TableName()).Where("id=?", sId).
		Updates(&systemUpdate).
		Error
	return err
}

func (r *SystemService) IndexSystem() (system *models.ZBaseSys, err error) {
	cacheKey := common.Conf.SystemIndexKey
	cacheRes, err := r.Context.CacheClient.Get(r.Context.GinContext.Request.Context(), cacheKey).Result()
	if err.Error() == redis.Nil.Error() {
		system, err := r.doCacheIndexSystem(cacheKey)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{
					"message": "service.IndexSystem1",
					"err":     err.Error(),
				},
			)
			return system, err
		}
		return system, nil
	}
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.IndexSystem2",
				"err":     err.Error(),
			},
		)
		return system, err
	}

	err = json.Unmarshal([]byte(cacheRes), &system)
	if err == nil {
		return system, nil
	}
	r.Context.Info(
		map[string]interface{}{
			"message": "service.IndexSystem3",
		},
	)
	system, err = r.doCacheIndexSystem(cacheKey)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.IndexSystem3",
				"err":     err,
			},
		)
		return nil, err
	}
	return system, nil

}

func (r *SystemService) doCacheIndexSystem(cacheKey string) (system *models.ZBaseSys, err error) {
	system = new(models.ZBaseSys)
	err = r.Context.Db.Table((&models.ZBaseSys{}).TableName()).Find(system).Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.doCacheIndexSystem2",
				"err":      err,
				"cacheKey": cacheKey,
			},
		)
		return system, err
	}
	jsonRes, err := json.Marshal(&system)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.doCacheIndexSystem2",
				"err":      err,
				"cacheKey": cacheKey,
			},
		)
		return system, err
	}
	err = r.Context.CacheClient.Set(r.Context.GinContext.Request.Context(), cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":               "service.doCacheIndexSystem3",
				"err":                   err,
				"cacheKey":              cacheKey,
				"jsonRes":               jsonRes,
				"DataCacheTimeDuration": common.Conf.DataCacheTimeDuration,
			},
		)
		return system, err
	}
	return system, nil
}
