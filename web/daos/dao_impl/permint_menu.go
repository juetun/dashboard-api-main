package dao_impl

import (
	"fmt"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type DaoPermitMenuImpl struct {
	base.ServiceDao
}

func (r *DaoPermitMenuImpl) GetMenuMap(menuId ...int64) (res map[int64]*models.AdminMenu, err error) {
	res = map[int64]*models.AdminMenu{}
	var mapMenu []*models.AdminMenu
	if mapMenu, err = r.GetMenu(menuId...); err != nil {
		return
	}
	for _, item := range mapMenu {
		res[item.Id] = item
	}
	return
}

func (r *DaoPermitMenuImpl) GetAdminMenuByModule(module string) (res models.AdminMenu, err error) {
	var list []models.AdminMenu
	if list, err = r.GetMenuByPermitKey("", module); err != nil {
		return
	}
	if len(list) > 0 {
		res = list[0]
	}

	return
}

func (r *DaoPermitMenuImpl) GetAllSystemList() (res []*models.AdminMenu, err error) {
	res, err = r.GetByCondition(map[string]interface{}{"module": wrappers.DefaultPermitModule}, []wrappers.DaoOrderBy{
		{
			Column:     "sort_value",
			SortFormat: "asc",
		},
	}, 0)
	return
}

func (r *DaoPermitMenuImpl) GetMenuByPermitKey(module string, permitKey ...string) (res []models.AdminMenu, err error) {

	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"module":    module,
			"permitKey": permitKey,
			"err":       err.Error(),
		}, "DaoPermitMenuImplGetMenuByPermitKey")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenu
		actRes = r.GetDefaultActErrorHandlerResult(m)
		db := actRes.Db.Table(actRes.TableName).Scopes(base.ScopesDeletedAt())
		lPk := len(permitKey)
		var f bool
		if lPk != 0 {
			f = true
			db = db.Where("permit_key IN(?)", permitKey)
		}
		if module != "" {
			f = true
			db = db.Where("module = ?", module)
		}
		if !f {
			return
		}
		actRes.Err = db.Limit(len(permitKey)).
			Order("sort_value desc").
			Find(&res).
			Error
		return
	})
	return
}

func (r *DaoPermitMenuImpl) Add(data *models.AdminMenu) (err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err.Error(),
		}, "DaoPermitMenuImplAdd")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(data)
		actRes.Err = actRes.Db.Create(data).Error
		return
	})
	return
}

func (r *DaoPermitMenuImpl) UpdateMenuByCondition(condition interface{}, data map[string]interface{}) (err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"data":      data,
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitUpdateMenuByCondition")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenu
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where(condition).
			Updates(data).
			Error
		return
	})
	return
}

func (r *DaoPermitMenuImpl) Save(id int64, data map[string]interface{}) (err error) {
	if id == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"id":   id,
			"data": data,
			"err":  err.Error(),
		}, "DaoPermitMenuImplSave")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenu
		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.
			Table(actRes.TableName).
			Where("id = ?", id).
			Updates(data).
			Error
		return
	})
	return
}
func (r *DaoPermitMenuImpl) GetByCondition(condition map[string]interface{}, orderBy []wrappers.DaoOrderBy, limit int) (res []*models.AdminMenu, err error) {

	if len(condition) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "daoPermitGetByCondition")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()

	var actHandler = func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenu
		actRes = r.GetDefaultActErrorHandlerResult(m)
		db := actRes.Db.Table(actRes.TableName).
			Where(condition)

		if len(orderBy) > 0 {
			var orderString string
			for _, item := range orderBy {
				orderString += fmt.Sprintf("%s %s", item.Column, item.SortFormat)
			}
			if orderString != "" {
				db = db.Order(orderString)
			}
		}
		if limit > 0 {
			db = db.Limit(limit)
		}
		actRes.Err = db.Find(&res).Error
		return
	}

	err = r.ActErrorHandler(actHandler)
	return
}

func (r *DaoPermitMenuImpl) GetMenuByMenuId(menuId int64) (res *models.AdminMenu, err error) {
	res = &models.AdminMenu{}
	var li []*models.AdminMenu
	if li, err = r.GetMenu(menuId); err != nil {
		return
	}
	if len(li) > 0 {
		res = li[0]
	}
	return
}

func (r *DaoPermitMenuImpl) GetMenu(menuId ...int64) (res []*models.AdminMenu, err error) {
	l := len(menuId)
	res = make([]*models.AdminMenu, 0, l)
	if l == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err.Error(),
		}, "daoPermitGetMenu")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		var m models.AdminMenu
		actRes = r.GetDefaultActErrorHandlerResult(&m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Where("id IN (?)", menuId).
			Limit(len(menuId)).
			Find(&res).
			Error
		return
	})
	return
}

func (r *DaoPermitMenuImpl) GetMenuByCondition(condition interface{}) (res []models.AdminMenu, err error) {
	if condition == nil {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err.Error(),
		}, "DaoPermitMenuImplGetMenuByCondition")
	}()
	var actHandler = func() (actRes *base.ActErrorHandlerResult) {
		var m *models.AdminMenu

		actRes = r.GetDefaultActErrorHandlerResult(m)
		actRes.Err = actRes.Db.Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where(condition).
			Limit(1000).
			Find(&res).
			Error
		return
	}
	err = r.ActErrorHandler(actHandler)
	return
}

func NewDaoPermitMenu(ctx ...*base.Context) (res daos.DaoPermitMenu) {
	p := &DaoPermitMenuImpl{}
	p.SetContext(ctx...)
	return p
}
