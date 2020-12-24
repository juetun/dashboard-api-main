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

	"github.com/go-redis/redis"
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
	cacheRes, err := r.Context.CacheClient.HGet(postKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.IndexPost", "err", err.Error())
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.IndexPost", "err", err.Error())
		return indexPostIndex, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.IndexPost", "err", err.Error())
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &indexPostIndex)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.IndexPost", "err", err.Error())
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.IndexPost", "err", err.Error())
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
		r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
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
			r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postList, err = postSrv.PostTagList(tag.Id, limit, offset)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postCount, err = postSrv.PostTagListCount(tag.Id, limit, offset)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
	case IndexTypeTwo:
		cate := new(models.ZCategories)
		err = r.Context.Db.Table((&models.ZCategories{}).TableName()).
			Where("Name = ?", name).
			Find(cate).Error
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postList, err = postSrv.PostCateList(cate.Id, limit, offset)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postCount, err = postSrv.PostCateListCount(cate.Id, limit, offset)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
	case IndexTypeThree:

		dba, postCount, err = postSrv.ConsolePostCount(limit, offset, false)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		if postCount > 0 {
			postList, err = postSrv.ConsolePostIndex(dba, limit, offset, false)
			if err != nil {
				r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
				return
			}
		}
	default:
		dba, postCount, err = postSrv.ConsolePostCount(limit, offset, false)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		if postCount > 0 {
			postList, err = postSrv.ConsolePostIndex(dba, limit, offset, false)
			if err != nil {
				r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
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
		r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
		return
	}
	err = r.Context.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
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
		r.Context.Log.Logger.Errorln("message", "service.Index.IndexPostDetail", "err", err.Error())
		return
	}
	cacheRes, err := r.Context.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		postDetail, err := r.doCacheIndexPostDetail(postSrv, cacheKey, field, postIdInt)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
			return postDetail, err
		}
		return postDetail, nil
	}
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
		return postDetail, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		postDetail, err = r.doCacheIndexPostDetail(postSrv, cacheKey, field, postIdInt)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
			return postDetail, err
		}
		return postDetail, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &postDetail)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
		postDetail, err = r.doCacheIndexPostDetail(postSrv, cacheKey, field, postIdInt)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
			return postDetail, err
		}
		return postDetail, nil
	}
	return

}
func (r *IndexService) doCacheIndexPostDetail(postSrv *ConsolePostService, cacheKey string, field string, postIdInt int) (postDetail wrappers.IndexPostDetail, err error) {
	postDetail, err = postSrv.IndexPostDetailDao(postIdInt)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.doCacheIndexPostDetail", "err", err.Error())
		return
	}
	jsonRes, err := json.Marshal(&postDetail)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostDetail", "err", err.Error())
		return
	}
	err = r.Context.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.doCacheIndexPostDetail", "err", err.Error())
		return
	}
	return
}

func (r *IndexService) PostViewAdd(postIdStr string) {
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.Index.PostViewAdd", "err", err.Error())
		return
	}
	err = r.Context.Db.Table((&models.ZPostViews{}).TableName()).
		Where("post_id = ?", postIdInt).Update("num", gorm.Expr("num + ?", 1)).Error
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.Index.PostViewAdd", "err", err.Error())
		return
	}
	return
}

func (r *ConsolePostService) PostArchives() (archivesList map[string][]*models.ZPosts, err error) {
	cacheKey := common.Conf.ArchivesKey
	field := ":all:"

	cacheRes, err := r.Context.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.PostArchives", "err", err.Error())
			return archivesList, err
		}
		return archivesList, nil
	}
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.PostArchives", "err", err.Error())
		return archivesList, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.PostArchives", "err", err.Error())
			return archivesList, err
		}
		return archivesList, nil
	}

	archivesList = make(map[string][]*models.ZPosts)
	err = json.Unmarshal([]byte(cacheRes), &archivesList)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.PostArchives", "err", err.Error())
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Context.Log.Logger.Errorln("message", "service.index.PostArchives", "err", err.Error())
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
		r.Context.Log.Logger.Errorln("message", "service.Index.doCacheArchives", "err", err.Error())
		return
	}
	archivesList = make(map[string][]*models.ZPosts)
	for _, v := range posts {
		date := v.CreatedAt.Format("2006-01")
		archivesList[date] = append(archivesList[date], v)
	}

	jsonRes, err := json.Marshal(&archivesList)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.doCacheArchives", "err", err.Error())
		return
	}
	err = r.Context.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Context.Log.Logger.Errorln("message", "service.index.doCacheArchives", "err", err.Error())
		return
	}
	return
}
