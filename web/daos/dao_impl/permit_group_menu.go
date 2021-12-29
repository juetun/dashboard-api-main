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
	var m models.AdminUserGroupMenu
	db := r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
		Where("group_id = ? AND module = ?", groupId, module)

	if len(pageMenuId) > 0 {
		db = db.Where("menu_id IN (?)", pageMenuId)
	}
	if err = db.Delete(&m).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupId":    groupId,
			"pageMenuId": pageMenuId,
			"module":     module,
			"err":        err,
		}, "DaoPermitGroupMenuImplDeleteGroupPermitByMenuIds")
	}
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

	var m models.AdminUserGroupMenu
	if err = r.ActErrorHandler(func() (actErrorHandlerResult *base.ActErrorHandlerResult) {
		actErrorHandlerResult = r.GetDefaultActErrorHandlerResult(&m)

		db := r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
			Where("group_id = ? ", groupId)
		if len(menuIds) > 0 {
			db = db.Where("menu_id IN(?)", menuIds)
		}
		if actErrorHandlerResult.Err = db.
			Delete(&models.AdminImport{}).Error; actErrorHandlerResult.Err != nil {
			r.Context.Error(map[string]interface{}{
				"groupId": groupId,
				"menuIds": menuIds,
				"err":     actErrorHandlerResult.Err.Error(),
			}, "DaoPermitGroupImportImplDeleteGroupMenuByGroupIdAndIds")
			return
		}
		return
	}); err != nil {
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}

	return
}

func (r *DaoPermitGroupMenuImpl) BatchAddData(tableName string, list []base.ModelBase) (err error) {
	var data = base.BatchAddDataParameter{
		Data: list,
		CommonDb: base.CommonDb{
			TableName: tableName,
			DbName:    r.Context.DbName,
			Db:        r.Context.Db,
		},
	}
	if err = r.BatchAdd(&data); err != nil {
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err.Error(),
		}, "DaoPermitGroupMenuImplBatchAddData")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func NewDaoPermitGroupMenu(c ...*base.Context) (res daos.DaoPermitGroupMenu) {
	p := &DaoPermitGroupMenuImpl{}
	p.SetContext(c...)
	res = p
	return
}
