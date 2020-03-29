/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-12
 * Time: 21:03
 */
package services

import (
	"database/sql"
	"errors"
	"html/template"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type ConsolePostService struct {
	base.ServiceBase
}

func NewConsolePostService(context ...*base.Context) (p *ConsolePostService) {
	p = &ConsolePostService{}
	p.SetContext(context)
	return
}

func (r *ConsolePostService) ConsolePostCount(limit int, offset int, isTrash bool) (count int64, err error) {
	if isTrash {
		err = r.getDbaTable().
			Unscoped().
			Where("`deleted_at` IS NOT NULL OR `deleted_at`=?", "0001-01-01 00:00:00").
			Order("id desc").
			Limit(limit).
			Offset(offset).
			Count(&count).Error
	} else {
		err = r.getDbaTable().
			Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
			Order("id desc").
			Limit(limit).
			Offset(offset).
			Count(&count).
			Error
	}
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.ConsolePostCount", "err": err.Error(),
		})
		return 0, err
	}
	return count, nil
}

func (r *ConsolePostService) getDbaTable() *gorm.DB {
	return r.Context.Db.Table((&models.ZPosts{}).TableName())
}
func (r *ConsolePostService) ConsolePostIndex(limit, offset int, isTrash bool) (postListArr []*pojos.ConsolePostList, err error) {
	postListArr = make([]*pojos.ConsolePostList, 0)
	var rows *sql.Rows
	if isTrash {
		rows, err = r.getDbaTable().Unscoped().
			// Where("`deleted_at` IS NOT NULL OR `deleted_at`=?", "0001-01-01 00:00:00").
			Order("id desc").
			Limit(limit).
			Offset(offset).
			Rows()
	} else {
		rows, err = r.getDbaTable().
			// Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
			Order("id desc").
			Limit(limit).
			Offset(offset).
			Rows()
	}

	if err != nil {
		r.Context.Log.Error(map[string]string{
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
			r.Context.Log.Error(map[string]string{
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
		srv := NewCategoryService()

		// category
		cates, err := srv.GetPostCateByPostId(post.Id)
		if err != nil {
			r.Context.Log.Error(map[string]string{
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

		srvTag := NewTagService()
		// tag
		tagIds, err := srvTag.GetPostTagsByPostId(post.Id)
		if err != nil {
			r.Context.Log.Error(map[string]string{
				"message": "service.ConsolePostIndex", "err": err.Error(),
			})
			return nil, err
		}

		tags, err := srvTag.GetTagsByIds(tagIds)
		if err != nil {
			r.Context.Log.Error(map[string]string{
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
			r.Context.Log.Error(map[string]string{
				"message": "service.ConsolePostIndex", "err": err.Error(),
			})
			return nil, err
		}
		consoleView := pojos.ConsoleView{
			Num: view.Num,
		}

		srvUser := NewUserService()

		// user
		user, err := srvUser.GetUserById(post.UserId)
		if err != nil {
			r.Context.Log.Error(map[string]string{
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

func (r *ConsolePostService) getZPostViewsDbaTable() *gorm.DB {
	return r.Context.Db.Table((&models.ZPostViews{}).TableName())
}
func (r *ConsolePostService) PostView(postId int) (postV *models.ZPostViews, err error) {
	postV = new(models.ZPostViews)
	err = r.getZPostViewsDbaTable().Where("post_id = ?", postId).Select("num").Find(postV).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
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

	session := r.getDbaTable().Begin()
	defer session.Close()
	err := session.Create(postCreate).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostStore", "err": err.Error(),
		})
		_ = session.Rollback()
		return
	}
	if postCreate.Id < 1 {
		r.Context.Log.Info(map[string]string{
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
		err := session.Create(postCateCreate).Error
		if err != nil {
			r.Context.Log.Error(map[string]string{
				"message": "service.PostStore", "err": err.Error(),
			})
			_ = session.Rollback()
			return
		}

		if postCateCreate.Id < 1 {
			r.Context.Log.Error(map[string]string{
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
			err = session.Create(postTagCreate).Error
			if err != nil {
				r.Context.Log.Error(map[string]string{
					"message": "service.PostStore post tag insert err", "err": err.Error(),
				})
				_ = session.Rollback()
				return
			}
			if postTagCreate.Id < 1 {
				r.Context.Log.Error(map[string]string{
					"message": "service.PostStore", "err": "post tag store not succeed",
				})
				_ = session.Rollback()
				return
			}

			err = session.Where("id=?", v).
				Update("num", gorm.Expr("price + ?", 1)).
				Error
			if err != nil {
				r.Context.Log.Error(map[string]string{
					"message": "service.PostStore post tag incr err", "err": err.Error(),
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

	err = session.Create(postView).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostStore", "err": err.Error(),
		})
		_ = session.Rollback()
		return
	}

	if postView.Id < 1 {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostStore", "err": "post view store no succeed",
		})
		_ = session.Rollback()
		return
	}

	_ = session.Commit()

	uid, err := common.ZHashId.Encode([]int{postCreate.Id})
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostStore create uid error", "err": err.Error(),
		})
		return
	}

	newPostCreate := models.ZPosts{
		Uid: uid,
	}
	err = session.Where("id = ?", postCreate.Id).Update(newPostCreate).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostStore",
			"err":     err.Error(),
		})
		return
	}

	return
}

func (r *ConsolePostService) PostDetail(postId int) (p *models.ZPosts, err error) {
	post := new(models.ZPosts)
	err = r.getDbaTable().Where("id = ?", postId).Find(post).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostDetail",
			"err":     err.Error(),
		})
		return post, err
	}
	return post, nil
}

func (r *ConsolePostService) IndexPostDetailDao(postId int) (postDetail pojos.IndexPostDetail, err error) {
	post := new(models.ZPosts)
	err = r.getDbaTable().
		Where("id = ?", postId).
		Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
		Find(post).
		Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
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
		r.Context.Log.Error(map[string]string{
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
		r.Context.Log.Error(map[string]string{
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
		r.Context.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	View := pojos.ConsoleView{
		Num: view.Num,
	}
	srvUser := NewUserService()
	// user
	user, err := srvUser.GetUserById(post.UserId)
	if err != nil {
		r.Context.Log.Error(map[string]string{
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
		r.Context.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}

	// next post
	nextPost, err := r.NextPost(postId)
	if err != nil {
		r.Context.Log.Error(map[string]string{
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
	err = r.getDbaTable().
		Where("id < ?", postId).
		Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Order("id desc").
		Find(post).
		Error
	return
}

func (r *ConsolePostService) NextPost(postId int) (post *models.ZPosts, err error) {
	post = new(models.ZPosts)
	err = r.getDbaTable().Where("id > ?", postId).
		Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
		Order("id asc").
		Find(post).
		Error
	return
}

func (r *ConsolePostService) PostIdTags(postId int) (tags []*models.ZTags, err error) {
	tagIds, err := r.PostIdTag(postId)
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostIdTags",
			"err":     err.Error(),
		})
		return
	}
	// tags = make([]models.ZTags,0)
	err = r.Context.Db.Table((&models.ZTags{}).TableName()).
		Where("id in (?)", tagIds).
		Find(&tags).
		Error
	return
}

func (r *ConsolePostService) PostIdTag(postId int) (tagIds []int, err error) {
	postTag := make([]models.ZPostTag, 0)
	err = r.Context.Db.
		Where("post_id = ?", postId).
		Find(&postTag).
		Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
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
		r.Context.Log.Error(map[string]string{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return
	}
	cate = new(models.ZCategories)
	err = r.Context.Db.Where("id =?", cateId).Find(cate).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return
	}
	return
}

func (r *ConsolePostService) PostUpdate(postId int, ps pojos.PostStore) (err error) {
	postUpdate := &models.ZPosts{
		Title:    ps.Title,
		UserId:   1,
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	postUpdate.Content = string(html)
	session := r.Context.Db.Begin()
	defer session.Close()
	err = session.Table((&models.ZPosts{}).TableName()).Where("id = ?", postId).Update(postUpdate).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     err.Error(),
		})
		_ = session.Rollback()
		return
	}
	postCate := new(models.ZPostCate)
	err = session.Where("post_id = ?", postId).Delete(postCate).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
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

		err = session.Create(postCateCreate).Error
		if err != nil {
			r.Context.Log.Error(map[string]string{
				"message": "service.PostUpdate",
				"err":     err.Error(),
			})
			_ = session.Rollback()
			return
		}

		if postCateCreate.Id < 1 {
			r.Context.Log.Error(map[string]string{
				"message": "service.PostUpdate",
				"err":     "post cate update no succeed",
			})
			_ = session.Rollback()
			return
		}
	}

	postTag := make([]models.ZPostTag, 0)
	err = session.Where("post_id = ?", postId).
		Find(&postTag).
		Error

	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     "get post tag  no succeed",
		})
		_ = session.Rollback()
		return
	}

	if len(postTag) > 0 {
		for _, v := range postTag {
			err = session.Where("id=?", v.TagId).
				Update("num", gorm.Expr("num  + ?", 1)).
				Error
			if err != nil {
				r.Context.Log.Error(map[string]string{
					"message": "service.PostUpdate post tag decr  err",
					"err":     err.Error(),
				})
				_ = session.Rollback()
				return
			}

		}

		err = session.Where("post_id = ?", postId).Delete(new(models.ZPostTag)).Error

		if err != nil {
			r.Context.Log.Error(map[string]string{
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
			err = session.Create(postTagCreate).Error
			if err != nil {
				r.Context.Log.Error(map[string]string{
					"message": "service.PostUpdate post tag insert err",
					"err":     err.Error(),
				})
				session.Rollback()
				return
			}

			err = session.Where("id=?", v).
				Update("num", gorm.Expr("price  + ?", 1)).
				Error
			if err != nil {
				r.Context.Log.Error(map[string]string{
					"message": "service.PostStore post tag incr err",
					"err":     err.Error(),
				})
				session.Rollback()
				return
			}
		}
	}
	session.Commit()
	return
}

func (r *ConsolePostService) PostDestroy(postId int) (bool, error) {
	post := new(models.ZPosts)
	toBeCharge := time.Now().Format(time.RFC3339)
	timeLayout := time.RFC3339
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, toBeCharge, loc)
	post.DeletedAt = &theTime
	err = r.Context.Db.Where("id =?", postId).Update(post).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostDestroy",
			"err":     err.Error(),
		})
		return false, err
	}
	return true, nil
}

func (r *ConsolePostService) PostUnTrash(postId int) (res bool, err error) {
	post := new(models.ZPosts)
	theTime, _ := time.Parse("2006-01-02 15:04:05", "")
	post.DeletedAt = &theTime
	err = r.Context.Db.Table((&models.ZPosts{}).TableName()).Where("id =?", postId).Update(post).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostUnTrash",
			"err":     err.Error(),
		})
		return false, err
	}
	return true, nil
}

func (r *ConsolePostService) PostCnt() (cnt int64, err error) {
	err = r.getDbaTable().Table((&models.ZPosts{}).TableName()).Count(&cnt).Error
	return
}

func (r *ConsolePostService) PostTagListCount(tagId int, limit int, offset int) (count int64, err error) {
	err = r.Context.Db.Table((&models.ZPostTag{}).TableName()).
		Where("tag_id = ?", tagId).
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Count(&count).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostTagListCount",
			"err":     err.Error(),
		})
		return 0, err
	}
	return
}

func (r *ConsolePostService) PostTagList(tagId int, limit int, offset int) (postListArr []*pojos.ConsolePostList, err error) {
	postListArr = make([]*pojos.ConsolePostList, 0)
	rows, err := r.Context.Db.Table((&models.ZPostTag{}).TableName()).
		Where("tag_id = ?", tagId).Order("id desc").
		Limit(limit).
		Offset(offset).Rows()

	if err != nil {
		r.Context.Log.Error(map[string]string{
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
			r.Context.Log.Error(map[string]string{
				"message": "service.Index.PostTagList",
				"err":     err.Error(),
			})
			return nil, err
		}

		post := new(models.ZPosts)
		err = r.Context.Db.Table((&models.ZPosts{}).TableName()).
			Where("id=?", postTag.PostId).
			Find(post).Error

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
	err = r.Context.Db.Table((&models.ZPostCate{}).TableName()).
		Where("cate_id = ?", cateId).Order("id desc").Limit(limit).Offset(offset).Count(&count).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostCateListCount",
			"err":     err.Error(),
		})
		return 0, err
	}
	return
}

func (r *ConsolePostService) PostCateList(cateId int, limit int, offset int) (postListArr []*pojos.ConsolePostList, err error) {
	rows, err := r.Context.Db.Where("cate_id = ?", cateId).
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Rows()

	if err != nil {
		r.Context.Log.Error(map[string]string{
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
			r.Context.Log.Error(map[string]string{
				"message": "service.Index.PostCateList",
				"err":     err.Error(),
			})
			return nil, err
		}

		post := new(models.ZPosts)
		err = r.Context.Db.Table((&models.ZPosts{}).TableName()).
			Where("id =?", postCate.PostId).
			Find(post).
			Error

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
