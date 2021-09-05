// Package dao_impl
/**
* @Author:changjiang
* @Description:
* @File:server
* @Version: 1.0.0
* @Date 2020/6/11 10:16 上午
 */
package dao_impl

import (
	"context"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"gorm.io/gorm"
)

type DaoAppPath struct {
	base.ServiceDao
}

func NewDaoAppPath(c ...*base.Context) (p *DaoAppPath) {
	p = &DaoAppPath{}
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

func (r *DaoAppPath) GetAllApp() (res *[]models.ZAppPath, err error) {
	res = &[]models.ZAppPath{}
	var m models.ZAppPath
	err = r.Context.Db.
		Table(m.TableName()).
		Find(res).
		Error
	return
}
