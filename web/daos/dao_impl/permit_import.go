// Package dao_impl
/**
* @Author:ChangJiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/6/5 11:35 下午
 */
package dao_impl

import (
	"fmt"
	"strings"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type DaoPermitImportImpl struct {
	base.ServiceDao
}

func (r *DaoPermitImportImpl) GetMenuImportByMenuIdAndImportIds(menuId int64, importIds ...int64) (res []*models.AdminMenuImport, err error) {
	if menuId == 0 || len(importIds) == 0 {
		return
	}

	if err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		var m models.AdminMenuImport
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)
		actErrorHandlerResult.Err = actErrorHandlerResult.
			Db.
			Table(actErrorHandlerResult.TableName).
			Where("menu_id = ? AND import_id IN(?)", menuId, importIds).
			Find(&res).
			Error
		return
	}); err != nil {
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitImportImpl) GetImportFromDbByIds(ids ...int64) (res map[int64]*models.AdminImport, err error) {
	l := len(ids)
	res = make(map[int64]*models.AdminImport, l)
	if l == 0 {
		return
	}
	var m models.AdminImport
	var data []*models.AdminImport
	logContent := map[string]interface{}{
		"ids": ids,
	}
	defer func() {
		if err != nil {
			r.Context.Error(logContent, "DaoPermitImportImplGetImportFromDbByIds")
		}
	}()

	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)

		actErrorHandlerResult.Err = r.Context.Db.Table(m.TableName()).
			Where("id IN (?)", ids).
			Scopes(base.ScopesDeletedAt()).
			Find(&data).Error
		return
	})
	if err != nil {
		logContent["err"] = err.Error()
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	for _, item := range data {
		res[item.Id] = item
	}
	return
}
func (r *DaoPermitImportImpl) GetImportForGateway(arg *wrapper_intranet.ArgGetImportPermit) (list []models.AdminImport, err error) {
	condition := map[string]interface{}{}
	if arg.AppName != "" {
		condition["app_name"] = arg.AppName
	}
	var m models.AdminImport
	db := r.Context.Db.Table(m.TableName()).
		Where(condition).Scopes(base.ScopesDeletedAt()).
		Where("need_login = ? OR need_sign = ?", models.NeedLoginNot, models.NeedSignNot)
	var e error
	if e = db.
		Find(&list).
		Error; e == nil {
		return
	}

	logContent := map[string]interface{}{
		"err": e.Error(),
	}

	defer func() {
		r.Context.Error(logContent, "DaoPermitImportGetImportForGateway")
	}()

	if err = r.CreateTableWithError(r.Context.Db, m.TableName(), e, &m, base.TableSetOption{"COMMENT": m.GetTableComment()}); err != nil {
		logContent["err1"] = err.Error()
		logContent["desc"] = "CreateTableWithError"
		return
	}

	if list, err = r.GetImportForGateway(arg); err != nil {
		logContent["err2"] = err.Error()
		logContent["desc"] = "ReGetImportForGateway"
		return
	}
	return
}

func (r *DaoPermitImportImpl) AddData(adminImport *models.AdminImport) (err error) {
	logContent := map[string]interface{}{
		"adminImport": adminImport,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitImportAddData")
		}
	}()
	if err = adminImport.InitRegexpString(); err != nil {
		logContent["desc"] = "InitRegexpString"
		return
	}
	if err = r.Context.Db.Create(adminImport).Error; err != nil {
		logContent["desc"] = "Create"
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitImportImpl) getColumnName(s string) (res string) {
	li := strings.Split(s, ";")
	res = s
	for _, s2 := range li {
		if s2 == "" {
			return
		}
		li1 := strings.Split(s2, ":")
		if len(li1) > 1 && li1[0] == "column" {
			res = li1[1]
		}
	}
	if res == "primary_key" {
		res = "id"
	}
	return
}
func (r *DaoPermitImportImpl) BatchMenuImport(tableName string, list []*models.AdminMenuImport) (err error) {
	if len(list) == 0 {
		return
	}
	var dtList []base.ModelBase
	for _, menuImport := range list {
		dtList = append(dtList, menuImport)
	}
	var batchData = base.BatchAddDataParameter{
		CommonDb: base.CommonDb{
			DbName:    r.Context.DbName,
			Db:        r.Context.Db,
			TableName: tableName,
		},
		Data: dtList,
	}
	if err = r.BatchAdd(&batchData); err != nil {
		r.Context.Error(map[string]interface{}{
			"data": list,
			"err":  err,
		}, "DaoPermitImportBatchMenuImport")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}

	return
}

func (r *DaoPermitImportImpl) DeleteByCondition(condition interface{}) (res bool, err error) {
	var m models.AdminImport
	if condition == nil {
		err = fmt.Errorf("您没有选择要删除的数据")
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitImportDeleteByCondition0")
		return
	}
	if err = r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).Where(condition).
		Delete(&models.AdminImport{}).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitImportDeleteByCondition1")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitImportImpl) UpdateByCondition(condition interface{}, data map[string]interface{}) (res bool, err error) {
	var m models.AdminImport
	if condition == nil || len(data) == 0 {
		err = fmt.Errorf("您没有选择要更新的数据")
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "DaoPermitImportUpdateByCondition0")
		return
	}
	if err = r.Context.Db.Table(m.TableName()).
		Where(condition).Scopes(base.ScopesDeletedAt()).
		Updates(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err.Error(),
		}, "DaoPermitImportUpdateByCondition1")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	res = true
	return
}

func (r *DaoPermitImportImpl) GetChildImportByMenuId(menuIds ...int64) (list []models.AdminMenuImport, err error) {

	list = []models.AdminMenuImport{}
	if len(menuIds) == 0 {
		return
	}
	var m models.AdminMenuImport
	if err = r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
		Where("menu_id IN (?)", menuIds).
		Find(&list).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"iIds": menuIds,
			"err":  err.Error(),
		}, "DaoPermitImportGetChildImportByMenuId")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}
func (r *DaoPermitImportImpl) GetImportMenuByImportIds(iIds ...int64) (list []models.AdminMenuImport, err error) {
	list = []models.AdminMenuImport{}
	if len(iIds) == 0 {
		return
	}
	var m models.AdminMenuImport
	if err = r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
		Where("import_id IN (?)", iIds).
		Find(&list).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"iIds": iIds,
			"err":  err.Error(),
		}, "DaoPermitImportGetImportMenuByImportIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}
func (r *DaoPermitImportImpl) UpdateMenuImport(condition string, data map[string]interface{}) (err error) {

	var m models.AdminMenuImport
	var e error
	var fetchData = base.FetchDataParameter{
		SourceDb:  r.Context.Db,
		DbName:    r.Context.DbName,
		TableName: m.TableName(),
	}
	if e = fetchData.SourceDb.Table(fetchData.TableName).Scopes(base.ScopesDeletedAt()).
		Where(condition).
		Updates(data).
		Error; e == nil {
		return
	}

	logContent := map[string]interface{}{
		"condition": condition,
		"e":         e.Error(),
	}
	defer func() {
		r.Context.Error(logContent, "DaoPermitImportUpdateMenuImport")
	}()
	if err = r.CreateTableWithError(fetchData.SourceDb, fetchData.TableName, e, &m, base.TableSetOption{"COMMENT": m.GetTableComment()}); err != nil {
		return
	}
	if err = r.UpdateMenuImport(condition, data); err != nil {
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}

	return
}

func (r *DaoPermitImportImpl) BatchAddData(list []models.AdminImport) (err error) {

	var (
		m        models.AdminImport
		data     *base.BatchAddDataParameter
		dataList []base.ModelBase
		e        error
	)

	for _, adminImport := range list {
		dataList = append(dataList, &adminImport)
	}
	if data, err = r.GetDefaultBatchAddDataParameter(dataList...); err != nil {
		return
	}

	if e = r.BatchAdd(data); e == nil {
		return
	}
	logContent := map[string]interface{}{
		"list": list,
		"e":    e.Error(),
	}
	defer func() {
		r.Context.Error(logContent, "DaoPermitImportBatchAddData")
	}()
	if err = r.CreateTableWithError(data.Db, data.TableName, e, &m, base.TableSetOption{"COMMENT": m.GetTableComment()}); err != nil {
		logContent["err1"] = err.Error()
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	if err = r.BatchAddData(list); err != nil {
		logContent["err2"] = err.Error()
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}
func (r *DaoPermitImportImpl) DeleteImportByIds(id ...int64) (err error) {
	if len(id) == 0 {
		return
	}
	var m models.AdminImport
	if err = r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).Unscoped().
		Where("id IN(?)", id).
		Delete(&models.AdminImport{}).Error; err == nil {
		r.Context.Error(map[string]interface{}{
			"id":  id,
			"err": err,
		}, "daoPermitDeleteImportByIds0")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	daoGroupImport := NewDaoPermitGroupImport(r.Context)
	if err = daoGroupImport.DeleteGroupImportWithMenuId(id...); err != nil {
		return
	}

	return
}

func (r *DaoPermitImportImpl) GetImportByCondition(condition map[string]interface{}) (list []models.AdminImport, err error) {
	list = []models.AdminImport{}
	if len(condition) == 0 {
		return
	}
	var m models.AdminImport

	if err = r.Context.Db.Table(m.TableName()).
		Where(condition).
		Scopes(base.ScopesDeletedAt()).
		Find(&list).
		Limit(1000).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitImportGetImportByCondition")
		return
	}
	return
}

func (r *DaoPermitImportImpl) GetDefaultOpenImportByMenuIds(menuId ...int64) (res []models.AdminImport, err error) {
	res = []models.AdminImport{}
	if len(menuId) == 0 {
		return
	}
	var m models.AdminImport
	var ami models.AdminMenuImport
	if err = r.Context.Db.Table(ami.TableName()).
		Unscoped().
		Select(`b.*`).
		Joins(fmt.Sprintf("AS a LEFT JOIN %s AS b  ON  b.`id` = a.import_id", m.TableName())).
		Where("a.menu_id IN(?) AND  b.default_open = ? AND `a`.`deleted_at` IS NULL ",
			menuId,
			models.DefaultOpen).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err,
		}, "GetDefaultOpenImportByMenuIds")
		return
	}
	return
}
func NewDaoPermitImport(c ...*base.Context) (res daos.DaoPermitImport) {
	p := &DaoPermitImportImpl{}
	p.SetContext(c...)
	res = p
	return
}
