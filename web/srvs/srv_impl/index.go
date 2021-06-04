/**
 * Created by GoLand.
 * User: zhu
 * Email: ylsc633@gmail.com
 * Date: 2019-05-16
 * Time: 20:17
 */
package srv_impl

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type IndexType string

const (
	IndexTypeOne   IndexType = "tag"
	IndexTypeTwo   IndexType = "cate"
	IndexTypeThree IndexType = "default"
)

type IndexService struct {
	base.ServiceBase
}

func NewIndexService(context ...*base.Context) (p *IndexService) {
	p = &IndexService{}
	p.SetContext(context...)
	return
}

func (r *IndexService) IndexPost(page string, limit string, indexType IndexType, name string) (indexPostIndex wrappers.IndexPostList, err error) {
	var postKey string
	switch indexType {
	case IndexTypeOne:
		postKey = common.Conf.TagPostIndexKey
	case IndexTypeTwo:
		postKey = common.Conf.CatePostIndexKey
	case IndexTypeThree:
		postKey = common.Conf.PostIndexKey
		name = "default"
	default:
		postKey = common.Conf.PostIndexKey
	}

	field := ":name:" + name + ":page:" + page + ":limit:" + limit
	cacheRes, err := r.Context.CacheClient.HGet(r.Context.GinContext.Request.Context(),postKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.IndexPost",
				"err":     err,
			}, )
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.index.IndexPost1",
			"err":     err,
		}, )
		return indexPostIndex, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.IndexPost2",
				"err":     err,
			}, )
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &indexPostIndex)
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.index.IndexPost3",
			"err":     err,
		}, )
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.IndexPost4",
				"err":     err,
			}, )
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}
	return
}

func (r *IndexService) doCacheIndexPostList(cacheKey string, field string, indexType IndexType, name string, queryPage string, queryLimit string) (res wrappers.IndexPostList, err error) {
	limit, offset := common.Offset(queryPage, queryLimit)
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message":    "service.index.doCacheIndexPostList",
			"cacheKey":   cacheKey,
			"field":      field,
			"indexType":  indexType,
			"name":       name,
			"queryPage":  queryPage,
			"queryLimit": queryLimit,
			"err":        err,
		}, )
		return
	}
	var postList *[]wrappers.ConsolePostList
	var postCount int64
	var dba *gorm.DB

	postSrv := NewConsolePostService(r.Context)

	switch indexType {
	case IndexTypeOne:
		tag := new(models.ZTags)
		err = r.Context.Db.Table((&models.ZTags{}).TableName()).Where("Name = ?", name).Find(tag).Error
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message":    "service.index.doCacheIndexPostList1",
				"cacheKey":   cacheKey,
				"field":      field,
				"indexType":  indexType,
				"name":       name,
				"queryPage":  queryPage,
				"queryLimit": queryLimit,
				"err":        err,
			}, )
			return
		}
		postList, err = postSrv.PostTagList(tag.Id, limit, offset)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message":    "service.index.doCacheIndexPostList2",
				"cacheKey":   cacheKey,
				"field":      field,
				"indexType":  indexType,
				"name":       name,
				"queryPage":  queryPage,
				"queryLimit": queryLimit,
				"err":        err,
			}, )
			return
		}
		postCount, err = postSrv.PostTagListCount(tag.Id, limit, offset)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.doCacheIndexPostList3",
				"tagId":   tag.Id,
				"limit":   limit,
				"offset":  offset,
				"err":     err,
			}, )
			return
		}
	case IndexTypeTwo:
		cate := new(models.ZCategories)
		err = r.Context.Db.Table((&models.ZCategories{}).TableName()).
			Where("Name = ?", name).
			Find(cate).Error
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.doCacheIndexPostList4",
				"name":    name,
				"err":     err,
			}, )
			return
		}
		postList, err = postSrv.PostCateList(cate.Id, limit, offset)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.doCacheIndexPostList5",
				"limit":   limit,
				"cateId":  cate.Id,
				"offset":  offset,
				"err":     err,
			}, )
			return
		}
		postCount, err = postSrv.PostCateListCount(cate.Id, limit, offset)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.doCacheIndexPostList6",
				"limit":   limit,
				"cateId":  cate.Id,
				"offset":  offset,
				"err":     err,
			}, )
			return
		}
	case IndexTypeThree:

		dba, postCount, err = postSrv.ConsolePostCount(limit, offset, false)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.doCacheIndexPostList7",
				"limit":   limit,
				"offset":  offset,
				"err":     err,
			}, )
			return
		}
		if postCount > 0 {
			postList, err = postSrv.ConsolePostIndex(dba, limit, offset, false)
			if err != nil {
				r.Context.Error(map[string]interface{}{
					"message": "service.index.doCacheIndexPostList8",
					"limit":   limit,
					"offset":  offset,
					"err":     err,
				}, )
				return
			}
		}
	default:
		dba, postCount, err = postSrv.ConsolePostCount(limit, offset, false)
		if err != nil {
			r.Context.Error(map[string]interface{}{
				"message": "service.index.doCacheIndexPostList9",
				"limit":   limit,
				"offset":  offset,
				"err":     err,
			}, )
			return
		}
		if postCount > 0 {
			postList, err = postSrv.ConsolePostIndex(dba, limit, offset, false)
			if err != nil {
				r.Context.Error(map[string]interface{}{
					"message": "service.index.doCacheIndexPostList10",
					"limit":   limit,
					"offset":  offset,
					"err":     err,
				}, )
				return
			}
		}
	}

	paginate := utils.MyPaginate(postCount, limit, queryPageInt)

	res = wrappers.IndexPostList{
		PostListArr: postList,
		Paginate:    paginate,
	}

	jsonRes, err := json.Marshal(&res)
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.index.doCacheIndexPostList11",
			"err":     err,
		}, )
		return
	}
	err = r.Context.CacheClient.HSet(r.Context.GinContext.Request.Context(),cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message":  "service.index.doCacheIndexPostList12",
			"cacheKey": cacheKey,
			"field":    field,
			"jsonRes":  jsonRes,
			"err":      err,
		}, )
		return
	}
	return
}

func (r *IndexService) IndexPostDetail(postIdStr string) (postDetail wrappers.IndexPostDetail, err error) {
	cacheKey := common.Conf.PostDetailIndexKey
	field := ":post:id:" + postIdStr

	postSrv := NewConsolePostService()
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{"message": "service.Index.IndexPostDetail",
				"err": err.Error(),
			},
		)
		return
	}
	cacheRes, err := r.Context.CacheClient.HGet(r.Context.GinContext.Request.Context(),cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		postDetail, err := r.doCacheIndexPostDetail(postSrv, cacheKey, field, postIdInt)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{"message": "service.index.IndexPostDetail",
					"err":       err.Error(),
					"cacheKey":  cacheKey,
					"field":     field,
					"postIdInt": postIdInt,
				},
			)
			return postDetail, err
		}
		return postDetail, nil
	}
	if err != nil {
		r.Context.Error(
			map[string]interface{}{"message": "service.index.IndexPostDetail",
				"err":      err.Error(),
				"cacheKey": cacheKey,
				"field":    field,
			},
		)
		return postDetail, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		postDetail, err = r.doCacheIndexPostDetail(postSrv, cacheKey, field, postIdInt)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{"message": "service.index.IndexPostDetail",
					"err":       err.Error(),
					"postSrv":   postSrv,
					"cacheKey":  cacheKey,
					"field":     field,
					"postIdInt": postIdInt,
				},
			)
			return postDetail, err
		}
		return postDetail, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &postDetail)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{"message": "service.index.IndexPostDetail2",
				"err":      err.Error(),
				"cacheRes": cacheRes,
			},
		)
		postDetail, err = r.doCacheIndexPostDetail(postSrv, cacheKey, field, postIdInt)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{"message": "service.index.IndexPostDetail3",
					"err":       err.Error(),
					"postSrv":   postSrv,
					"cacheKey":  cacheKey,
					"field":     field,
					"postIdInt": postIdInt,
				},
			)
			return postDetail, err
		}
		return postDetail, nil
	}
	return

}
func (r *IndexService) doCacheIndexPostDetail(postSrv *ConsolePostService, cacheKey string, field string, postIdInt int) (postDetail wrappers.IndexPostDetail, err error) {
	postDetail, err = postSrv.IndexPostDetailDao(postIdInt)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":   "service.doCacheIndexPostDetail",
				"err":       err.Error(),
				"cacheKey":  cacheKey,
				"field":     field,
				"postIdInt": postIdInt,
			},
		)
		return
	}
	jsonRes, err := json.Marshal(&postDetail)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message": "service.doCacheIndexPostDetail",
				"err":     err.Error(),
			},
		)
		return
	}
	err = r.Context.CacheClient.HSet(r.Context.GinContext.Request.Context(),cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.doCacheIndexPostDetail2",
				"err":      err.Error(),
				"cacheKey": cacheKey,
				"field":    field,
			},
		)
		return
	}
	return
}

func (r *IndexService) PostViewAdd(postIdStr string) {
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":   "service.Index.PostViewAdd",
				"err":       err.Error(),
				"postIdStr": postIdStr,
			},
		)
		return
	}
	err = r.Context.Db.Table((&models.ZPostViews{}).TableName()).
		Where("post_id = ?", postIdInt).Update("num", gorm.Expr("num + ?", 1)).Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":   "service.Index.PostViewAdd",
				"err":       err.Error(),
				"postIdInt": postIdInt,
			},
		)
		return
	}
	return
}

func (r *ConsolePostService) PostArchives() (archivesList map[string][]*models.ZPosts, err error) {
	cacheKey := common.Conf.ArchivesKey
	field := ":all:"

	cacheRes, err := r.Context.CacheClient.HGet(r.Context.GinContext.Request.Context(),cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{
					"message":  "service.Index.PostArchives",
					"err":      err.Error(),
					"cacheKey": cacheKey,
					"field":    field,
				},
			)
			return archivesList, err
		}
		return archivesList, nil
	}
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.Index.PostArchives1",
				"err":      err.Error(),
				"cacheKey": cacheKey,
				"field":    field,
			},
		)
		return archivesList, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{
					"message":  "service.Index.PostArchives2",
					"err":      err.Error(),
					"cacheKey": cacheKey,
					"field":    field,
				},
			)
			return archivesList, err
		}
		return archivesList, nil
	}

	archivesList = make(map[string][]*models.ZPosts)
	err = json.Unmarshal([]byte(cacheRes), &archivesList)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.Index.PostArchives3",
				"err":      err.Error(),
				"cacheRes": cacheRes,
			},
		)
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Context.Error(
				map[string]interface{}{
					"message":  "service.Index.PostArchives4",
					"err":      err.Error(),
					"cacheRes": cacheRes,
					"field":    field,
				},
			)
			return archivesList, err
		}
		return archivesList, nil
	}
	return
}

func (r *ConsolePostService) doCacheArchives(cacheKey string, field string) (archivesList map[string][]*models.ZPosts, err error) {
	posts := make([]*models.ZPosts, 0)
	err = r.Context.Db.Table((&models.ZPosts{}).TableName()).
		Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
		Order("created_at desc").
		Find(&posts).Error
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.Index.doCacheArchives",
				"err":      err.Error(),
				"cacheKey": cacheKey,
				"field":    field,
			},
		)
		return
	}
	archivesList = make(map[string][]*models.ZPosts)
	for _, v := range posts {
		date := v.CreatedAt.Format("2006-01")
		archivesList[date] = append(archivesList[date], v)
	}

	jsonRes, err := json.Marshal(&archivesList)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.Index.doCacheArchives3",
				"err":      err.Error(),
				"cacheKey": cacheKey,
				"field":    field,
			},
		)
		return
	}
	err = r.Context.CacheClient.HSet(r.Context.GinContext.Request.Context(),cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Context.Error(
			map[string]interface{}{
				"message":  "service.Index.doCacheArchives4",
				"err":      err.Error(),
				"cacheKey": cacheKey,
				"field":    field,
			},
		)
		return
	}
	return
}
