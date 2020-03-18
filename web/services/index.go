/**
 * Created by GoLand.
 * User: zhu
 * Email: ylsc633@gmail.com
 * Date: 2019-05-16
 * Time: 20:17
 */
package services

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/juetun/study/app-dashboard/lib/app_log"
	"github.com/juetun/study/app-dashboard/lib/base"
	"github.com/juetun/study/app-dashboard/lib/common"
	"github.com/juetun/study/app-dashboard/web/models"
	"github.com/juetun/study/app-dashboard/web/pojos"
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

func NewIndexService() *IndexService {
	return &IndexService{}
}

func (r *IndexService) IndexPost(page string, limit string, indexType IndexType, name string) (indexPostIndex pojos.IndexPostList, err error) {
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
	cacheRes, err := r.CacheClient.HGet(postKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Log.Errorln("message", "service.index.IndexPost", "err", err.Error())
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}
	if err != nil {
		r.Log.Errorln("message", "service.index.IndexPost", "err", err.Error())
		return indexPostIndex, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Log.Errorln("message", "service.index.IndexPost", "err", err.Error())
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &indexPostIndex)
	if err != nil {
		r.Log.Errorln("message", "service.index.IndexPost", "err", err.Error())
		indexPostIndex, err := r.doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {
			r.Log.Errorln("message", "service.index.IndexPost", "err", err.Error())
			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}
	return
}

func (r *IndexService) doCacheIndexPostList(cacheKey string, field string, indexType IndexType, name string, queryPage string, queryLimit string) (res pojos.IndexPostList, err error) {
	limit, offset := common.Offset(queryPage, queryLimit)
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
		return
	}
	var postList []*pojos.ConsolePostList
	var postCount int64
	postSrv := NewConsolePostService()
	switch indexType {
	case IndexTypeOne:
		tag := new(models.ZTags)
		_, err = r.Db.Where("Name = ?", name).Get(tag)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postList, err = postSrv.PostTagList(tag.Id, limit, offset)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postCount, err = postSrv.PostTagListCount(tag.Id, limit, offset)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
	case IndexTypeTwo:
		cate := new(models.ZCategories)
		_, err = r.Db.Where("Name = ?", name).Get(cate)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postList, err = postSrv.PostCateList(cate.Id, limit, offset)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postCount, err = postSrv.PostCateListCount(cate.Id, limit, offset)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
	case IndexTypeThree:
		postList, err = postSrv.ConsolePostIndex(limit, offset, false)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
		postCount, err = postSrv.ConsolePostCount(limit, offset, false)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
	default:
		postList, err = postSrv.ConsolePostIndex(limit, offset, false)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}

		postCount, err = postSrv.ConsolePostCount(limit, offset, false)
		if err != nil {
			r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
			return
		}
	}

	paginate := common.MyPaginate(postCount, limit, queryPageInt)

	res = pojos.IndexPostList{
		PostListArr: postList,
		Paginate:    paginate,
	}

	jsonRes, err := json.Marshal(&res)
	if err != nil {
		r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
		return
	}
	err = r.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Log.Errorln("message", "service.index.doCacheIndexPostList", "err", err.Error())
		return
	}
	return
}

func (r *IndexService) IndexPostDetail(postIdStr string) (postDetail pojos.IndexPostDetail, err error) {
	cacheKey := common.Conf.PostDetailIndexKey
	field := ":post:id:" + postIdStr

	postSrv:=NewConsolePostService()
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		app_log.GetLog().Errorln("message", "service.Index.IndexPostDetail", "err", err.Error())
		return
	}
	cacheRes, err := r.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		postDetail, err := r.doCacheIndexPostDetail(postSrv,cacheKey, field, postIdInt)
		if err != nil {
			r.Log.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
			return postDetail, err
		}
		return postDetail, nil
	}
	if err != nil {
		r.Log.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
		return postDetail, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		postDetail, err = r.doCacheIndexPostDetail(postSrv,cacheKey, field, postIdInt)
		if err != nil {
			r.Log.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
			return postDetail, err
		}
		return postDetail, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &postDetail)
	if err != nil {
		r.Log.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
		postDetail, err = r.doCacheIndexPostDetail(postSrv,cacheKey, field, postIdInt)
		if err != nil {
			r.Log.Errorln("message", "service.index.IndexPostDetail", "err", err.Error())
			return postDetail, err
		}
		return postDetail, nil
	}
	return

}
func (r *IndexService) doCacheIndexPostDetail(postSrv *ConsolePostService, cacheKey string, field string, postIdInt int) (postDetail pojos.IndexPostDetail, err error) {
	postDetail, err = postSrv.IndexPostDetailDao(postIdInt)
	if err != nil {
		r.Log.Errorln("message", "service.doCacheIndexPostDetail", "err", err.Error())
		return
	}
	jsonRes, err := json.Marshal(&postDetail)
	if err != nil {
		r.Log.Errorln("message", "service.index.doCacheIndexPostDetail", "err", err.Error())
		return
	}
	err = r.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Log.Errorln("message", "service.index.doCacheIndexPostDetail", "err", err.Error())
		return
	}
	return
}

func (r *IndexService) PostViewAdd(postIdStr string) {
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {
		r.Log.Errorln("message", "service.Index.PostViewAdd", "err", err.Error())
		return
	}
	_, err = r.Db.Where("post_id = ?", postIdInt).Incr("num").Update(models.ZPostViews{})
	if err != nil {
		r.Log.Errorln("message", "service.Index.PostViewAdd", "err", err.Error())
		return
	}
	return
}

func (r *ConsolePostService) PostArchives() (archivesList map[string][]*models.ZPosts, err error) {
	cacheKey := common.Conf.ArchivesKey
	field := ":all:"

	cacheRes, err := r.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Log.Errorln("message", "service.index.PostArchives", "err", err.Error())
			return archivesList, err
		}
		return archivesList, nil
	}
	if err != nil {
		r.Log.Errorln("message", "service.index.PostArchives", "err", err.Error())
		return archivesList, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Log.Errorln("message", "service.index.PostArchives", "err", err.Error())
			return archivesList, err
		}
		return archivesList, nil
	}

	archivesList = make(map[string][]*models.ZPosts)
	err = json.Unmarshal([]byte(cacheRes), &archivesList)
	if err != nil {
		r.Log.Errorln("message", "service.index.PostArchives", "err", err.Error())
		archivesList, err := r.doCacheArchives(cacheKey, field)
		if err != nil {
			r.Log.Errorln("message", "service.index.PostArchives", "err", err.Error())
			return archivesList, err
		}
		return archivesList, nil
	}
	return
}

func (r *ConsolePostService) doCacheArchives(cacheKey string, field string) (archivesList map[string][]*models.ZPosts, err error) {
	posts := make([]*models.ZPosts, 0)
	err = r.Db.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Desc("created_at").Find(&posts)
	if err != nil {
		r.Log.Errorln("message", "service.Index.doCacheArchives", "err", err.Error())
		return
	}
	archivesList = make(map[string][]*models.ZPosts)
	for _, v := range posts {
		date := v.CreatedAt.Format("2006-01")
		archivesList[date] = append(archivesList[date], v)
	}

	jsonRes, err := json.Marshal(&archivesList)
	if err != nil {
		r.Log.Errorln("message", "service.index.doCacheArchives", "err", err.Error())
		return
	}
	err = r.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {
		r.Log.Errorln("message", "service.index.doCacheArchives", "err", err.Error())
		return
	}
	return
}
