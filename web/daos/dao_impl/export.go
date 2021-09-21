// Package dao_impl
/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:49 上午
 */
package dao_impl

import (
	"context"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type DaoExportImpl struct {
	base.ServiceDao
}

func NewDaoExportImpl(c ...*base.Context) (res daos.DaoExport) {
	p := &DaoExportImpl{}
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
	return p
}
func (r DaoExportImpl) Update(model *models.ZExportData) (err error) {
	var m models.ZExportData
	if err = r.Context.Db.Table(m.TableName()).
		Where("hid = ?", model.Hid).
		Updates(model).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"model": model,
			"err":   err.Error(),
		}, "DaoExportImplUpdate")
		return
	}
	return
}
func (r DaoExportImpl) UpdateByHIds(data map[string]interface{}, hIds *[]string) (err error) {
	if len(*hIds) < 1 {
		return
	}
	var m models.ZExportData
	if err = r.Context.Db.Table(m.TableName()).
		Where("hid in(?)", *hIds).
		Updates(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"data": data,
			"hIds": hIds,
			"err":  err.Error(),
		}, "daoExportImplUpdateByHIds")
		return
	}
	return
}
func (r DaoExportImpl) Create(model *models.ZExportData) (res bool, err error) {
	res = true
	err = utils.CreateForHID(r.Context.Db, model)
	if err != nil {
		res = false
		r.Context.Error(map[string]interface{}{
			"model": model,
			"err":   err.Error(),
		}, "daoExportImplCreate")
		return
	}
	return
}
func (r DaoExportImpl) Progress(args *wrappers.ArgumentsExportProgress) (res *[]models.ZExportData, err error) {
	res = &[]models.ZExportData{}
	if len(args.IdString) < 1 {
		return
	}
	var m models.ZExportData
	if err = r.Context.Db.Table(m.TableName()).
		Where("create_user_hid=? AND hid in (?)", args.User.UserId, args.IdString).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"args": args,
			"err":  err.Error(),
		}, "daoExportImplProgress")
		return
	}
	return
}
func (r DaoExportImpl) GetListByUser(userHid string, limit int) (res *[]models.ZExportData, err error) {
	if limit == 0 {
		limit = 12
	}
	var m models.ZExportData
	res = &[]models.ZExportData{}
	if err = r.Context.Db.Table(m.TableName()).
		Where("create_user_hid=?", userHid).
		Order("created_at desc").
		Limit(limit).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"userHid": userHid,
			"limit":   limit,
			"err":     err.Error(),
		}, "daoExportImplGetListByUser")
		return
	}
	return
}
