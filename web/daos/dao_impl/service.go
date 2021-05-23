/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/23 10:34 上午
 */
package dao_impl

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type DaoServiceImpl struct {
	base.ServiceDao
}

func (r *DaoServiceImpl) Update(condition, data map[string]interface{}) (err error) {
	var m models.AdminApp
	if err = r.Context.Db.Model(&m).Where(condition).Update(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "DaoServiceImplUpdate")
		return
	}
	return
}

func (r *DaoServiceImpl) Create(app *models.AdminApp) (err error) {
	var m models.AdminApp
	if err = r.Context.Db.Model(&m).Create(app).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"app": app,
			"err": err.Error(),
		}, "DaoServiceImplCreate")
		return
	}
	return
}

func (r *DaoServiceImpl) GetByIds(id ...int) (res []models.AdminApp, err error) {
	res = []models.AdminApp{}
	if len(id) == 0 {
		return
	}
	var m models.AdminApp
	if err = r.Context.Db.Model(&m).Where("id IN(?)", id).Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err.Error(),
		}, "DaoServiceImplGetByIds")
		return
	}
	return
}

func (r *DaoServiceImpl) fetchGetDb(db *gorm.DB, arg *wrappers.ArgServiceList) (dba *gorm.DB) {
	var m models.AdminApp
	if db == nil {
		db = r.Context.Db
	}
	dba = db.Model(&m)
	if arg.Name != "" {
		dba = dba.Where("name LIKE ?", fmt.Sprintf("%%%s%%", arg.Name))
	}
	if arg.Id > 0 {
		dba = dba.Where("id = ?", arg.Id)
	}
	if arg.UniqueKey != "" {
		dba = dba.Where("unique_key LIKE ?", fmt.Sprintf("%%%s%%", arg.UniqueKey))
	}
	if arg.Port > 0 {
		dba = dba.Where("port = ?", arg.Port)
	}
	if arg.IsStop > 0 {
		dba = dba.Where("is_stop = ?", arg.IsStop)
	}
	if arg.Desc != "" {
		dba = dba.Where("desc LIKE ?", fmt.Sprintf("%%%s%%", arg.Desc))
	}
	dba = dba.Unscoped().Where("deleted_at IS NULL")
	return
}
func (r *DaoServiceImpl) GetCount(db *gorm.DB, arg *wrappers.ArgServiceList) (total int, dba *gorm.DB, err error) {
	dba = r.fetchGetDb(db, arg)
	if err = dba.Count(&total).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "DaoServiceImplGetCount")
		return
	}
	return
}

func (r *DaoServiceImpl) GetList(db *gorm.DB, arg *wrappers.ArgServiceList, page *response.Pager) (list []models.AdminApp, err error) {
	if db == nil {
		db = r.fetchGetDb(db, arg)
	}
	if page.Order != "" {
		db = db.Order(page.Order)
	}
	if err = db.Offset(page.GetOffset()).Limit(page.PageSize).Find(&list).Error; err != nil {
		return
	}
	return
}

func NewDaoServiceImpl(context ...*base.Context) daos.DaoService {
	p := &DaoServiceImpl{}
	p.SetContext(context...)
	return p
}
