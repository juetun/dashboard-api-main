/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-12
 * Time: 21:03
 */
package services

import (
	"errors"
	"html/template"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services/bak"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type ConsolePostService struct {
	Db *xorm.Engine
	base.ServiceBase
}

func NewConsolePostService() *ConsolePostService {
	return &ConsolePostService{}
}

func (r *ConsolePostService) ConsolePostCount(limit int, offset int, isTrash bool) (count int64, err error) {
	post := new(models.ZPosts)
	if isTrash {
		count, err = r.Db.Unscoped().Where("`deleted_at` IS NOT NULL OR `deleted_at`=?", "0001-01-01 00:00:00").Desc("id").Limit(limit, offset).Count(post)
	} else {
		count, err = r.Db.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Desc("id").Limit(limit, offset).Count(post)
	}
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.ConsolePostCount", "err": err.Error(),
		})
		return 0, err
	}
	return count, nil
}

func (r *ConsolePostService) ConsolePostIndex(limit, offset int, isTrash bool) (postListArr []*pojos.ConsolePostList, err error) {
	post := new(models.ZPosts)

	var rows *xorm.Rows
	if isTrash {
		rows, err = r.Db.Unscoped().Where("`deleted_at` IS NOT NULL OR `deleted_at`=?", "0001-01-01 00:00:00").Desc("id").Limit(limit, offset).Rows(post)
	} else {
		rows, err = r.Db.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Desc("id").Limit(limit, offset).Rows(post)
	}

	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.ConsolePostIndex", "err": err.Error(),
		})
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		// post
		post := new(models.ZPosts)
		err = rows.Scan(post)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.ConsolePostIndex", "err": err.Error(),
			})
			return nil, err
		}

		consolePost := pojos.ConsolePost{
			Id:        post.Id,
			Uid:       post.Uid,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		// category
		cates, err := bak.GetPostCateByPostId(post.Id)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.ConsolePostIndex", "err": err.Error(),
			})
			return nil, err
		}
		consoleCate := pojos.ConsoleCate{
			Id:          cates.Id,
			Name:        cates.Name,
			DisplayName: cates.DisplayName,
			SeoDesc:     cates.SeoDesc,
		}

		// tag
		tagIds, err := bak.GetPostTagsByPostId(post.Id)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.ConsolePostIndex", "err": err.Error(),
			})
			return nil, err
		}
		tags, err := bak.GetTagsByIds(tagIds)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.ConsolePostIndex", "err": err.Error(),
			})
			return nil, err
		}
		var consoleTags []pojos.ConsoleTag
		for _, v := range tags {
			consoleTag := pojos.ConsoleTag{
				Id:          v.Id,
				Name:        v.Name,
				DisplayName: v.DisplayName,
				SeoDesc:     v.SeoDesc,
				Num:         v.Num,
			}
			consoleTags = append(consoleTags, consoleTag)
		}

		// view
		view, err := r.PostView(post.Id)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.ConsolePostIndex", "err": err.Error(),
			})
			return nil, err
		}
		consoleView := pojos.ConsoleView{
			Num: view.Num,
		}

		// user
		user, err := bak.GetUserById(post.UserId)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.ConsolePostIndex", "err": err.Error(),
			})
			return nil, err
		}
		consoleUser := pojos.ConsoleUser{
			Id:     user.Id,
			Name:   user.Name,
			Email:  user.Email,
			Status: user.Status,
		}

		postList := pojos.ConsolePostList{
			Post:     consolePost,
			Category: consoleCate,
			Tags:     consoleTags,
			View:     consoleView,
			Author:   consoleUser,
		}
		postListArr = append(postListArr, &postList)
	}

	return postListArr, nil
}

func (r *ConsolePostService) PostView(postId int) (*models.ZPostViews, error) {
	postV := new(models.ZPostViews)
	_, err := r.Db.Where("post_id = ?", postId).Cols("num").Get(postV)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostView", "err": err.Error(),
		})
	}
	return postV, nil
}

func (r *ConsolePostService) PostStore(ps pojos.PostStore, userId int) {
	postCreate := &models.ZPosts{
		Title:    ps.Title,
		UserId:   userId,
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	postCreate.Content = string(html)

	session := r.Db.NewSession()
	defer session.Close()
	affected, err := session.Insert(postCreate)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostStore", "err": err.Error(),
		})
		_ = session.Rollback()
		return
	}
	if affected < 1 {

		r.Log.Info(map[string]string{
			"message": "service.PostStore", "err": "post store no succeed",
		})
		_ = session.Rollback()
		return
	}

	if ps.Category > 0 {
		postCateCreate := models.ZPostCate{
			PostId: postCreate.Id,
			CateId: ps.Category,
		}
		affected, err := session.Insert(postCateCreate)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.PostStore", "err": err.Error(),
			})
			_ = session.Rollback()
			return
		}

		if affected < 1 {
			r.Log.Error(map[string]string{
				"message": "service.PostStore", "err": "post cate store not succeed",
			})
			_ = session.Rollback()
			return
		}
	}

	if len(ps.Tags) > 0 {
		for _, v := range ps.Tags {
			postTagCreate := models.ZPostTag{
				PostId: postCreate.Id,
				TagId:  v,
			}
			affected, err := session.Insert(postTagCreate)
			if err != nil {
				r.Log.Error(map[string]string{
					"message": "service.PostStore post tag insert err", "err": err.Error(),
				})
				_ = session.Rollback()
				return
			}
			if affected < 1 {
				r.Log.Error(map[string]string{
					"message": "service.PostStore", "err": "post tag store not succeed",
				})
				_ = session.Rollback()
				return
			}

			affected, err = session.ID(v).Incr("num").Update(models.ZTags{})
			if err != nil {
				r.Log.Error(map[string]string{
					"message": "service.PostStore post tag incr err", "err": err.Error(),
				})
				_ = session.Rollback()
				return
			}
			if affected < 1 {
				r.Log.Error(map[string]string{
					"message": "service.PostStore", "err": "post tag incr not succeed",
				})
				_ = session.Rollback()
				return
			}
		}
	}

	postView := models.ZPostViews{
		PostId: postCreate.Id,
		Num:    1,
	}

	affected, err = session.Insert(postView)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostStore", "err": err.Error(),
		})
		_ = session.Rollback()
		return
	}

	if affected < 1 {
		r.Log.Error(map[string]string{
			"message": "service.PostStore", "err": "post view store no succeed",
		})
		_ = session.Rollback()
		return
	}

	_ = session.Commit()

	uid, err := common.ZHashId.Encode([]int{postCreate.Id})
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostStore create uid error", "err": err.Error(),
		})
		return
	}

	newPostCreate := models.ZPosts{
		Uid: uid,
	}
	affected, err = session.Where("id = ?", postCreate.Id).Update(newPostCreate)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostStore",
			"err":     err.Error(),
		})
		return
	}

	if affected < 1 {
		r.Log.Error(map[string]string{
			"message": "service.PostStore",
			"err":     "post view store not succeed",
		})
		return
	}

	return
}

func (r *ConsolePostService) PostDetail(postId int) (p *models.ZPosts, err error) {
	post := new(models.ZPosts)
	_, err = r.Db.Where("id = ?", postId).Get(post)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostDetail",
			"err":     err.Error(),
		})
		return post, err
	}
	return post, nil
}

func (r *ConsolePostService) IndexPostDetailDao(postId int) (postDetail pojos.IndexPostDetail, err error) {
	post := new(models.ZPosts)
	_, err = r.Db.Where("id = ?", postId).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Get(post)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	if post.Id <= 0 {
		return postDetail, errors.New("Post do not exists ")
	}
	Post := pojos.IndexPost{
		Id:        post.Id,
		Uid:       post.Uid,
		Title:     post.Title,
		Summary:   post.Summary,
		Original:  post.Original,
		Content:   template.HTML(post.Content),
		Password:  post.Password,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	tags, err := r.PostIdTags(postId)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	var Tags []pojos.ConsoleTag
	for _, v := range tags {
		consoleTag := pojos.ConsoleTag{
			Id:          v.Id,
			Name:        v.Name,
			DisplayName: v.DisplayName,
			SeoDesc:     v.SeoDesc,
			Num:         v.Num,
		}
		Tags = append(Tags, consoleTag)
	}

	cate, err := r.PostCates(postId)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	Cate := pojos.ConsoleCate{
		Id:          cate.Id,
		Name:        cate.Name,
		DisplayName: cate.DisplayName,
		SeoDesc:     cate.SeoDesc,
	}

	// view
	view, err := r.PostView(post.Id)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	View := pojos.ConsoleView{
		Num: view.Num,
	}

	// user
	user, err := bak.GetUserById(post.UserId)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	Author := pojos.ConsoleUser{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}

	// last post
	lastPost, err := r.LastPost(postId)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}

	// next post
	nextPost, err := r.NextPost(postId)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}

	postDetail = pojos.IndexPostDetail{
		Post:     Post,
		Category: Cate,
		Tags:     Tags,
		View:     View,
		Author:   Author,
		LastPost: lastPost,
		NextPost: nextPost,
	}

	return postDetail, nil
}

func (r *ConsolePostService) LastPost(postId int) (post *models.ZPosts, err error) {
	post = new(models.ZPosts)
	_, err = r.Db.Where("id < ?", postId).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Desc("id").Get(post)
	return
}

func (r *ConsolePostService) NextPost(postId int) (post *models.ZPosts, err error) {
	post = new(models.ZPosts)
	_, err = r.Db.Where("id > ?", postId).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Asc("id").Get(post)
	return
}

func (r *ConsolePostService) PostIdTags(postId int) (tags []*models.ZTags, err error) {
	tagIds, err := r.PostIdTag(postId)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostIdTags",
			"err":     err.Error(),
		})
		return
	}
	// tags = make([]models.ZTags,0)
	err = r.Db.In("id", tagIds).Find(&tags)
	return
}

func (r *ConsolePostService) PostIdTag(postId int) (tagIds []int, err error) {
	postTag := make([]models.ZPostTag, 0)
	err = r.Db.Where("post_id = ?", postId).Find(&postTag)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostIdTag",
			"err":     err.Error(),
		})
		return
	}

	for _, v := range postTag {
		tagIds = append(tagIds, v.TagId)
	}
	return tagIds, nil
}

func (r *ConsolePostService) PostCates(postId int) (cate *models.ZCategories, err error) {
	srv := NewCategoryService()
	cateId, err := srv.PostCate(postId)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return
	}
	cate = new(models.ZCategories)
	_, err = r.Db.Id(cateId).Get(cate)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return
	}
	return
}

func (r *ConsolePostService) PostUpdate(postId int, ps pojos.PostStore) {
	postUpdate := &models.ZPosts{
		Title:    ps.Title,
		UserId:   1,
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	postUpdate.Content = string(html)
	session := r.Db.NewSession()
	defer session.Close()
	affected, err := session.Where("id = ?", postId).Update(postUpdate)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     err.Error(),
		})
		_ = session.Rollback()
		return
	}
	if affected < 1 {
		r.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     err.Error(),
		})
		_ = session.Rollback()
		return
	}

	postCate := new(models.ZPostCate)
	_, err = session.Where("post_id = ?", postId).Delete(postCate)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     "post cate delete no succeed",
		})
		_ = session.Rollback()
		return
	}

	if ps.Category > 0 {
		postCateCreate := models.ZPostCate{
			PostId: postId,
			CateId: ps.Category,
		}

		affected, err := session.Insert(postCateCreate)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.PostUpdate",
				"err":     err.Error(),
			})
			_ = session.Rollback()
			return
		}

		if affected < 1 {
			r.Log.Error(map[string]string{
				"message": "service.PostUpdate",
				"err":     "post cate update no succeed",
			})
			_ = session.Rollback()
			return
		}
	}

	postTag := make([]models.ZPostTag, 0)
	err = session.Where("post_id = ?", postId).Find(&postTag)

	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     "get post tag  no succeed",
		})
		_ = session.Rollback()
		return
	}

	if len(postTag) > 0 {
		for _, v := range postTag {
			affected, err = session.ID(v.TagId).Decr("num").Update(models.ZTags{})
			if err != nil {
				r.Log.Error(map[string]string{
					"message": "service.PostUpdate post tag decr  err",
					"err":     err.Error(),
				})
				_ = session.Rollback()
				return
			}
			if affected < 1 {
				r.Log.Error(map[string]string{
					"message": "service.PostUpdate",
					"err":     "post cate decr no succeed",
				})
				_ = session.Rollback()
				return
			}
		}

		_, err = session.Where("post_id = ?", postId).Delete(new(models.ZPostTag))

		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.PostUpdate",
				"err":     "delete post tag  no succeed",
			})
			_ = session.Rollback()
			return
		}
	}

	if len(ps.Tags) > 0 {
		for _, v := range ps.Tags {
			postTagCreate := models.ZPostTag{
				PostId: postId,
				TagId:  v,
			}
			affected, err := session.Insert(postTagCreate)
			if err != nil {
				r.Log.Error(map[string]string{
					"message": "service.PostUpdate post tag insert err",
					"err":     err.Error(),
				})
				_ = session.Rollback()
				return
			}
			if affected < 1 {
				r.Log.Error(map[string]string{
					"message": "service.PostUpdate",
					"err":     "post cate update no succeed",
				})
				_ = session.Rollback()
				return
			}
			affected, err = session.ID(v).Incr("num").Update(models.ZTags{})
			if err != nil {
				r.Log.Error(map[string]string{
					"message": "service.PostStore post tag incr err",
					"err":     err.Error(),
				})
				_ = session.Rollback()
				return
			}
			if affected < 1 {
				r.Log.Error(map[string]string{
					"message": "service.PostStore",
					"err":     "post tag incr no succeed",
				})
				_ = session.Rollback()
				return
			}
		}
	}
	_ = session.Commit()

	return
}

func (r *ConsolePostService) PostDestroy(postId int) (bool, error) {
	post := new(models.ZPosts)
	toBeCharge := time.Now().Format(time.RFC3339)
	timeLayout := time.RFC3339
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, toBeCharge, loc)
	post.DeletedAt = &theTime
	_, err = r.Db.Id(postId).Update(post)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostDestroy",
			"err":     err.Error(),
		})
		return false, err
	}
	return true, nil
}

func (r *ConsolePostService) PostUnTrash(postId int) (bool, error) {
	post := new(models.ZPosts)
	theTime, _ := time.Parse("2006-01-02 15:04:05", "")
	post.DeletedAt = &theTime
	_, err := r.Db.Id(postId).Update(post)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostUnTrash",
			"err":     err.Error(),
		})
		return false, err
	}
	return true, nil
}

func (r *ConsolePostService) PostCnt() (cnt int64, err error) {
	post := new(models.ZPosts)
	cnt, err = r.Db.Count(post)
	return
}

func (r *ConsolePostService) PostTagListCount(tagId int, limit int, offset int) (count int64, err error) {
	postTag := new(models.ZPostTag)
	count, err = r.Db.Where("tag_id = ?", tagId).Desc("id").Limit(limit, offset).Count(postTag)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostTagListCount",
			"err":     err.Error(),
		})
		return 0, err
	}
	return
}

func (r *ConsolePostService) PostTagList(tagId int, limit int, offset int) (postListArr []*pojos.ConsolePostList, err error) {
	postTag := new(models.ZPostTag)
	rows, err := r.Db.Where("tag_id = ?", tagId).Desc("id").Limit(limit, offset).Rows(postTag)

	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.Index.PostTagList",
			"err":     err.Error(),
		})
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		// post
		postTag := new(models.ZPostTag)
		err = rows.Scan(postTag)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.Index.PostTagList",
				"err":     err.Error(),
			})
			return nil, err
		}

		post := new(models.ZPosts)
		_, err = r.Db.Id(postTag.PostId).Get(post)

		consolePost := pojos.ConsolePost{
			Id:        post.Id,
			Uid:       post.Uid,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		postList := pojos.ConsolePostList{
			Post: consolePost,
		}
		postListArr = append(postListArr, &postList)
	}

	return postListArr, nil
}

func (r *ConsolePostService) PostCateListCount(cateId int, limit int, offset int) (count int64, err error) {
	postCate := new(models.ZPostCate)
	count, err = r.Db.Where("cate_id = ?", cateId).Desc("id").Limit(limit, offset).Count(postCate)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.PostCateListCount",
			"err":     err.Error(),
		})
		return 0, err
	}
	return
}

func (r *ConsolePostService) PostCateList(cateId int, limit int, offset int) (postListArr []*pojos.ConsolePostList, err error) {
	postCate := new(models.ZPostCate)
	rows, err := r.Db.Where("cate_id = ?", cateId).Desc("id").Limit(limit, offset).Rows(postCate)

	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.Index.PostCateList",
			"err":     err.Error(),
		})
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		// post
		postCate := new(models.ZPostCate)
		err = rows.Scan(postCate)
		if err != nil {
			r.Log.Error(map[string]string{
				"message": "service.Index.PostCateList",
				"err":     err.Error(),
			})
			return nil, err
		}

		post := new(models.ZPosts)
		_, err = r.Db.Id(postCate.PostId).Get(post)

		consolePost := pojos.ConsolePost{
			Id:        post.Id,
			Uid:       post.Uid,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		postList := pojos.ConsolePostList{
			Post: consolePost,
		}
		postListArr = append(postListArr, &postList)
	}

	return postListArr, nil
}
