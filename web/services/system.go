/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-07
 * Time: 22:12
 */
package services

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/pojos"

	"github.com/go-redis/redis"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/web/models"

	"time"
)

type SystemService struct {
	base.ServiceBase
}

func NewSystemService(context ...*base.Context) (p *SystemService) {
	p = &SystemService{}
	p.SetContext(context)
	return
}
func (r *SystemService) GetSystemList() (system *models.ZBaseSys, err error) {
	system = new(models.ZBaseSys)
	err = r.Context.Db.Table(system.TableName()).Find(system).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) { // 如果没有查询到数据
			err = nil
			return
		}
		r.Context.Log.Errorln("message", "service.GetSystemList", "err", err.Error())
		return
	}
	if system.Id <= 0 {
		return
	}
	systemInsert := models.ZBaseSys{
		Theme:        1,
		Title:        "汽车家",
		Keywords:     "汽车信息",
		Description:  "阿斯顿发生",
		RecordNumber: "京ICP-0001",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		// Theme:       "jekyll-theme-cayman",
		// Title:        common.Conf.Title,
		// Keywords:     common.Conf.Keywords,
		// Description:  common.Conf.Description,
		// RecordNumber: common.Conf.RecordNumber,
	}
	dba := r.Context.Db.Table((&models.ZBaseSys{}).TableName())
	err = dba.Create(&systemInsert).Error
	if err != nil {
		r.Context.Log.Errorln("message", "service.GetSystemList", "err", err.Error())
		return
	}
	err = dba.Find(system).Error
	if err != nil {
		r.Context.Log.Errorln("message", "service.GetSystemList", "err", err.Error())
		return
	}

	return
}

func (r *SystemService) SystemUpdate(sId int, ss pojos.ConsoleSystem) error {
	systemUpdate := models.ZBaseSys{
		Title:        ss.Title,
		Keywords:     ss.Keywords,
		Description:  ss.Description,
		RecordNumber: ss.RecordNumber,
		Theme:        ss.Theme,
	}
	err := r.Context.Db.Where("id=?", sId).
		Update(&systemUpdate).
		Error
	return err
}

func (r *SystemService) IndexSystem() (system *models.ZBaseSys, err error) {
	cacheKey := common.Conf.SystemIndexKey
	cacheRes, err := r.Context.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		system, err := r.doCacheIndexSystem(cacheKey)
		if err != nil {
			r.Context.Log.Infoln("message", "service.IndexSystem", "err", err.Error())
			return system, err
		}
		return system, nil
	}
	if err != nil {
		r.Context.Log.Infoln("message", "service.IndexSystem", "err", err.Error())
		return system, err
	}

	err = json.Unmarshal([]byte(cacheRes), &system)
	if err == nil {
		return system, nil
	}
	r.Context.Log.Errorln("message", "service.IndexSystem", "err", err.Error())
	system, err = r.doCacheIndexSystem(cacheKey)
	if err != nil {
		r.Context.Log.Errorln("message", "service.IndexSystem", "err", err.Error())
		return nil, err
	}
	return system, nil

}

func (r *SystemService) doCacheIndexSystem(cacheKey string) (system *models.ZBaseSys, err error) {
	system = new(models.ZBaseSys)
	err = r.Context.Db.Table((&models.ZBaseSys{}).TableName()).Find(system).Error
	if err != nil {
		r.Context.Log.Infoln("message", "service.doCacheIndexSystem", "err", err.Error())
		return system, err
	}
	jsonRes, err := json.Marshal(&system)
	if err != nil {
		r.Context.Log.Errorln("message", "service.doCacheIndexSystem", "err", err.Error())
		return system, err
	}
	err = r.Context.CacheClient.Set(cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Context.Log.Errorln("message", "service.doCacheIndexSystem", "err", err.Error())
		return system, err
	}
	return system, nil
}
