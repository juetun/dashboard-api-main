package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroupMenuImpl struct {
	base.ServiceDao
}

func (r *DaoPermitGroupMenuImpl) DeleteGroupMenuByGroupIds(groupId ...int64) (err error) {
	if len(groupId) == 0 {
		return
	}
	var m models.AdminUserGroupMenu
	db := r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
		Where("group_id IN (?)", groupId)

	db = db.Where(db)
	if err = db.
		Delete(&m).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupId": groupId,
			"err":     err,
		}, "DaoPermitGroupMenuImplDeleteGroupMenuByGroupIds")
	}
	return
}

func (r *DaoPermitGroupMenuImpl) DeleteGroupPermitByMenuIds(groupId int64, module string, pageMenuId []int64) (err error) {
	if len(pageMenuId) == 0 {
		return
	}
	var m models.AdminUserGroupMenu
	db := r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt()).
		Where("group_id =? AND module=?", groupId, module)

	if len(pageMenuId) > 0 {
		db = db.Where("menu_id IN (?) AND path_type=?", pageMenuId, models.PathTypePage)
	}

	db = db.Where(db)
	if err = db.
		Delete(&m).
		Error; err != nil {
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
	if err = r.Context.Db.
		Table(m.TableName()).
		Select("distinct `menu_id`,`group_id`,`id`").
		Where("module = ? AND `group_id` in(?) ", module, groupIds).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupIds": groupIds,
			"err":      err,
		}, "DaoPermitGroupMenuImplGetMenuIdsByPermitByGroupIds")
		return
	}
	return
}

func (r *DaoPermitGroupMenuImpl) DeleteGroupMenuByGroupIdAndIds(groupId int64, ids ...int64) (err error) {
	if groupId == 0 {
		return
	}
	var m1 models.AdminUserGroupImport
	db := r.Context.Db.Table(m1.TableName()).Scopes(base.ScopesDeletedAt()).
		Where("group_id = ? ", groupId)
	if len(ids) > 0 {
		db = db.Where(" id IN( ? )", ids)
	}
	if err = db.
		Delete(&models.AdminImport{}).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"groupId": groupId,
			"ids":     ids,
			"err":     err.Error(),
		}, "DaoPermitGroupImportImplDeleteGroupMenuByGroupIdAndIds")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitGroupMenuImpl) BatchAddData(tableName string, list []base.ModelBase) (err error) {
	var data = base.BatchAddDataParameter{
		Data:      list,
		TableName: tableName,
		DbName:    r.Context.DbName,
		Db:        r.Context.Db,
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
