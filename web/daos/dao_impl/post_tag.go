// Package dao_impl
/**
* @Author:changjiang
* @Description:
* @File:tag
* @Version: 1.0.0
* @Date 2020/4/5 8:21 下午
 */
package dao_impl

import (
	"context"
	"errors"
	"strings"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type DaoPostTag struct {
	base.ServiceDao
}

func NewDaoPostTag(c ...*base.Context) (p *DaoPostTag) {
	p = &DaoPostTag{}
	p.SetContext(c...)
	s, ctx := p.Context.GetTraceId()
	p.Context.Db, p.Context.DbName, _ = base.GetDbClient(&base.GetDbClientData{
		Context:     p.Context,
		DbNameSpace: daos.DatabaseDefault,
		CallBack: func(db *gorm.DB, dbName string) (dba *gorm.DB, err error) {
			dba = db.WithContext(context.WithValue(ctx, app_obj.DbContextValueKey, base.DbContextValue{
				TraceId: s,
				DbName:  dbName,
			}))
			return
		},
	})
	return
}

func (r *DaoPostTag) InsertPostTag(data *[]map[string]interface{}) (err error) {
	field := []string{"style_id", "desc"}
	dataMsg := make([]string, 0)
	dataTmp := make([]interface{}, 0)
	duplicate := make([]string, 0)
	dataMsgA := make([]string, 0)
	if len(*data) == 0 {
		return errors.New("插入数据异常")
	}

	field = []string{
		"post_id",
		"tag_id",
		"created_at",
		"updated_at",
		"deleted_at",
	}
	for key, value := range *data {
		for _, v := range field {
			if key == 0 {
				if v != "created_at" {
					duplicate = append(duplicate, "`"+v+"`=VALUES(`"+v+"`)")
				}
				dataMsgA = append(dataMsgA, "?")
			}
			var v1 interface{}
			if _, ok := value[v]; ok {
				v1 = value[v]
			}
			dataTmp = append(dataTmp, v1)
		}
		dataMsg = append(dataMsg, "("+strings.Join(dataMsgA, ",")+")")
	}
	sql := "INSERT INTO `" + ((&models.ZPostTag{}).TableName()) + "`(`" + strings.Join(field, "`,`") +
		"`) VALUES" + strings.Join(dataMsg, ",") + " ON DUPLICATE KEY UPDATE " + strings.Join(duplicate, ",")
	err = r.Context.Db.Exec(sql, dataTmp...).Error
	return
}

func (r *DaoPostTag) DeleteDataByPostId(postId int) (err error) {
	err = r.Context.Db.Where("post_id = ?", postId).Unscoped().
		Delete(&models.ZPostTag{}).
		Error
	return
}

func (r *DaoPostTag) GetEveryTagCountByTagIds(postId int, tagIds *[]int) (list *[]wrappers.TagCount, err error) {
	list = &[]wrappers.TagCount{}
	err = r.Context.Db.Table((&models.ZPostTag{}).TableName()).
		Select("post_id,tag_id,count(*) as count").
		Where("post_id =? AND tag_id in (?) AND deleted_at IS NULL", postId, *tagIds).
		Unscoped().
		Group("post_id,tag_id").Find(list).
		Error
	return
}
func (r *DaoPostTag) GetListByPostId(postId int) (postTag *[]models.ZPostTag, err error) {
	postTag = &[]models.ZPostTag{}
	err = r.Context.Db.Where("post_id = ?", postId).
		Find(postTag).
		Error

	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "service.PostUpdate",
			"postId":  postId,
			"err":     "get post tag  no succeed",
		})
		return
	}
	return
}
