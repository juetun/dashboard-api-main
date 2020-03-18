/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-07
 * Time: 22:12
 */
package services

import (
	"encoding/json"

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

func NewSystemService() *SystemService {
	return &SystemService{}
}
func (r *SystemService) GetSystemList() (system *models.ZSystems, err error) {
	system = new(models.ZSystems)
	_, err = r.Db.Get(system)
	if err != nil {
		r.Log.Errorln("message", "service.GetSystemList", "err", err.Error())
		return
	}
	if system.Id <= 0 {
		systemInsert := models.ZSystems{
			Theme:        common.Conf.Theme,
			Title:        common.Conf.Title,
			Keywords:     common.Conf.Keywords,
			Description:  common.Conf.Description,
			RecordNumber: common.Conf.RecordNumber,
		}
		_, err = r.Db.Insert(systemInsert)
		if err != nil {
			r.Log.Errorln("message", "service.GetSystemList", "err", err.Error())
			return
		}
		_, err = r.Db.Get(system)
		if err != nil {
			r.Log.Errorln("message", "service.GetSystemList", "err", err.Error())
			return
		}
	}
	return
}

func (r *SystemService) SystemUpdate(sId int, ss pojos.ConsoleSystem) error {
	systemUpdate := models.ZSystems{
		Title:        ss.Title,
		Keywords:     ss.Keywords,
		Description:  ss.Description,
		RecordNumber: ss.RecordNumber,
		Theme:        ss.Theme,
	}
	_, err := r.Db.Id(sId).Update(&systemUpdate)
	return err
}

func (r *SystemService) IndexSystem() (system *models.ZSystems, err error) {
	cacheKey := common.Conf.SystemIndexKey
	cacheRes, err := r.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		system, err := r.doCacheIndexSystem(cacheKey)
		if err != nil {
			r.Log.Infoln("message", "service.IndexSystem", "err", err.Error())
			return system, err
		}
		return system, nil
	} else if err != nil {
		r.Log.Infoln("message", "service.IndexSystem", "err", err.Error())
		return system, err
	}

	err = json.Unmarshal([]byte(cacheRes), &system)
	if err != nil {
		r.Log.Errorln("message", "service.IndexSystem", "err", err.Error())
		system, err = r.doCacheIndexSystem(cacheKey)
		if err != nil {
			r.Log.Errorln("message", "service.IndexSystem", "err", err.Error())
			return nil, err
		}
		return system, nil
	}
	return system, nil
}

func (r *SystemService) doCacheIndexSystem(cacheKey string) (system *models.ZSystems, err error) {
	system = new(models.ZSystems)
	_, err = r.Db.Get(system)
	if err != nil {
		r.Log.Infoln("message", "service.doCacheIndexSystem", "err", err.Error())
		return system, err
	}
	jsonRes, err := json.Marshal(&system)
	if err != nil {
		r.Log.Errorln("message", "service.doCacheIndexSystem", "err", err.Error())
		return system, err
	}
	err = r.CacheClient.Set(cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Log.Errorln("message", "service.doCacheIndexSystem", "err", err.Error())
		return system, err
	}
	return system, nil
}
