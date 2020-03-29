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
	"github.com/jinzhu/gorm"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
)

type TagService struct {
	base.ServiceBase
}


func NewTagService(context ...*base.Context) (p *TagService) {
	p = &TagService{}
	p.SetContext(context)
	return
}
func (r *TagService) TagStore(ts pojos.TagStore) (err error) {
	tag := new(models.ZTags)
	dba := r.Context.Db.
		Table((&models.ZTags{}).TableName())
	err = dba.
		Where("name = ?", ts.Name).
		Find(tag).
		Error
	if err != nil {
		r.Context.Log.Errorln("message", "service.TagStore", "error", err.Error())
		return err
	}

	if tag.Id > 0 {
		r.Context.Log.Errorln("message", "service.TagStore", "error", "Tag has exists")
		return errors.New("Tag has exists")
	}

	tagInsert := &models.ZTags{
		Name:        ts.Name,
		DisplayName: ts.DisplayName,
		SeoDesc:     ts.SeoDesc,
		Num:         0,
	}
	err = dba.Create(tagInsert).Error
	r.Context.CacheClient.Del(common.Conf.TagListKey)
	return
}

func (r *TagService) GetPostTagsByPostId(postId int) (tagsArr []int, err error) {
	rows, err := r.Context.Db.Table((&models.ZPostTag{}).TableName()).Where("post_id = ?", postId).Select("tag_id").Rows()
	if err != nil {
		r.Context.Log.Errorln("message", "service.GetPostTagsByPostId", "error", err.Error())
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
	err = r.getTableDb().Where("id=?", tagId).Find(tag).Error
	return
}
func (r *TagService) getTableDb() *gorm.DB {
	return r.Context.Db.Table((&models.ZTags{}).TableName())
}

func (r *TagService) TagUpdate(tagId int, ts pojos.TagStore) (err error) {
	tagUpdate := &models.ZTags{
		Name:        ts.Name,
		DisplayName: ts.DisplayName,
		SeoDesc:     ts.SeoDesc,
	}
	err = r.getTableDb().Where("id=?", tagId).
		Update(tagUpdate).Error
	return err
}

func (r *TagService) GetTagsByIds(tagIds []int) (tags []*models.ZTags, err error) {
	tags = make([]*models.ZTags, 0)
	err = r.getTableDb().
		Where("id in (?)", tagIds).
		Select("id,name,display_name,seo_desc,num").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagService) TagsIndex(limit int, offset int) (num int64, tags []*models.ZTags, err error) {
	tags = make([]*models.ZTags, 0)
	dba := r.getTableDb()
	err = dba.Count(&num).Error
	if err != nil {
		return
	}
	if num > 0 {
		err = dba.Order("num desc ").
			Limit(limit).
			Offset(offset).
			Find(&tags).
			Error
	}
	return
}

func (r *TagService) DelTagRel(tagId int) {
	session := r.Context.Db.Begin().Table((&models.ZPostTag{}).TableName())
	defer session.Close()
	postTag := new(models.ZPostTag)
	err := session.Where("tag_id = ?", tagId).Delete(postTag).Error
	if err != nil {
		_ = session.Rollback()
		r.Context.Log.Errorln("message", "service.DelTagRel", "err", err.Error())
		return
	}
	tag := new(models.ZTags)
	err = session.Where("id=?", tagId).Delete(tag).Error
	if err != nil {
		_ = session.Rollback()
		r.Context.Log.Errorln("message", "service.DelTagRel", "err", err.Error())
		return
	}
	session.Commit()
	r.Context.CacheClient.Del(common.Conf.TagListKey)
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
		r.Context.Log.Errorln("message", "service.Index.CommonData", "err", err.Error())
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
		r.Context.Log.Errorln("message", "service.Index.CommonData", "err", err.Error())
		return
	}

	srvLink := NewLinkService()
	links, err := srvLink.AllLink()
	if err != nil {
		r.Context.Log.Errorln("message", "service.Index.CommonData", "err", err.Error())
		return
	}
	srvSystem := NewSystemService()

	system, err := srvSystem.IndexSystem()
	if err != nil {
		r.Context.Log.Errorln("message", "service.Index.CommonData", "err", err.Error())
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
	cacheRes, err := r.Context.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		tags, err := r.doCacheTagList(cacheKey)
		if err != nil {
			r.Context.Log.Infoln("message", "service.AllTags", "err", err.Error())
			return tags, err
		}
		return tags, nil
	}
	if err != nil {
		r.Context.Log.Infoln("message", "service.AllTags", "err", err.Error())
		return nil, err
	}

	var cacheTag []models.ZTags
	err = json.Unmarshal([]byte(cacheRes), &cacheTag)
	if err != nil {
		r.Context.Log.Errorln("message", "service.AllTags", "err", err.Error())
		tags, err := r.doCacheTagList(cacheKey)
		if err != nil {
			r.Context.Log.Errorln("message", "service.AllTags", "err", err.Error())
			return nil, err
		}
		return tags, nil
	}
	return cacheTag, nil
}

func (r *TagService) doCacheTagList(cacheKey string) ([]models.ZTags, error) {
	tags, err := r.tags()
	if err != nil {
		r.Context.Log.Infoln("message", "service.doCacheTagList", "err", err.Error())
		return tags, err
	}
	jsonRes, err := json.Marshal(&tags)
	if err != nil {
		r.Context.Log.Errorln("message", "service.doCacheTagList", "err", err.Error())
		return nil, err
	}
	err = r.Context.CacheClient.Set(cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Context.Log.Errorln("message", "service.doCacheTagList", "err", err.Error())
		return nil, err
	}
	return tags, nil
}

func (r *TagService) tags() ([]models.ZTags, error) {
	tags := make([]models.ZTags, 0)
	err := r.getTableDb().Find(&tags).Error
	if err != nil {
		r.Context.Log.Infoln("message", "service.Tags", "err", err.Error())
		return tags, err
	}

	return tags, nil
}

func (r *TagService) TagCnt() (cnt int64, err error) {
	err = r.getTableDb().Count(&cnt).Error
	return
}
