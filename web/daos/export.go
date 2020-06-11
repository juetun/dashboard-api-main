/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:49 上午
 */
package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/basic/utils"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/pojos"
)

type DaoExport struct {
	base.ServiceDao
}

func NewDaoExport(context ...*base.Context) (p *DaoExport) {
	p = &DaoExport{}
	p.SetContext(context)
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
func (r DaoExport) Create(model *models.ZExportData) (res bool, err error) {
	res = true
	err = utils.CreateForHID(r.Context.Db, model)
	if err != nil {
		res = false
		return
	}
	return
}
func (r DaoExport) Progress(args *pojos.ArgumentsExportProgress) (res *[]models.ZExportData, err error) {
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
