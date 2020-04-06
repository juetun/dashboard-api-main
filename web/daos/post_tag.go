/**
* @Author:changjiang
* @Description:
* @File:tag
* @Version: 1.0.0
* @Date 2020/4/5 8:21 下午
 */
package daos

import (
	"errors"
	"strings"

	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/base-wrapper/lib/base"
)

type DaoPostTag struct {
 	base.ServiceDao
}

func NewDaoPostTag(context ...*base.Context) (p *DaoPostTag) {
	p = &DaoPostTag{}
	p.SetContext(context)
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

func (r *DaoPostTag) GetEveryTagCountByTagIds(postId int, tagIds *[]int) (list *[]pojos.TagCount, err error) {
	list = &[]pojos.TagCount{}
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
		r.Context.Log.Error(map[string]string{
			"message": "service.PostUpdate",
			"err":     "get post tag  no succeed",
		})
		return
	}
	return
}
