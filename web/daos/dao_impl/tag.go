/**
* @Author:changjiang
* @Description:
* @File:tag
* @Version: 1.0.0
* @Date 2020/4/5 8:21 下午
 */
package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type DaoTag struct {
	base.ServiceDao
}

func NewDaoTag(context ...*base.Context) (p *DaoTag) {
	p = &DaoTag{}
	p.SetContext(context...)
	return
}
func (r *DaoTag) UpdateTagNumById(tagCount *wrappers.TagCount) (err error) {
	err = r.Context.Db.Table((&models.ZTags{}).TableName()).
		Where("id=?", tagCount.TagId).
		Update("num", tagCount.Count).
		Error
	return
}
