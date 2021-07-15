// Package srv_impl
/**
 * Created by GoLand.
 * Date: 2019-01-12
 * Time: 01:57
 */
package srv_impl

import (
	"encoding/json"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type TagService struct {
	base.ServiceBase
}

func NewTagService(context ...*base.Context) (p *TagService) {
	p = &TagService{}
	p.SetContext(context...)
	return
}
func (r *TagService) TagStore(ts wrappers.TagStore) (err error) {
	tag := new(models.ZTags)
	dba := r.Context.Db.
		Table((&models.ZTags{}).TableName())
	err = dba.
		Where("name = ?", ts.Name).
		Find(tag).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.TagStore", "error": err,
			},
		)
		return
	}

	if tag.Id > 0 {
		err = errors.New("您输入的标签名已存在")
		r.Context.Error(
			map[string]interface{}{
				"message": "service.TagStore1", "error": err,
			},
		)
		return
	}

	tagInsert := &models.ZTags{
		Name:        ts.Name,
		DisplayName: ts.DisplayName,
		SeoDesc:     ts.SeoDesc,
		Num:         0,
	}
	err = dba.Create(tagInsert).Error
	r.Context.CacheClient.Del(r.Context.GinContext.Request.Context(), common.Conf.TagListKey)
	return
}
func (r *TagService) GetPostTagsByPostIds(postIds []string) (res *map[int][]wrappers.ConsoleTag, err error) {
	res = &map[int][]wrappers.ConsoleTag{}
	if len(postIds) == 0 {
		return
	}
	var dt []models.ZPostTag
	err = r.Context.Db.Table((&models.ZPostTag{}).TableName()).Where("post_id in (?)", postIds).
		Find(&dt).Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.GetPostTagsByPostId",
				"error":   err,
				"postIds": postIds,
			},
		)
		return
	}
	tagIds := r.uniqueTagId(&dt)
	var mp *map[int]models.ZTags
	if mp, err = r.GetTagsMapByIds(tagIds); err != nil {
		return
	}
	for _, value := range dt {
		if _, ok := (*res)[value.PostId]; !ok {
			(*res)[value.PostId] = make([]wrappers.ConsoleTag, 0)
		}
		p := wrappers.PostTagShow{
			ZPostTag: value,
			ZTags:    models.ZTags{},
		}
		if _, ok := (*mp)[value.TagId]; ok {
			p.ZTags = (*mp)[value.TagId]
		}
		(*res)[value.PostId] = append((*res)[value.PostId], wrappers.ConsoleTag{
			Id:          p.ZTags.Id,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			SeoDesc:     p.SeoDesc,
			Num:         p.Num,
		})
	}
	return

}
func (r *TagService) GetTagsMapByIds(ids *[]int) (res *map[int]models.ZTags, err error) {
	res = &map[int]models.ZTags{}
	if len(*ids) == 0 {
		return
	}
	var dt []models.ZTags
	err = r.Context.Db.Table((&models.ZTags{}).TableName()).
		Where("id in (?)", *ids).
		Find(&dt).
		Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.GetTagsMapByIds",
				"error":   err,
				"ids":     ids,
			},
		)
		return
	}

	for _, value := range dt {
		(*res)[value.Id] = value
	}
	return
}

func (r *TagService) uniqueTagId(dt *[]models.ZPostTag) *[]int {
	ids := make([]int, 0)
	mapIds := make(map[int]int)
	for _, value := range *dt {
		if _, ok := mapIds[value.TagId]; !ok {
			ids = append(ids, value.TagId)
			mapIds[value.TagId] = value.TagId
		}
	}
	return &ids
}

func (r *TagService) GetPostTagsByPostId(postId int) (tagsArr []int, err error) {
	rows, err := r.Context.Db.Table((&models.ZPostTag{}).TableName()).
		Where("post_id = ?", postId).
		Select("tag_id").
		Rows()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.GetPostTagsByPostId",
				"error":   err,
				"postId":  postId,
			},
		)
		return
	}
	defer func() {
		err = rows.Close()
	}()
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

func (r *TagService) TagUpdate(tagId int, ts wrappers.TagStore) (err error) {
	var tag = &models.ZTags{}
	dba := r.getTableDb()
	err = dba.
		Where("name = ? AND id!=?", ts.Name, tagId).
		Find(tag).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.GetPostTagsByPostId",
				"error":   err,
				"tagId":   tagId,
				"ts":      ts,
			},
		)
		return err
	}
	if tag.Id > 0 {

		err = errors.New("您输入的标签名已存在")
		r.Context.Error(
			map[string]interface{}{
				"message": "service.GetPostTagsByPostId1",
				"error":   err,
				"tagId":   tagId,
				"ts":      ts,
			},
		)
		return
	}

	tagUpdate := &models.ZTags{
		Name:        ts.Name,
		DisplayName: ts.DisplayName,
		SeoDesc:     ts.SeoDesc,
	}
	err = r.getTableDb().Where("id=?", tagId).
		Updates(tagUpdate).Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":   "service.GetPostTagsByPostId2",
				"error":     err,
				"tagUpdate": tagUpdate,
			},
		)
	}
	return err
}

func (r *TagService) GetTagsByIds(tagIds []int) (tags []*models.ZTags, err error) {
	tags = make([]*models.ZTags, 0)
	err = r.getTableDb().
		Where("id in (?)", tagIds).
		Select("id,name,display_name,seo_desc,num").Find(&tags).Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.GetTagsByIds",
				"error":   err,
				"tagIds":  tagIds,
			},
		)
		return
	}
	return
}

func (r *TagService) TagsIndex(limit int, offset int) (num int64, tags []*models.ZTags, err error) {
	tags = make([]*models.ZTags, 0)
	dba := r.getTableDb().Table((&models.ZTags{}).TableName()).Unscoped().Where("deleted_at IS NULL")
	err = dba.Count(&num).Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.TagsIndex",
				"error":   err,
				"limit":   limit,
				"offset":  offset,
			},
		)
		return
	}
	if num > 0 {
		err = dba.Order("num desc ").
			Limit(limit).
			Offset(offset).
			Find(&tags).
			Error
		if err != nil {
			r.Context.Error(
				map[string]interface{}{
					"message": "service.TagsIndex",
					"error":   err,
					"limit":   limit,
					"offset":  offset,
				},
			)
		}
	}
	return
}

func (r *TagService) DelTagRel(tagId int) {
	session := r.Context.Db.Begin().Table((&models.ZPostTag{}).TableName())
	defer session.Commit()
	postTag := new(models.ZPostTag)
	err := session.Where("tag_id = ?", tagId).Delete(postTag).Error
	if err != nil {
		_ = session.Rollback()
		r.Context.Error(
			map[string]interface{}{
				"message": "service.DelTagRel",
				"error":   err,
				"tagId":   tagId,
			},
		)
		return
	}
	tag := new(models.ZTags)
	err = session.Where("id=?", tagId).Delete(tag).Error
	if err != nil {
		_ = session.Rollback()
		r.Context.Error(
			map[string]interface{}{
				"message": "service.DelTagRel1",
				"error":   err,
				"tagId":   tagId,
			},
		)
		return
	}
	r.Context.CacheClient.Del(r.Context.GinContext.Request.Context(), common.Conf.TagListKey)
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
		r.Context.Error(
			map[string]interface{}{
				"message": "service.Index.CommonData",
				"error":   err,
			},
		)
		return
	}
	var catess []wrappers.IndexCategory
	for _, v := range cates {
		c := wrappers.IndexCategory{
			Cates: v.Cates,
			Html:  template.HTML(v.Html),
		}
		catess = append(catess, c)
	}

	tags, err := r.AllTags()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.Index.CommonData1",
				"error":   err,
			},
		)
		return
	}

	srvLink := NewLinkService()
	links, err := srvLink.AllLink()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.Index.CommonData2",
				"error":   err,
			},
		)
		return
	}
	srvSystem := NewSystemService()

	system, err := srvSystem.IndexSystem()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.Index.CommonData3",
				"error":   err,
			},
		)
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
	cacheRes, err := r.Context.CacheClient.Get(r.Context.GinContext.Request.Context(), cacheKey).Result()
	if err == redis.Nil {
		tags, err := r.doCacheTagList(cacheKey)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{
					"message": "service.AllTags",
					"error":   err,
				},
			)
			return tags, err
		}
		return tags, nil
	}
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.AllTags1",
				"error":   err,
			},
		)
		return nil, err
	}

	var cacheTag []models.ZTags
	err = json.Unmarshal([]byte(cacheRes), &cacheTag)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.AllTags2",
				"error":   err,
			},
		)
		tags, err := r.doCacheTagList(cacheKey)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{
					"message": "service.AllTags3",
					"error":   err,
				},
			)
			return nil, err
		}
		return tags, nil
	}
	return cacheTag, nil
}

func (r *TagService) doCacheTagList(cacheKey string) ([]models.ZTags, error) {
	tags, err := r.tags()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.doCacheTagList",
				"error":    err,
				"cacheKey": cacheKey,
			},
		)
		return tags, err
	}
	jsonRes, err := json.Marshal(&tags)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.doCacheTagList1",
				"error":    err,
				"cacheKey": cacheKey,
			},
		)
		return nil, err
	}
	err = r.Context.CacheClient.Set(r.Context.GinContext.Request.Context(), cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.doCacheTagList2",
				"error":    err,
				"cacheKey": cacheKey,
			},
		)
		return nil, err
	}
	return tags, nil
}

func (r *TagService) tags() ([]models.ZTags, error) {
	tags := make([]models.ZTags, 0)
	err := r.getTableDb().Find(&tags).Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.Tags",
				"error":   err,
			},
		)
		return tags, err
	}

	return tags, nil
}

func (r *TagService) TagCnt() (cnt int64, err error) {
	err = r.getTableDb().Count(&cnt).Error
	return
}
