/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:49 上午
 */
package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/basic/utils"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type DaoExport struct {
	base.ServiceDao
}

func NewDaoExport(context ...*base.Context) (p *DaoExport) {
	p = &DaoExport{}
	p.SetContext(context...)
	p.Context.Db = base.GetDbClient(&base.GetDbClientData{
		Context:     p.Context,
		DbNameSpace: daos.DatabaseDefault,
	})
	return
}
func (r DaoExport) Update(model *models.ZExportData) (err error) {
	var m models.ZExportData
	err = r.Context.Db.Table(m.TableName()).
		Where("hid = ?", model.Hid).
		Update(model).
		Error
	return
}
func (r DaoExport) UpdateByHIds(data map[string]interface{}, hIds *[]string) (err error) {
	if len(*hIds) < 1 {
		return
	}
	var m models.ZExportData
	err = r.Context.Db.Table(m.TableName()).
		Where("hid in(?)", *hIds).
		Update(data).
		Error
	return
}
func (r DaoExport) Create(model *models.ZExportData) (res bool, err error) {
	res = true
	err = utils.CreateForHID(r.Context.Db, model)
	if err != nil {
		res = false
		return
	}
	return
}
func (r DaoExport) Progress(args *wrappers.ArgumentsExportProgress) (res *[]models.ZExportData, err error) {
	res = &[]models.ZExportData{}
	if len(args.IdString) < 1 {
		return
	}
	var m models.ZExportData
	err = r.Context.Db.Table(m.TableName()).
		Where("create_user_hid=? AND hid in (?)", args.User.UserId, args.IdString).
		Find(&res).
		Error
	return
}
func (r DaoExport) GetListByUser(userHid string, limit int) (res *[]models.ZExportData, err error) {
	if limit == 0 {
		limit = 12
	}
	var m models.ZExportData
	res = &[]models.ZExportData{}
	err = r.Context.Db.Table(m.TableName()).
		Where("create_user_hid=?", userHid).
		Order("created_at desc").
		Limit(limit).
		Find(&res).
		Error
	return
}
