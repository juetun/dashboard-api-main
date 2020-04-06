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
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/pojos"
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

func (r *ConsolePostService) ConsolePostCount(limit int, offset int, isTrash bool) (dba *gorm.DB, count int64, err error) {
	dba = r.getDbaTable()
	if isTrash {
		dba = dba.Unscoped().Where("deleted_at IS NOT NULL")
	} else {
		dba = dba.Unscoped().Where("deleted_at IS  NULL")
	}
	err = dba.Count(&count).
		Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.ConsolePostCount", "err": err.Error(),
		})
	}
	return
}

func (r *ConsolePostService) getDbaTable() *gorm.DB {
	return r.Context.Db.Table((&models.ZPosts{}).TableName())
}
func (r *ConsolePostService) ConsolePostIndex(dba *gorm.DB, limit, offset int, isTrash bool) (postListArr *[]pojos.ConsolePostList, err error) {
	postListArr = &[]pojos.ConsolePostList{}
	if dba == nil {
		dba = r.getDbaTable().Unscoped().Where("deleted_at NOT NULL")
	}

	var dt []models.ZPosts
	err = dba.Order("id desc").
		Limit(limit).
		Offset(offset).
		Find(&dt).Error

	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.ConsolePostIndex", "err": err.Error(),
		})
		return
	}
	srv := NewCategoryService(r.Context)
	srvTag := NewTagService(r.Context)
	srvUser := NewUserService(r.Context)
	ids, userId := r.uniquePostId(&dt)
	var mapCates *map[string]pojos.PostShow
	mapCates, err = srv.GetPostCateByPostIds(ids)
	if err != nil {
		return
	}
	var mapTags *map[int][]pojos.ConsoleTag
	mapTags, err = srvTag.GetPostTagsByPostIds(ids)
	if err != nil {
		return
	}
	var mapUser *map[int]models.ZUsers
	mapUser, err = srvUser.GetUserMapByIds(userId)
	if err != nil {
		return
	}
	var mapView *map[string]models.ZPostViews
	mapView, err = r.PostView(ids)
	if err != nil {
		return
	}

	for _, post := range dt {
		postList := pojos.ConsolePostList{
			Post: pojos.ConsolePost{
				Id:        post.Id,
				Uid:       post.Uid,
				Title:     post.Title,
				Summary:   post.Summary,
				Original:  post.Original,
				Content:   post.Content,
				Password:  post.Password,
				CreatedAt: post.CreatedAt,
				UpdatedAt: post.UpdatedAt,
			},
			Category: pojos.ConsoleCate{},
			Tags:     make([]pojos.ConsoleTag, 0),
			View:     pojos.ConsoleView{},
			Author:   pojos.ConsoleUser{},
		}
		pid := strconv.Itoa(post.Id)
		if _, ok := (*mapCates)[pid]; ok {
			postList.Category = pojos.ConsoleCate{
				Id:          (*mapCates)[pid].ZCategories.Id,
				Name:        (*mapCates)[pid].ZCategories.Name,
				DisplayName: (*mapCates)[pid].ZCategories.DisplayName,
				SeoDesc:     (*mapCates)[pid].ZCategories.SeoDesc,
			}
		}
		if _, ok := (*mapUser)[post.UserId]; ok {
			postList.Author = pojos.ConsoleUser{
				Id:     (*mapUser)[post.UserId].Id,
				Name:   (*mapUser)[post.UserId].Name,
				Email:  (*mapUser)[post.UserId].Email,
				Status: (*mapUser)[post.UserId].Status,
			}
		}
		if _, ok := (*mapView)[pid]; ok {
			postList.View = pojos.ConsoleView{Num: (*mapView)[pid].Num}
		}
		if _, ok := (*mapTags)[post.Id]; ok {
			postList.Tags = (*mapTags)[post.Id]
		}

		*postListArr = append(*postListArr, postList)
	}
	return
}
func (r *ConsolePostService) uniquePostId(dt *[]models.ZPosts) (ids *[]string, userId *[]int) {
	ids = &[]string{}
	userId = &[]int{}
	mUid := make(map[int]int)
	mId := make(map[string]string)
	for _, post := range *dt {
		if _, ok := mUid[post.UserId]; !ok {
			*userId = append(*userId, post.UserId)
			mUid[post.UserId] = post.UserId
		}
		pid := strconv.Itoa(post.Id)
		if _, ok := mId[pid]; !ok {
			*ids = append(*ids, pid)
			mId[pid] = pid
		}
	}
	return
}

func (r *ConsolePostService) getZPostViewsDbaTable() *gorm.DB {
	return r.Context.Db.Table((&models.ZPostViews{}).TableName())
}

func (r *ConsolePostService) PostView(postId *[]string) (postV *map[string]models.ZPostViews, err error) {
	postV = &map[string]models.ZPostViews{}
	var views []models.ZPostViews
	if len(*postId) == 0 {
		return
	}
	err = r.getZPostViewsDbaTable().Where("post_id in (?)", *postId).
		Find(&views).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostView", "err": err.Error(),
		})
	}
	for _, value := range views {
		(*postV)[value.PostId] = value
	}
	return
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
	defer session.Commit()
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

	if ps.Category != "0" && ps.Category != "" {
		postCateCreate := models.ZPostCate{
			PostId: strconv.Itoa(postCreate.Id),
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
		PostId: strconv.Itoa(postCreate.Id),
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
	pid := strconv.Itoa(post.Id)
	view, err := r.PostView(&[]string{pid})
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	View := pojos.ConsoleView{}
	if _, ok := (*view)[pid]; ok {
		View.Num = (*view)[pid].Num
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
	session := r.Context.Db.Begin()
	defer session.Commit()

	postUpdate := &models.ZPosts{
		Title:    ps.Title,
		UserId:   1,
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	postUpdate.Content = string(html)

	err = session.Table((&models.ZPosts{}).TableName()).
		Where("id = ?", postId).
		Update(postUpdate).
		Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     err.Error(),
		})
		_ = session.Rollback()
		return
	}

	// 删除之前的类型
	err = session.Where("post_id = ?", postId).
		Unscoped().
		Delete(&models.ZPostCate{}).
		Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     "post cate delete no succeed",
		})
		_ = session.Rollback()
		return
	}

	// 重新添加类型
	if ps.Category != "" && ps.Category != "0" {
		err = r.postCategoryLogic(postId, &ps, session)
		if err != nil {
			_ = session.Rollback()
			return
		}
	}
	err = r.postTagLogic(postId, &ps, session)
	if err != nil {
		_ = session.Rollback()
		return
	}
	return
}

func (r *ConsolePostService) postTagLogic(postId int, ps *pojos.PostStore, session *gorm.DB) (err error) {
	dao := daos.NewDaoPostTag(r.Context)
	dao.Context.Db = session
	var postTag *[]models.ZPostTag
	postTag, err = dao.GetListByPostId(postId)
	if err != nil {
		return
	}
	tagId := ps.Tags
	for _, value := range *postTag {
		tagId = append(tagId, value.TagId)
	}
	err = dao.DeleteDataByPostId(postId)
	if err != nil {
		return
	}

	// 添加当前帖子加入的话题
	insertTagRelation := make([]map[string]interface{}, 0)
	for _, value := range ps.Tags {
		insertTagRelation = append(insertTagRelation,
			map[string]interface{}{
				"post_id":    postId,
				"tag_id":     value,
				"created_at": base.TimeNormal{Time: time.Now()},
				"updated_at": base.TimeNormal{Time: time.Now()},
				"deleted_at": nil,
			},
		)
	}
	err = dao.InsertPostTag(&insertTagRelation)
	if err != nil {
		return
	}

	var countList *[]pojos.TagCount
	countList, err = dao.GetEveryTagCountByTagIds(postId, &tagId)
	if err != nil {
		return
	}
	daoTag := daos.NewDaoTag(r.Context)
	daoTag.Context.Db = session
	for _, value := range *countList {
		err = daoTag.UpdateTagNumById(&value)
		if err != nil {
			return
		}
	}
	return
}

// 帖子分类 逻辑
func (r *ConsolePostService) postCategoryLogic(postId int, ps *pojos.PostStore, session *gorm.DB) (err error) {
	postCateCreate := models.ZPostCate{
		PostId: strconv.Itoa(postId),
		CateId: ps.Category,
		Model: base.Model{
			Id:        0,
			CreatedAt: base.TimeNormal{Time: time.Now()},
			UpdatedAt: base.TimeNormal{Time: time.Now()},
			DeletedAt: nil,
		},
	}

	err = session.Create(&postCateCreate).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     err.Error(),
		})
	}
	return
}
func (r *ConsolePostService) PostDestroy(postId int) (res bool, err error) {
	post := new(models.ZPosts)

	err = r.Context.Db.Where("id =?", postId).Delete(post).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.PostDestroy",
			"err":     err.Error(),
		})
		return
	}
	res = true
	return
}

func (r *ConsolePostService) PostUnTrash(postId int) (res bool, err error) {
	err = r.Context.Db.Table((&models.ZPosts{}).TableName()).
		Where("id =?", postId).
		Update("deleted_at", nil).
		// Delete(&models.ZPosts{}).
		Error
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

func (r *ConsolePostService) PostTagList(tagId int, limit int, offset int) (postListArr *[]pojos.ConsolePostList, err error) {
	postListArr = &[]pojos.ConsolePostList{}
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
		*postListArr = append(*postListArr, postList)
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

func (r *ConsolePostService) PostCateList(cateId int, limit int, offset int) (postListArr *[]pojos.ConsolePostList, err error) {
	postListArr = &[]pojos.ConsolePostList{}
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
		return
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
			return
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
		*postListArr = append(*postListArr, postList)
	}

	return
}
