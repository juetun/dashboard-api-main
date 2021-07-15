// Package srv_impl
/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-12
 * Time: 21:03
 */
package srv_impl

import (
	"errors"
	"html/template"
	"strconv"
	"time"

	"gorm.io/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type ConsolePostService struct {
	base.ServiceBase
}

func NewConsolePostService(context ...*base.Context) (p *ConsolePostService) {
	p = &ConsolePostService{}
	p.SetContext(context...)
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
		r.Context.Error(map[string]interface{}{
			"message": "service.ConsolePostCount",
			"limit":   limit,
			"offset":  offset,
			"isTrash": isTrash,
			"err":     err.Error(),
		})
	}
	return
}

func (r *ConsolePostService) getDbaTable() *gorm.DB {
	return r.Context.Db.Table((&models.ZPosts{}).TableName())
}
func (r *ConsolePostService) ConsolePostIndex(dba *gorm.DB, limit, offset int, isTrash bool) (postListArr *[]wrappers.ConsolePostList, err error) {
	postListArr = &[]wrappers.ConsolePostList{}
	if dba == nil {
		dba = r.getDbaTable().Unscoped().Where("deleted_at NOT NULL")
	}

	var dt []models.ZPosts
	err = dba.Order("id desc").
		Limit(limit).
		Offset(offset).
		Find(&dt).Error

	if err != nil {

		r.Context.Error(map[string]interface{}{
			"message": "service.ConsolePostIndex",
			"limit":   limit,
			"offset":  offset,
			"isTrash": isTrash,
			"err":     err.Error(),
		})
		return
	}
	srv := NewCategoryService(r.Context)
	srvTag := NewTagService(r.Context)
	srvUser := NewUserService(r.Context)
	ids, userId := r.uniquePostId(dt)
	var mapCates *map[string]wrappers.PostShow
	mapCates, err = srv.GetPostCateByPostIds(ids)
	if err != nil {
		return
	}
	var mapTags *map[int][]wrappers.ConsoleTag
	mapTags, err = srvTag.GetPostTagsByPostIds(ids)
	if err != nil {
		return
	}
	var mapUser *map[string]models.UserMain
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
		postList := wrappers.ConsolePostList{
			Post: wrappers.ConsolePost{
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
			Category: wrappers.ConsoleCate{},
			Tags:     make([]wrappers.ConsoleTag, 0),
			View:     wrappers.ConsoleView{},
			Author:   wrappers.ConsoleUser{},
		}
		pid := strconv.Itoa(post.Id)
		if _, ok := (*mapCates)[pid]; ok {
			postList.Category = wrappers.ConsoleCate{
				Id:          (*mapCates)[pid].ZCategories.Id,
				Name:        (*mapCates)[pid].ZCategories.Name,
				DisplayName: (*mapCates)[pid].ZCategories.DisplayName,
				SeoDesc:     (*mapCates)[pid].ZCategories.SeoDesc,
			}
		}
		if _, ok := (*mapUser)[post.UserHId]; ok {
			postList.Author = wrappers.ConsoleUser{
				UserHid: (*mapUser)[post.UserHId].UserHid,
				Name:    (*mapUser)[post.UserHId].Name,
				Email:   (*mapUser)[post.UserHId].Email,
				Status:  (*mapUser)[post.UserHId].Status,
			}
		}
		if _, ok := (*mapView)[pid]; ok {
			postList.View = wrappers.ConsoleView{Num: (*mapView)[pid].Num}
		}
		if _, ok := (*mapTags)[post.Id]; ok {
			postList.Tags = (*mapTags)[post.Id]
		}

		*postListArr = append(*postListArr, postList)
	}
	return
}
func (r *ConsolePostService) uniquePostId(dt []models.ZPosts) (ids []string, userId []string) {
	ids = make([]string, 0, len(dt))
	userId = make([]string, 0, len(dt))
	mUid := make(map[string]string)
	mId := make(map[string]string)
	for _, post := range dt {
		if _, ok := mUid[post.UserHId]; !ok {
			userId = append(userId, post.UserHId)
			mUid[post.UserHId] = post.UserHId
		}
		pid := strconv.Itoa(post.Id)
		if _, ok := mId[pid]; !ok {
			ids = append(ids, pid)
			mId[pid] = pid
		}
	}
	return
}

func (r *ConsolePostService) getZPostViewsDbaTable() *gorm.DB {
	return r.Context.Db.Table((&models.ZPostViews{}).TableName())
}

func (r *ConsolePostService) PostView(postId []string) (postV *map[string]models.ZPostViews, err error) {
	postV = &map[string]models.ZPostViews{}
	var views []models.ZPostViews
	if len(postId) == 0 {
		return
	}
	err = r.getZPostViewsDbaTable().Where("post_id in (?)", postId).
		Find(&views).Error
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.PostView",
			"postId":  postId,
			"err":     err.Error(),
		})
	}
	for _, value := range views {
		(*postV)[value.PostId] = value
	}
	return
}

func (r *ConsolePostService) PostStore(ps wrappers.PostStore, userId string) {
	postCreate := &models.ZPosts{
		Title:    ps.Title,
		UserHId:  userId,
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

		r.Context.Error(map[string]interface{}{
			"message": "service.PostStore",
			"err":     err.Error(),
		})
		_ = session.Rollback()
		return
	}
	if postCreate.Id < 1 {
		r.Context.Info(map[string]interface{}{
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

			r.Context.Error( map[string]interface{}{
				"message": "service.PostStore", "err": err.Error(),
			})
			_ = session.Rollback()
			return
		}

		if postCateCreate.Id < 1 {
			r.Context.Error( map[string]interface{}{
				"message": "service.PostStore",
				"err": "post cate store not succeed",
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
				r.Context.Error( map[string]interface{}{
					"message": "service.PostStore post tag insert err", "err": err.Error(),
				})
				_ = session.Rollback()
				return
			}
			if postTagCreate.Id < 1 {
				r.Context.Error( map[string]interface{}{
					"message": "service.PostStore", "err": "post tag store not succeed",
				})
				_ = session.Rollback()
				return
			}

			err = session.Where("id=?", v).
				Update("num", gorm.Expr("price + ?", 1)).
				Error
			if err != nil {
				r.Context.Error( map[string]interface{}{
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
		r.Context.Error( map[string]interface{}{
			"message": "service.PostStore", "err": err.Error(),
		})
		_ = session.Rollback()
		return
	}

	if postView.Id < 1 {
		r.Context.Error( map[string]interface{}{
			"message": "service.PostStore", "err": "post view store no succeed",
		})
		_ = session.Rollback()
		return
	}

	_ = session.Commit()

	uid, err := common.ZHashId.Encode([]int{postCreate.Id})
	if err != nil {
		r.Context.Error( map[string]interface{}{
			"message": "service.PostStore create uid error", "err": err.Error(),
		})
		return
	}

	newPostCreate := models.ZPosts{
		Uid: uid,
	}
	err = session.Where("id = ?", postCreate.Id).Updates(newPostCreate).Error
	if err != nil {
		r.Context.Error( map[string]interface{}{
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
		r.Context.Error( map[string]interface{}{
			"message": "service.PostDetail",
			"err":     err.Error(),
		})
		return post, err
	}
	return post, nil
}

func (r *ConsolePostService) IndexPostDetailDao(postId int) (postDetail wrappers.IndexPostDetail, err error) {
	post := new(models.ZPosts)
	err = r.getDbaTable().
		Where("id = ?", postId).
		Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
		Find(post).
		Error
	if err != nil {
		r.Context.Error( map[string]interface{}{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	if post.Id <= 0 {
		return postDetail, errors.New("Post do not exists ")
	}
	Post := wrappers.IndexPost{
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
		r.Context.Error( map[string]interface{}{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	var Tags []wrappers.ConsoleTag
	for _, v := range tags {
		consoleTag := wrappers.ConsoleTag{
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
		r.Context.Error( map[string]interface{}{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	Cate := wrappers.ConsoleCate{
		Id:          cate.Id,
		Name:        cate.Name,
		DisplayName: cate.DisplayName,
		SeoDesc:     cate.SeoDesc,
	}

	// view
	pid := strconv.Itoa(post.Id)
	view, err := r.PostView([]string{pid})
	if err != nil {
		r.Context.Error( map[string]interface{}{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	View := wrappers.ConsoleView{}
	if _, ok := (*view)[pid]; ok {
		View.Num = (*view)[pid].Num
	}
	srvUser := NewUserService()
	// user
	user, err := srvUser.GetUserById(post.UserHId)
	if err != nil {
		r.Context.Error( map[string]interface{}{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}
	Author := wrappers.ConsoleUser{
		UserHid: user.UserHid,
		Name:    user.Name,
		Email:   user.Email,
		Status:  user.Status,
	}

	// last post
	lastPost, err := r.LastPost(postId)
	if err != nil {
		r.Context.Error( map[string]interface{}{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}

	// next post
	nextPost, err := r.NextPost(postId)
	if err != nil {
		r.Context.Error( map[string]interface{}{
			"message": "service.IndexPostDetailDao",
			"err":     err.Error(),
		})
		return
	}

	postDetail = wrappers.IndexPostDetail{
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
		r.Context.Error( map[string]interface{}{
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
		r.Context.Error( map[string]interface{}{
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
		r.Context.Error( map[string]interface{}{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return
	}
	cate = new(models.ZCategories)
	err = r.Context.Db.Where("id =?", cateId).Find(cate).Error
	if err != nil {
		r.Context.Error( map[string]interface{}{
			"message": "service.PostCates",
			"err":     err.Error(),
		})
		return
	}
	return
}

func (r *ConsolePostService) PostUpdate(postId int, ps wrappers.PostStore) (err error) {
	session := r.Context.Db.Begin()
	defer session.Commit()

	postUpdate := &models.ZPosts{
		Title:    ps.Title,
		UserHId:  "",
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	postUpdate.Content = string(html)

	err = session.Table((&models.ZPosts{}).TableName()).
		Where("id = ?", postId).
		Updates(postUpdate).
		Error
	if err != nil {
		r.Context.Error( map[string]interface{}{
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
		r.Context.Error( map[string]interface{}{
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

func (r *ConsolePostService) postTagLogic(postId int, ps *wrappers.PostStore, session *gorm.DB) (err error) {
	dao := dao_impl.NewDaoPostTag(r.Context)
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

	var countList *[]wrappers.TagCount
	countList, err = dao.GetEveryTagCountByTagIds(postId, &tagId)
	if err != nil {
		return
	}
	daoTag := dao_impl.NewDaoTag(r.Context)
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
func (r *ConsolePostService) postCategoryLogic(postId int, ps *wrappers.PostStore, session *gorm.DB) (err error) {
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
		r.Context.Error( map[string]interface{}{
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
		r.Context.Error( map[string]interface{}{
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
		r.Context.Error( map[string]interface{}{
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
		r.Context.Error( map[string]interface{}{
			"message": "service.PostTagListCount",
			"err":     err.Error(),
		})
		return 0, err
	}
	return
}

func (r *ConsolePostService) PostTagList(tagId int, limit int, offset int) (postListArr *[]wrappers.ConsolePostList, err error) {
	postListArr = &[]wrappers.ConsolePostList{}
	rows, err := r.Context.Db.Table((&models.ZPostTag{}).TableName()).
		Where("tag_id = ?", tagId).Order("id desc").
		Limit(limit).
		Offset(offset).Rows()

	if err != nil {
		r.Context.Error( map[string]interface{}{
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
			r.Context.Error( map[string]interface{}{
				"message": "service.Index.PostTagList",
				"err":     err.Error(),
			})
			return nil, err
		}

		post := new(models.ZPosts)
		err = r.Context.Db.Table((&models.ZPosts{}).TableName()).
			Where("id=?", postTag.PostId).
			Find(post).Error

		consolePost := wrappers.ConsolePost{
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

		postList := wrappers.ConsolePostList{
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
		r.Context.Error( map[string]interface{}{
			"message": "service.PostCateListCount",
			"err":     err.Error(),
		})
		return 0, err
	}
	return
}

func (r *ConsolePostService) PostCateList(cateId int, limit int, offset int) (postListArr *[]wrappers.ConsolePostList, err error) {
	postListArr = &[]wrappers.ConsolePostList{}
	rows, err := r.Context.Db.Where("cate_id = ?", cateId).
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Rows()

	if err != nil {
		r.Context.Error( map[string]interface{}{
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
			r.Context.Error( map[string]interface{}{
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

		consolePost := wrappers.ConsolePost{
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

		postList := wrappers.ConsolePostList{
			Post: consolePost,
		}
		*postListArr = append(*postListArr, postList)
	}

	return
}
