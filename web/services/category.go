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
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
)

type CategoryService struct {
	base.ServiceBase
}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}
func (r *CategoryService) GetCateById(cateId int) (cate *models.ZCategories, err error) {
	cate = new(models.ZCategories)
	_, err = r.Db.Id(cateId).Get(cate)
	return
}

func (r *CategoryService) GetCateByParentId(parentId int) (cate *models.ZCategories, err error) {
	cate = new(models.ZCategories)
	_, err = r.Db.Where("parent_id = ?", parentId).Get(cate)
	return
}

func (r *CategoryService) DelCateRel(cateId int) {
	session := r.Db.NewSession()
	defer session.Close()
	postCate := new(models.ZPostCate)
	_, err := session.Where("cate_id = ?", cateId).Delete(postCate)
	if err != nil {
		_ = session.Rollback()
		r.Error("message", "service.DelCateRel", "err", err.Error())
		return
	}
	cate := new(models.ZCategories)
	_, err = session.ID(cateId).Delete(cate)
	if err != nil {
		_ = session.Rollback()
		r.Error("message", "service.DelCateRel", "err", err.Error())
		return
	}
	err = session.Commit()
	if err != nil {
		_ = session.Rollback()
		r.Error("message", "service.DelCateRel", "err", err.Error())
		return
	}
	r.CacheClient.Del(common.Conf.CateListKey)
	return
}

func (r *CategoryService) CateStore(cs pojos.CateStore) (bool, error) {

	defaultCate := new(models.ZCategories)
	_, err := r.Db.Where("name = ?", cs.Name).Get(defaultCate)
	if err != nil {
		r.Error("message", "service.CateStore", "err", err.Error())
		return false, err
	}
	if defaultCate.Id > 0 {
		r.Error("message", "service.CateStore", "err", "Cate has exists ")
		return false, errors.New("Tag has exists ")
	}

	if cs.ParentId > 0 {
		cate := new(models.ZCategories)
		_, err := r.Db.Id(cs.ParentId).Get(cate)
		if err != nil {
			r.Error("message", "service.CateStore", "err", err.Error())
			return false, err
		}
		if cate.Id <= 0 {
			r.Error("message", "service.CateStore", err, "The parent id has not data ")
			return false, errors.New("The parent id has not data ")
		}
	}

	cate := models.ZCategories{
		Name:        cs.Name,
		DisplayName: cs.DisplayName,
		SeoDesc:     cs.SeoDesc,
		ParentId:    cs.ParentId,
	}
	_, err = r.Db.Insert(cate)
	if err != nil {
		r.Error("message", "service.CateStore", "err", err.Error())
		return false, err
	}
	r.CacheClient.Del(common.Conf.CateListKey)
	return true, nil
}

func (r *CategoryService) CateUpdate(cateId int, cs pojos.CateStore) (bool, error) {
	cate := new(models.ZCategories)
	if cs.ParentId != 0 {
		res, err := r.Db.Id(cs.ParentId).Get(cate)
		if err != nil {
			r.Error("message", "service.CateUpdate", "err", err.Error())
			return false, err
		}
		if !res || cate.Id < 1 {
			r.Error("message", "service.CateUpdate", err, "the parent id is not exists ")
			return false, errors.New("the parent id is not exists ")
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
	_, err := r.Db.Id(cateId).Cols("name", "display_name", "seo_desc", "parent_id").Update(cateUpdate)
	if err != nil {
		r.Error("message", "service.CateUpdate", "err", err.Error())
		return false, err
	}
	r.CacheClient.Del(common.Conf.CateListKey)
	return true, nil
}

func (r *CategoryService) GetSimilar(beginId []int, resIds []int, level int) (beginId2 []int, resIds2 []int, level2 int) {
	if len(beginId) != 0 {
		cates := make([]*models.ZCategories, 0)
		err := r.Db.In("parent_id", beginId).Find(&cates)
		if err != nil {
			r.Error("message", "service.GetSimilar", err, "the parent id data is not exists ")
			return []int{}, []int{}, 0
		}
		if len(cates) != 0 {
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
		if level == 0 && len(cates) == 0 {
			return beginId, beginId, level
		} else {
			return beginId, resIds, level
		}
	}
	return beginId, resIds, level
}

func (r *CategoryService) GetPostCateByPostId(postId int) (cates *models.ZCategories, err error) {
	postCate := new(models.ZPostCate)
	has, err := r.Db.Cols("cate_id").Where("post_id = ?", postId).Get(postCate)
	if err != nil {
		r.Error("message", "service.GetPostCateByPostId", "err", err.Error())
		return cates, err
	}
	if has {
		cates = new(models.ZCategories)
		has, err = r.Db.Where("id = ?", postCate.CateId).Cols("id", "name", "display_name", "seo_desc").Get(cates)
		if err != nil {
			r.Error("message", "service.GetPostCateByPostId", "err", err.Error())
			return cates, err
		}
		if !has {
			r.Error("message", "service.GetPostCateByPostId", "err", "there has not data")
			return cates, errors.New("can not get the post cate")
		}
	} else {
		r.Error("message", "service.GetPostCateByPostId", "err", "there has not data")
		return cates, errors.New("can not get the post cate")
	}

	return cates, nil

}

func (r *CategoryService) PostCate(postId int) (int, error) {
	postCate := new(models.ZPostCate)
	_, err := r.Db.Where("post_id = ?", postId).Get(postCate)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return 0, err
	}
	return postCate.CateId, nil
}
// Get the cate list what by parent sort
func (r *CategoryService) CateListBySort() ([]pojos.Category, error) {
	cacheKey := common.Conf.CateListKey
	cacheRes, err := r.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		cates, err := r.doCacheCateList(cacheKey)
		if err != nil {
			r.Error("message", "service.CateListBySort", "err", err.Error())
			return nil, err
		}
		return cates, nil
	} else if err != nil {
		r.Error("message", "service.CateListBySort", "err", err.Error())
		return nil, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		cates, err := r.doCacheCateList(cacheKey)
		if err != nil {
			r.Error("message", "service.CateListBySort", "err", err.Error())
			return nil, err
		}
		return cates, nil
	}

	var comCate []pojos.Category
	err = json.Unmarshal([]byte(cacheRes), &comCate)
	if err != nil {
		r.Error("message", "service.CateListBySort", "err", err.Error())
		cates, err := r.doCacheCateList(cacheKey)
		if err != nil {
			r.Error("message", "service.CateListBySort", "err", err.Error())
			return nil, err
		}
		return cates, nil
	}
	return comCate, nil
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
	err = r.CacheClient.Set(cacheKey, jsonRes, time.Duration(common.Conf.DataCacheTimeDuration)*time.Hour).Err()
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
	err := r.Db.Find(&cates)

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
		affected, err := r.Db.Insert(&cateCreate)
		if err != nil {
			r.Info("message", "service.AllCates", "err", err.Error())
			return cates, err
		}
		if affected < 1 {
			r.Info("message", "service.AllCates", err, "未成功插入数据")
			return cates, errors.New("插入默认分类数据失败")
		}
		err = r.Db.Find(&cates)

		if err != nil {
			r.Info("message", "service.AllCates", "err", err.Error())
			return cates, err
		}

		return cates, nil
	}

	return cates, nil
}

func (r *CategoryService) CateCnt() (cnt int64, err error) {
	cate := new(models.ZCategories)
	cnt, err = r.Db.Count(cate)
	return
}

func (r *CategoryService) Error(msg ...interface{}) {

}
func (r *CategoryService) Info(msg ...interface{}) {

}
