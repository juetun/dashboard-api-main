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
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"

	"time"
)

type LinkService struct {
	base.ServiceBase
}

func NewLinkService() *LinkService {
	return &LinkService{}
}
func (r *LinkService) LinkList(offset int, limit int) (links []models.ZLinks, cnt int64, err error) {
	links = make([]models.ZLinks, 0)
	cnt, err = r.Db.Asc("order").Limit(limit, offset).FindAndCount(&links)
	return
}

func (r *LinkService) LinkSore(ls pojos.LinkStore) (err error) {
	LinkInsert := models.ZLinks{
		Name:  ls.Name,
		Link:  ls.Link,
		Order: ls.Order,
	}
	_, err = r.Db.Insert(&LinkInsert)
	return
}

func (r *LinkService) LinkDetail(linkId int) (link *models.ZLinks, err error) {
	link = new(models.ZLinks)
	_, err = r.Db.Id(linkId).Get(link)
	return
}

func (r *LinkService) LinkUpdate(ls pojos.LinkStore, linkId int) (err error) {
	linkUpdate := models.ZLinks{
		Link:  ls.Link,
		Name:  ls.Name,
		Order: ls.Order,
	}
	_, err = r.Db.Id(linkId).Update(&linkUpdate)
	return
}

func (r *LinkService) LinkDestroy(linkId int) (err error) {
	link := new(models.ZLinks)
	_, err = r.Db.Id(linkId).Delete(link)
	return
}

func (r *LinkService) LinkCnt() (cnt int64, err error) {
	link := new(models.ZLinks)
	cnt, err = r.Db.Count(link)
	return
}

func (r *LinkService) AllLink() (links []models.ZLinks, err error) {

	cacheKey := common.Conf.LinkIndexKey
	cacheRes, err := r.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		links, err := r.doCacheLinkList(cacheKey)
		if err != nil {
			r.Log.Infoln("message", "service.AllLink", "err", err.Error())
			return links, err
		}
		return links, nil
	} else if err != nil {
		r.Log.Infoln("message", "service.AllLink", "err", err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(cacheRes), &links)
	if err != nil {
		r.Log.Errorln("message", "service.AllLink", "err", err.Error())
		links, err = r.doCacheLinkList(cacheKey)
		if err != nil {
			r.Log.Errorln("message", "service.AllLink", "err", err.Error())
			return nil, err
		}
		return links, nil
	}
	return links, nil
}

func (r *LinkService) doCacheLinkList(cacheKey string) (links []models.ZLinks, err error) {
	links = make([]models.ZLinks, 0)
	err = r.Db.Find(&links)
	if err != nil {
		r.Log.Errorln("message", "service.doCacheLinkList", "err", err.Error())
		return links, err
	}
	jsonRes, err := json.Marshal(&links)
	if err != nil {
		r.Log.Errorln("message", "service.doCacheLinkList", "err", err.Error())
		return nil, err
	}
	err = r.CacheClient.Set(cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Log.Errorln("message", "service.doCacheLinkList", "err", err.Error())
		return nil, err
	}
	return links, nil
}
