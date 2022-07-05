package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroupMenuImpl struct {
	base.ServiceDao
}

func (r *DaoPermitGroupMenuImpl) DeleteGroupImportByGroupIds(groupId ...int64) (err error) {
	if len(groupId) == 0 {
		return
	}
	var m models.AdminUserGroupImport

	logContent := map[string]interface{}{
		"groupId": groupId,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitGroupMenuImplDeleteGroupImportByGroupIds")
		}
	}()
	if err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)

		actErrorHandlerResult.Err = actErrorHandlerResult.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
			Where("group_id IN (?)", groupId).Delete(&m).Error
		return
	}); err != nil {
		return
	}
	return
}
func (r *DaoPermitGroupMenuImpl) DeleteGroupMenuByGroupIds(groupId ...int64) (err error) {
	if len(groupId) == 0 {
		return
	}
	var m models.AdminUserGroupMenu

	logContent := map[string]interface{}{
		"groupId": groupId,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "DaoPermitGroupMenuImplDeleteGroupMenuByGroupIds")
		}
	}()
	if err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)

		actErrorHandlerResult.Err = actErrorHandlerResult.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
			Where("group_id IN (?)", groupId).Delete(&m).Error
		return
	}); err != nil {
		return
	}

	return
}

func (r *DaoPermitGroupMenuImpl) DeleteGroupPermitByMenuIds(groupId int64, module string, pageMenuId []int64) (err error) {
	if len(pageMenuId) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"groupId":    groupId,
			"pageMenuId": pageMenuId,
			"module":     module,
			"err":        err,
		}, "DaoPermitGroupMenuImplDeleteGroupPermitByMenuIds")
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminUserGroupMenu
		actRes = r.GetDefaultActErrorHandlerResult(m)
		db := actRes.Db.Table(actRes.TableName).Scopes(base.ScopesDeletedAt()).
			Where("group_id = ? AND module = ?", groupId, module)
		if len(pageMenuId) > 0 {
			db = db.Where("menu_id IN (?)", pageMenuId)
		}
		actRes.Err = db.Delete(&m).Error
		return
	})
	return
}

func (r *DaoPermitGroupMenuImpl) GetMenuIdsByPermitByGroupIds(module string, groupIds ...int64) (res []models.AdminUserGroupMenu, err error) {
	if len(groupIds) == 0 {
		return
	}
	var m models.AdminUserGroupMenu
	if err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)

		if actErrorHandlerResult.Err = actErrorHandlerResult.Db.
			Table(actErrorHandlerResult.TableName).
			Select("distinct `menu_id`,`group_id`,`id`").
			Where("module = ? AND `group_id` in(?) ", module, groupIds).
			Find(&res).
			Error; actErrorHandlerResult.Err != nil {
			r.Context.Error(map[string]interface{}{
				"groupIds": groupIds,
				"err":      actErrorHandlerResult.Err.Error(),
			}, "DaoPermitGroupMenuImplGetMenuIdsByPermitByGroupIds")
			return
		}
		return
	}); err != nil {
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitGroupMenuImpl) DeleteGroupMenuByGroupIdAndIds(groupId int64, menuIds ...int64) (err error) {
	if groupId == 0 {
		return
	}

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"groupId": groupId,
			"menuIds": menuIds,
			"err":     err.Error(),
		}, "DaoPermitGroupImportImplDeleteGroupMenuByGroupIdAndIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		var m *models.AdminUserGroupMenu
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(m)
		db := actErrorHandlerResult.Db.Table(actErrorHandlerResult.TableName).Scopes(base.ScopesDeletedAt()).
			Where("group_id = ? ", groupId)
		if len(menuIds) > 0 {
			db = db.Where("menu_id IN(?)", menuIds)
		}
		actErrorHandlerResult.Err = db.
			Delete(&models.AdminImport{}).Error
		return
	})
	return
}

func (r *DaoPermitGroupMenuImpl) BatchAddDataUserGroupMenu(list []base.ModelBase) (err error) {
	if len(list) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"list": list,
			"err":  err.Error(),
		}, "DaoPermitGroupMenuImplBatchAddData")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(list[0])
		actRes.Err = r.BatchAdd(actRes.ParseBatchAddDataParameter(base.BatchAddDataParameterData(list)))
		return
	})
	return
}

func NewDaoPermitGroupMenu(c ...*base.Context) (res daos.DaoPermitGroupMenu) {
	p := &DaoPermitGroupMenuImpl{}
	p.SetContext(c...)
	res = p
	return
}
