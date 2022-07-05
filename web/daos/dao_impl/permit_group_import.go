package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroupImportImpl struct {
	base.ServiceDao
}

func (r *DaoPermitGroupImportImpl) GetMenuIdsByPermitByGroupIds(menuIds, groupIds []int64, ) (res []*models.AdminUserGroupImport, err error) {
	res = []*models.AdminUserGroupImport{}
	if len(menuIds) == 0 || len(groupIds) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"menuIds":  menuIds,
			"groupIds": groupIds,
			"err":      err.Error(),
		}, "DaoPermitGroupImportImplGetMenuIdsByPermitByGroupIds")
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminUserGroupImport
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("menu_id IN(?) AND group_id IN(?)", menuIds,
				groupIds).
			Find(&res).Error
		return
	})

	return
}

func (r *DaoPermitGroupImportImpl) DeleteGroupMenuPermitByGroupIdAndMenuIds(groupId int64, menuIds ...int64) (err error) {

	if len(menuIds) == 0 || groupId == 0 {
		return
	}
	var m1 models.AdminUserGroupImport

	logContent := map[string]interface{}{
		"groupId": groupId,
		"menuIds": menuIds,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitGroupImportDeleteGroupMenuPermitByGroupIdAndMenuIds")
			err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		}
	}()

	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m1)

		actErrorHandlerResult.Err = actErrorHandlerResult.Db.
			Table(actErrorHandlerResult.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("group_id =? AND  menu_id IN(?) ", groupId, menuIds).
			Delete(&m1).Error
		return
	})

	return
}

func (r *DaoPermitGroupImportImpl) DeleteGroupImportWithGroupIdAndImportIds(groupId int64, importIds ...int64) (err error) {
	if len(importIds) == 0 && groupId == 0 {
		return
	}
	var m1 models.AdminUserGroupImport
	db := r.Context.Db.Table(m1.TableName()).Scopes(base.ScopesDeletedAt())
	if groupId > 0 {
		db = db.Where("group_id =?", groupId)
	}
	if len(importIds) > 0 {
		db = db.Where(" import_id IN(?)", importIds)
	}
	if err = db.Delete(&models.AdminImport{}).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupId":   groupId,
			"importIds": importIds,
			"err":       err.Error(),
		}, "DaoPermitGroupImportImplDeleteGroupImportWithMenuId")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitGroupImportImpl) DeleteGroupImportWithMenuId(menuId ...int64) (err error) {
	var m1 models.AdminUserGroupImport
	if err = r.Context.Db.Table(m1.TableName()).Scopes(base.ScopesDeletedAt()).
		Where("menu_id IN(?) ", menuId).
		Delete(&models.AdminImport{}).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err.Error(),
		}, "DaoPermitGroupImportImplDeleteGroupImportWithMenuId")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitGroupImportImpl) BatchAddDataGroupImport(list []base.ModelBase) (err error) {
	if len(list) == 0 {
		return
	}
	logContent := map[string]interface{}{
		"list": list,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitGroupImportImplBatchAddData")
			err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		}
	}()
	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		var m models.AdminUserGroupImport
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)
		actErrorHandlerResult.Err = r.BatchAdd(
			actErrorHandlerResult.ParseBatchAddDataParameter(
				base.BatchAddDataParameterData(list),
			),
		)
		return
	})

	return
}
func (r *DaoPermitGroupImportImpl) GetSelectImportByImportId(groupId int64, importId ...int64) (res []models.AdminUserGroupImport, err error) {
	if len(importId) == 0 {
		return
	}
	var m models.AdminUserGroupImport
	if err = r.Context.Db.Table(m.TableName()).
		Where("import_id IN(?)  AND group_id=?", importId, groupId).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"importId": importId,
			"err":      err,
		}, "daoPermitGetSelectImportByImportId")
		return
	}
	return
}

func NewDaoPermitGroupImport(c ...*base.Context) (res daos.DaoPermitGroupImport) {
	p := &DaoPermitGroupImportImpl{}
	p.SetContext(c...)
	res = p
	return
}
