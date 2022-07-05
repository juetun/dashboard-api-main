// Package dao_impl
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

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type DaoServiceImpl struct {
	base.ServiceDao
}

func (r *DaoServiceImpl) GetByKeys(keys ...string) (res []models.AdminApp, err error) {
	res = make([]models.AdminApp, 0, len(keys))
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"keys": keys,
			"err":  err.Error(),
		}, "daoServiceImplGetByKeys")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminApp
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("unique_key IN(?)", keys).
			Find(&res).
			Error
		return
	})

	return
}

func (r *DaoServiceImpl) GetImportMenuByModule(module string) (res []wrappers.ImportMenu, err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"module": module,
			"err":    err.Error(),
		}, "daoServiceImplGetImportMenuByModule")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var ami *models.AdminMenuImport
		actRes = r.GetDefaultActErrorHandlerResult(ami)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Select("DISTINCT import_app_name as app_name").
			Where("`menu_module` = ?  AND `deleted_at` IS NULL", module).
			Find(&res).Error
		return
	})

	return
}

func (r *DaoServiceImpl) Update(condition, data map[string]interface{}) (err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "DaoServiceImplUpdate")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminApp
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where(condition).
			Updates(data).Error
		return
	})

	return
}

func (r *DaoServiceImpl) Create(app *models.AdminApp) (err error) {

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"app": app,
			"err": err.Error(),
		}, "DaoServiceImplCreate")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(app)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Create(app).Error
		return
	})

	return
}

func (r *DaoServiceImpl) GetByIds(id ...int) (res []models.AdminApp, err error) {
	res = []models.AdminApp{}
	if len(id) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err.Error(),
		}, "DaoServiceImplGetByIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminApp
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("id IN(?)", id).
			Find(&res).Error
		return
	})

	return
}

func (r *DaoServiceImpl) fetchGetDb(db *gorm.DB, arg *wrappers.ArgServiceList) (dba *gorm.DB) {
	var m models.AdminApp
	if db == nil {
		db = r.Context.Db
	}
	dba = db.Model(&m).Unscoped().Where("deleted_at IS NULL")
	if arg == nil {
		return
	}
	if len(arg.UniqueKeys) > 0 {
		dba = dba.Where("`unique_key` IN (?)", arg.UniqueKeys)
	}
	if arg.Name != "" {
		dba = dba.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", arg.Name))
	}
	if arg.Id > 0 {
		dba = dba.Where("`id` = ?", arg.Id)
	}
	if arg.UniqueKey != "" {
		dba = dba.Where("`unique_key` LIKE ?", fmt.Sprintf("%%%s%%", arg.UniqueKey))
	}
	if arg.Port > 0 {
		dba = dba.Where("port = ?", arg.Port)
	}
	if arg.IsStop > 0 {
		dba = dba.Where("`is_stop` = ?", arg.IsStop)
	}
	if arg.Desc != "" {
		dba = dba.Where("`desc` LIKE ?", fmt.Sprintf("%%%s%%", arg.Desc))
	}
	return
}
func (r *DaoServiceImpl) GetCount(db *gorm.DB, arg *wrappers.ArgServiceList) (total int64, dba *gorm.DB, err error) {
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
	if page != nil {
		if arg.PageQuery.Order != "" {
			db = db.Order(arg.PageQuery.Order)
		}
		db = db.Offset(arg.PageQuery.GetOffset()).Limit(page.PageSize)
	}
	if err = db.Find(&list).Error; err != nil {
		return
	}
	return
}

func NewDaoServiceImpl(ctxs ...*base.Context) daos.DaoService {
	p := &DaoServiceImpl{}
	p.SetContext(ctxs...)
	return p
}
