/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-08
 * Time: 22:35
 */
package services

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"

	"time"
)

type LinkService struct {
	base.ServiceBase
}

func NewLinkService(context ...*base.Context) (p *LinkService) {
	p = &LinkService{}
	p.SetContext(context)
	return
}
func (r *LinkService) LinkList(offset int, limit int) (links []models.ZLinks, cnt int64, err error) {
	links = make([]models.ZLinks, 0)
	dba := r.Context.Db.Table((&models.ZLinks{}).TableName()).Unscoped().Where("deleted_at IS NULL")
	err = dba.Count(&cnt).Error
	if cnt > 0 {
		err = dba.Order("`order` asc").Limit(limit).
			Offset(offset).
			Find(&links).Error
	}
	return
}

func (r *LinkService) LinkSore(ls pojos.LinkStore) (err error) {
	LinkInsert := models.ZLinks{
		Name:  ls.Name,
		Link:  ls.Link,
		Order: ls.Order,
	}
	err = r.Context.Db.Table((&models.ZLinks{}).TableName()).Create(&LinkInsert).Error
	return
}

func (r *LinkService) LinkDetail(linkId int) (link *models.ZLinks, err error) {
	link = &models.ZLinks{}
	err = r.Context.Db.Table((&models.ZLinks{}).TableName()).
		Where("id=?", linkId).
		Find(link).Error
	return
}

func (r *LinkService) LinkUpdate(ls pojos.LinkStore, linkId int) (err error) {
	linkUpdate := models.ZLinks{
		Link:  ls.Link,
		Name:  ls.Name,
		Order: ls.Order,
	}
	err = r.Context.Db.Table((&models.ZLinks{}).TableName()).
		Where("id=?", linkId).
		Update(&linkUpdate).Error
	return
}

func (r *LinkService) LinkDestroy(linkId int) (err error) {
	link := new(models.ZLinks)
	err = r.Context.Db.Table(link.TableName()).
		Where("id=?", linkId).
		Delete(link).Error
	return
}

func (r *LinkService) LinkCnt() (cnt int64, err error) {
	link := new(models.ZLinks)
	err = r.Context.Db.Table((&models.ZLinks{}).TableName()).
		Count(link).
		Error
	return
}

func (r *LinkService) AllLink() (links []models.ZLinks, err error) {

	cacheKey := common.Conf.LinkIndexKey
	cacheRes, err := r.Context.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		links, err := r.doCacheLinkList(cacheKey)
		if err != nil {
			r.Context.Log.Infoln("message", "service.AllLink", "err", err.Error())
			return links, err
		}
		return links, nil
	} else if err != nil {
		r.Context.Log.Infoln("message", "service.AllLink", "err", err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(cacheRes), &links)
	if err != nil {
		r.Context.Log.Errorln("message", "service.AllLink", "err", err.Error())
		links, err = r.doCacheLinkList(cacheKey)
		if err != nil {
			r.Context.Log.Errorln("message", "service.AllLink", "err", err.Error())
			return nil, err
		}
		return links, nil
	}
	return links, nil
}

func (r *LinkService) doCacheLinkList(cacheKey string) (links []models.ZLinks, err error) {
	links = make([]models.ZLinks, 0)
	err = r.Context.Db.Table((&models.ZLinks{}).TableName()).
		Find(&links).
		Error
	if err != nil {
		r.Context.Log.Errorln("message", "service.doCacheLinkList", "err", err.Error())
		return links, err
	}
	jsonRes, err := json.Marshal(&links)
	if err != nil {
		r.Context.Log.Errorln("message", "service.doCacheLinkList", "err", err.Error())
		return nil, err
	}
	err = r.Context.CacheClient.Set(cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Context.Log.Errorln("message", "service.doCacheLinkList", "err", err.Error())
		return nil, err
	}
	return links, nil
}
