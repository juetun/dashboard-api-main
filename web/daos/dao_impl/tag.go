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

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type DaoTag struct {
	base.ServiceDao
}

func NewDaoTag(c ...*base.Context) (p *DaoTag) {
	p = &DaoTag{}
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
func (r *DaoTag) UpdateTagNumById(tagCount *wrappers.TagCount) (err error) {
	err = r.Context.Db.Table((&models.ZTags{}).TableName()).
		Where("id=?", tagCount.TagId).
		Update("num", tagCount.Count).
		Error
	return
}
