/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-11
 * Time: 23:24
 */
package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
)

type CategoryService struct {
	base.ServiceBase
}

func NewCategoryService(context ...*base.Context) (p *CategoryService) {
	p = &CategoryService{}
	p.SetContext(context)
	return
}
func (r *CategoryService) GetCateById(cateId int) (cate *models.ZCategories, err error) {
	cate = &models.ZCategories{}
	err = r.Context.Db.Table((&models.ZCategories{}).TableName()).
		Where("id=?", cateId).
		Find(cate).Error
	return
}

func (r *CategoryService) GetCateByParentId(parentId int) (cate *models.ZCategories, err error) {
	cate = &models.ZCategories{}
	err = r.Context.Db.Table((&models.ZCategories{}).TableName()).
		Where("parent_id = ?", parentId).
		Find(cate).Error
	return
}

func (r *CategoryService) DelCateRel(cateId int) {
	session := r.Context.Db.Begin()
	defer session.Close()
	postCate := new(models.ZPostCate)
	err := session.Where("cate_id = ?", cateId).Delete(postCate).Error
	if err != nil {
		_ = session.Rollback()
		r.Error("message", "service.DelCateRel", "err", err.Error())
		return
	}
	cate := new(models.ZCategories)
	err = session.Where("id=?", cateId).Delete(cate).Error
	if err != nil {
		_ = session.Rollback()
		r.Error("message", "service.DelCateRel", "err", err.Error())
		return
	}
	session.Commit()
	r.Context.CacheClient.Del(common.Conf.CateListKey)
	return
}

func (r *CategoryService) CateStore(cs pojos.CateStore) (res bool, err error) {

	defaultCate := new(models.ZCategories)
	err = r.Context.Db.Where("name = ?", cs.Name).Find(defaultCate).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		r.Error("message", "service.CateStore", "err", err.Error())
		return
	}
	if defaultCate.Id > 0 {
		r.Error("message", "service.CateStore", "err", "Cate has exists ")
		return res, errors.New("你输入的分类名称已存在")
	}

	if cs.ParentId > 0 {
		var cate models.ZCategories
		err = r.Context.Db.Where("id=?", cs.ParentId).Find(&cate).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				err = errors.New("你输入的分类上级分类不存在或已删除")
				return
			}
			r.Error("message", "service.CateStore", "err", err.Error())
			return
		}

	}

	cate := models.ZCategories{
		Name:        cs.Name,
		DisplayName: cs.DisplayName,
		SeoDesc:     cs.SeoDesc,
		ParentId:    cs.ParentId,
		Model: base.Model{
			CreatedAt: base.TimeNormal{Time: time.Now()},
			UpdatedAt: base.TimeNormal{Time: time.Now()},
		},
	}
	err = r.Context.Db.Create(&cate).Error
	if err != nil {
		r.Error("message", "service.CateStore", "err", err.Error())
		return
	}
	r.Context.CacheClient.Del(common.Conf.CateListKey)
	res = true
	return
}

func (r *CategoryService) CateUpdate(cateId int, cs pojos.CateStore) (res bool, err error) {
	cate := new(models.ZCategories)
	if cs.ParentId != 0 {
		err = r.Context.Db.Where("id=?", cs.ParentId).
			Find(cate).
			Error
		if err != nil {
			r.Error("message", "service.CateUpdate", "err", err.Error())
			return false, err
		}
		if !res || cate.Id < 1 {
			r.Error("message", "service.CateUpdate", err, "the parent id is not exists ")
			return false, errors.New("上级分类不存在")
		}
		ids := []int{cateId}
		resIds := []int{0}
		_, res2, _ := r.GetSimilar(ids, resIds, 0)
		for _, v := range res2 {
			if v == cs.ParentId {
				return false, errors.New("Can not be you child node ")
			}
		}
	}
	cateUpdate := &models.ZCategories{
		Name:        cs.Name,
		DisplayName: cs.DisplayName,
		SeoDesc:     cs.SeoDesc,
		ParentId:    cs.ParentId,
	}
	err = r.Context.Db.Table((&models.ZCategories{}).TableName()).
		Where("id =?", cateId).Update(cateUpdate).Error
	if err != nil {
		r.Error("message", "service.CateUpdate", "err", err.Error())
		return false, err
	}
	r.Context.CacheClient.Del(common.Conf.CateListKey)
	return true, nil
}

func (r *CategoryService) GetSimilar(beginId []int, resIds []int, level int) (beginId2 []int, resIds2 []int, level2 int) {
	if len(beginId) != 0 {
		cates := make([]*models.ZCategories, 0)
		err := r.Context.Db.Table((&models.ZCategories{}).TableName()).Where("parent_id in(?)", beginId).
			Find(&cates).Error
		if err != nil {
			r.Error("message", "service.GetSimilar", err, "the parent id data is not exists ")
			return []int{}, []int{}, 0
		}
		if len(cates) == 0 {
			return beginId, resIds, level
		}
		if level == 0 {
			resIds2 = beginId
		} else {
			resIds2 = resIds
		}
		for _, v := range cates {
			id := v.Id
			beginId2 = append(beginId2, id)
			resIds2 = append(resIds2, id)
		}
		level2 = level + 1
		return r.GetSimilar(beginId2, resIds2, level2)
	}
	return beginId, beginId, level

}

// 根据文章ID获取文章的分类
func (r *CategoryService) GetPostCateByPostIds(postIds *[]int) (res *map[int]pojos.PostShow, err error) {
	res = &map[int]pojos.PostShow{}
	if len(*postIds) == 0 {
		return
	}
	var dt []models.ZPostCate
	err = r.Context.Db.Table((&models.ZPostCate{}).TableName()).
		Where("post_id in (?)", *postIds).
		Find(&dt).Error
	if err != nil {
		return
	}
	cateIds := r.uniqueCateId(&dt)
	mp, err := r.GetCategoryByIds(cateIds)
	for _, value := range dt {
		p := pojos.PostShow{
			ZPostCate: value,
		}
		if _, ok := (*mp)[value.CateId]; ok {
			p.ZCategories = (*mp)[value.CateId]
		}
		(*res)[value.PostId] = p
	}
	return
}

func (r *CategoryService) GetCategoryByIds(ids *[]int) (res *map[int]models.ZCategories, err error) {
	res = &map[int]models.ZCategories{}
	if len(*ids) == 0 {
		return
	}
	var dt []models.ZCategories
	err = r.Context.Db.Table((&models.ZCategories{}).TableName()).Where("id in (?)", *ids).Find(&dt).Error
	if err != nil {
		return
	}

	for _, value := range dt {
		(*res)[value.Id] = value
	}
	return
}

func (r *CategoryService) uniqueCateId(dt *[]models.ZPostCate) *[]int {
	cateIds := make([]int, 0)
	mapCateIds := make(map[int]int)
	for _, value := range *dt {
		if _, ok := mapCateIds[value.CateId]; !ok {
			cateIds = append(cateIds, value.CateId)
			mapCateIds[value.CateId] = value.CateId
		}
	}
	return &cateIds
}

func (r *CategoryService) GetPostCateByPostId(postId int) (cates *models.ZCategories, err error) {
	postCate := new(models.ZPostCate)
	err = r.Context.Db.Table((&models.ZPostCate{}).TableName()).Select("cate_id").Where("post_id = ?", postId).
		Find(postCate).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			r.Error("message", "service.GetPostCateByPostId", "err", "there has not data")
			return cates, errors.New("can not get the post cate")
		}
		r.Error("message", "service.GetPostCateByPostId", "err", err.Error())
		return cates, err
	}

	cates = new(models.ZCategories)
	err = r.Context.Db.Table((&models.ZCategories{}).TableName()).
		Where("id = ?", postCate.CateId).
		Select("id,name,display_name,seo_desc").
		Find(cates).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			r.Error("message", "service.GetPostCateByPostId", "err", "there has not data")
			return cates, errors.New("can not get the post cate")
		}
		r.Error("message", "service.GetPostCateByPostId", "err", err.Error())
		return cates, err
	}
	return cates, nil

}

func (r *CategoryService) PostCate(postId int) (res int, err error) {
	postCate := new(models.ZPostCate)
	err = r.Context.Db.Table(postCate.TableName()).
		Where("post_id = ?", postId).
		Find(postCate).
		Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return
	}
	return postCate.CateId, err
}
func (r *CategoryService) GetPostCates(postId *[]int) (res *map[int]models.ZPostCate, err error) {
	res = &map[int]models.ZPostCate{}
	if len(*postId) == 0 {
		return
	}
	var postCate []models.ZPostCate
	err = r.Context.Db.Table((&models.ZPostCate{}).TableName()).
		Where("post_id in (?)", *postId).
		Find(&postCate).
		Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return
	}
	for _, value := range postCate {
		(*res)[value.PostId] = value
	}
	return res, err
}

// Get the cate list what by parent sort
func (r *CategoryService) CateListBySort() (res []pojos.Category, err error) {
	res = make([]pojos.Category, 0)
	cacheKey := common.Conf.CateListKey

	if r.Context.CacheClient == nil {
		r.Error("message", "service.CateListBySort redis connect is err", "err", err.Error())
		return
	}
	cacheRes, err := r.Context.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		cates, err := r.doCacheCateList(cacheKey)
		if err != nil {
			r.Error("message", "service.CateListBySort", "err", err.Error())
		}
		return cates, nil
	}
	if err != nil {
		r.Error("message", "service.CateListBySort", "err", err.Error())
		return
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		res, err = r.doCacheCateList(cacheKey)
		if err != nil {
			r.Error("message", "service.CateListBySort", "err", err.Error())
		}
		return
	}
	err = json.Unmarshal([]byte(cacheRes), &res)
	if err != nil {
		r.Error("message", "service.CateListBySort", "err", err.Error())
		res, err = r.doCacheCateList(cacheKey)
		if err != nil {
			r.Error("message", "service.CateListBySort", "err", err.Error())
		}
	}
	return
}

// Get the all cate
// then set it to cache
func (r *CategoryService) doCacheCateList(cacheKey string) ([]pojos.Category, error) {
	allCates, err := r.allCates()
	if err != nil {
		r.Error("message", "service.CateListBySort", "err", err.Error())
		return nil, err
	}
	var cate pojos.Category
	var cates []pojos.Category
	for _, v := range allCates {
		cate.Cates = v
		cates = append(cates, cate)
	}
	res := r.tree(cates, 0, 0, 0)
	jsonRes, err := json.Marshal(&res)
	if err != nil {
		r.Error("message", "service.CateListBySort", "err", err.Error())
		return nil, err
	}
	err = r.Context.CacheClient.Set(cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {
		r.Error("message", "service.CateListBySort", "err", err.Error())
		return nil, err
	}
	return res, nil
}

// data recursion
func (r *CategoryService) tree(cate []pojos.Category, parent int, level int, key int) []pojos.Category {
	html := "-"
	var data []pojos.Category
	for _, v := range cate {
		var ParentId = v.Cates.ParentId
		var Id = v.Cates.Id
		if ParentId == parent {
			var newHtml string
			if level != 0 {
				newHtml = common.GoRepeat("&nbsp;&nbsp;&nbsp;&nbsp;", level) + "|"
			}
			v.Html = newHtml + common.GoRepeat(html, level)
			data = append(data, v)
			data = r.merge(data, r.tree(cate, Id, level+1, key+1))
		}
	}
	return data
}

// merge two arr
func (r *CategoryService) merge(arr1 []pojos.Category, arr2 []pojos.Category) []pojos.Category {
	for _, val := range arr2 {
		arr1 = append(arr1, val)
	}
	return arr1
}

// Get all cate
// if not exists
// create the default one
func (r *CategoryService) allCates() ([]models.ZCategories, error) {
	cates := make([]models.ZCategories, 0)
	err := r.Context.Db.Table((&models.ZCategories{}).TableName()).
		Find(&cates).
		Error

	if err != nil {
		r.Info("message", "service.AllCates", "err", err.Error())
		return cates, err
	}

	if len(cates) == 0 {
		cateCreate := models.ZCategories{
			Name:        "default",
			DisplayName: "默认分类",
			SeoDesc:     "默认的分类",
			ParentId:    0,
		}
		err := r.Context.Db.Create(&cateCreate).Error
		if err != nil {
			r.Info("message", "service.AllCates", "err", err.Error())
			return cates, err
		}
		if cateCreate.Id < 1 {
			r.Info("message", "service.AllCates", err, "未成功插入数据")
			return cates, errors.New("插入默认分类数据失败")
		}
		err = r.Context.Db.Table((&models.ZCategories{}).TableName()).Find(&cates).Error

		if err != nil {
			r.Info("message", "service.AllCates", "err", err.Error())
			return cates, err
		}

		return cates, nil
	}

	return cates, nil
}

func (r *CategoryService) CateCnt() (cnt int64, err error) {
	err = r.Context.Db.Table((&models.ZCategories{}).TableName()).Count(&cnt).Error
	return
}

func (r *CategoryService) Error(msg ...interface{}) {

}
func (r *CategoryService) Info(msg ...interface{}) {

}
