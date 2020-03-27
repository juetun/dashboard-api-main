/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-12
 * Time: 01:57
 */
package services

import (
	"encoding/json"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
)

type TagService struct {
	base.ServiceBase
}

func NewTagService() *TagService {
	return &TagService{}
}
func (r *TagService) TagStore(ts pojos.TagStore) (err error) {
	tag := new(models.ZTags)
	_, err = r.Db.Where("name = ?", ts.Name).Get(tag)
	if err != nil {
		r.Log.Errorln("message", "service.TagStore", "error", err.Error())
		return err
	}

	if tag.Id > 0 {
		r.Log.Errorln("message", "service.TagStore", "error", "Tag has exists")
		return errors.New("Tag has exists")
	}

	tagInsert := &models.ZTags{
		Name:        ts.Name,
		DisplayName: ts.DisplayName,
		SeoDesc:     ts.SeoDesc,
		Num:         0,
	}
	_, err = r.Db.Insert(tagInsert)
	r.CacheClient.Del(common.Conf.TagListKey)
	return
}

func (r *TagService) GetPostTagsByPostId(postId int) (tagsArr []int, err error) {
	postTag := new(models.ZPostTag)
	rows, err := r.Db.Where("post_id = ?", postId).Cols("tag_id").Rows(postTag)
	if err != nil {
		r.Log.Errorln("message", "service.GetPostTagsByPostId", "error", err.Error())
		return nil, nil
	}
	defer rows.Close()
	for rows.Next() {
		postTag := new(models.ZPostTag)
		err = rows.Scan(postTag)
		if err != nil {
			return nil, err
		}
		tagsArr = append(tagsArr, postTag.TagId)
	}
	return
}

func (r *TagService) GetTagById(tagId int) (tag *models.ZTags, err error) {
	tag = new(models.ZTags)
	_, err = r.Db.ID(tagId).Get(tag)
	return
}

func (r *TagService) TagUpdate(tagId int, ts pojos.TagStore) error {
	tagUpdate := &models.ZTags{
		Name:        ts.Name,
		DisplayName: ts.DisplayName,
		SeoDesc:     ts.SeoDesc,
	}
	_, err := r.Db.ID(tagId).Update(tagUpdate)
	return err
}

func (r *TagService) GetTagsByIds(tagIds []int) ([]*models.ZTags, error) {
	tags := make([]*models.ZTags, 0)
	err := r.Db.In("id", tagIds).Cols("id", "name", "display_name", "seo_desc", "num").Find(&tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagService) TagsIndex(limit int, offset int) (num int64, tags []*models.ZTags, err error) {
	tags = make([]*models.ZTags, 0)
	num, err = r.Db.Desc("num").Limit(limit, offset).FindAndCount(&tags)
	return
}

func (r *TagService) DelTagRel(tagId int) {
	session := r.Db.NewSession()
	defer session.Close()
	postTag := new(models.ZPostTag)
	_, err := session.Where("tag_id = ?", tagId).Delete(postTag)
	if err != nil {
		_ = session.Rollback()
		r.Log.Errorln("message", "service.DelTagRel", "err", err.Error())
		return
	}
	tag := new(models.ZTags)
	_, err = session.ID(tagId).Delete(tag)
	if err != nil {
		_ = session.Rollback()
		r.Log.Errorln("message", "service.DelTagRel", "err", err.Error())
		return
	}
	err = session.Commit()
	if err != nil {
		_ = session.Rollback()
		r.Log.Errorln("message", "service.DelTagRel", "err", err.Error())
		return
	}
	r.CacheClient.Del(common.Conf.TagListKey)
	return
}
func (r *TagService) CommonData() (h gin.H, err error) {
	h = gin.H{
		"themeJs":          common.Conf.ThemeJs,
		"themeCss":         common.Conf.ThemeCss,
		"themeImg":         common.Conf.ThemeImg,
		"themeFancyboxCss": common.Conf.ThemeFancyboxCss,
		"themeFancyboxJs":  common.Conf.ThemeFancyboxJs,
		"themeHLightCss":   common.Conf.ThemeHLightCss,
		"themeHLightJs":    common.Conf.ThemeHLightJs,
		"themeShareCss":    common.Conf.ThemeShareCss,
		"themeShareJs":     common.Conf.ThemeShareJs,
		"themeArchivesJs":  common.Conf.ThemeArchivesJs,
		"themeArchivesCss": common.Conf.ThemeArchivesCss,
		"themeNiceImg":     common.Conf.ThemeNiceImg,
		"themeAllCss":      common.Conf.ThemeAllCss,
		"themeIndexImg":    common.Conf.ThemeIndexImg,
		"themeCateImg":     common.Conf.ThemeCateImg,
		"themeTagImg":      common.Conf.ThemeTagImg,
		"title":            "",
		"tem":              "defaultList",
	}
	h["script"] = template.HTML(common.Conf.OtherScript)
	srv := NewCategoryService()
	cates, err := srv.CateListBySort()
	if err != nil {
		r.Log.Errorln("message", "service.Index.CommonData", "err", err.Error())
		return
	}
	var catess []pojos.IndexCategory
	for _, v := range cates {
		c := pojos.IndexCategory{
			Cates: v.Cates,
			Html:  template.HTML(v.Html),
		}
		catess = append(catess, c)
	}

	tags, err := r.AllTags()
	if err != nil {
		r.Log.Errorln("message", "service.Index.CommonData", "err", err.Error())
		return
	}

	srvLink := NewLinkService()
	links, err := srvLink.AllLink()
	if err != nil {
		r.Log.Errorln("message", "service.Index.CommonData", "err", err.Error())
		return
	}
	srvSystem := NewSystemService()

	system, err := srvSystem.IndexSystem()
	if err != nil {
		r.Log.Errorln("message", "service.Index.CommonData", "err", err.Error())
		return
	}
	h["cates"] = catess
	h["system"] = system
	h["links"] = links
	h["tags"] = tags
	return
}

func (r *TagService) AllTags() ([]models.ZTags, error) {
	cacheKey := common.Conf.TagListKey
	cacheRes, err := r.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		tags, err := r.doCacheTagList(cacheKey)
		if err != nil {
			r.Log.Infoln("message", "service.AllTags", "err", err.Error())
			return tags, err
		}
		return tags, nil
	}
	if err != nil {
		r.Log.Infoln("message", "service.AllTags", "err", err.Error())
		return nil, err
	}

	var cacheTag []models.ZTags
	err = json.Unmarshal([]byte(cacheRes), &cacheTag)
	if err != nil {
		r.Log.Errorln("message", "service.AllTags", "err", err.Error())
		tags, err := r.doCacheTagList(cacheKey)
		if err != nil {
			r.Log.Errorln("message", "service.AllTags", "err", err.Error())
			return nil, err
		}
		return tags, nil
	}
	return cacheTag, nil
}

func (r *TagService) doCacheTagList(cacheKey string) ([]models.ZTags, error) {
	tags, err := r.tags()
	if err != nil {
		r.Log.Infoln("message", "service.doCacheTagList", "err", err.Error())
		return tags, err
	}
	jsonRes, err := json.Marshal(&tags)
	if err != nil {
		r.Log.Errorln("message", "service.doCacheTagList", "err", err.Error())
		return nil, err
	}
	err = r.CacheClient.Set(cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Log.Errorln("message", "service.doCacheTagList", "err", err.Error())
		return nil, err
	}
	return tags, nil
}

func (r *TagService) tags() ([]models.ZTags, error) {
	tags := make([]models.ZTags, 0)
	err := r.Db.Find(&tags)
	if err != nil {
		r.Log.Infoln("message", "service.Tags", "err", err.Error())
		return tags, err
	}

	return tags, nil
}

func (r *TagService) TagCnt() (cnt int64, err error) {
	tag := new(models.ZTags)
	cnt, err = r.Db.Count(tag)
	return
}
