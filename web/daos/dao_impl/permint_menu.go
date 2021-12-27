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

	var m models.AdminMenu

	db := r.Context.Db.Table(m.TableName()).Scopes(base.ScopesDeletedAt())
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

	if err = db.Limit(len(permitKey)).
		Order("sort_value desc").
		Find(&res).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"permitKey": permitKey,
			"err":       err.Error(),
		}, "daoPermitGetMenuByPermitKey")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
		return
	}
	return
}

func (r *DaoPermitMenuImpl) Add(data *models.AdminMenu) (err error) {

	if err = r.Context.Db.Create(data).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"data": data,
			"err":  err,
		}, "DaoPermitMenuImplAdd")
	}
	return
}

func (r *DaoPermitMenuImpl) UpdateMenuByCondition(condition interface{}, data map[string]interface{}) (err error) {
	var m models.AdminMenu
	if err = r.Context.Db.
		Table(m.TableName()).
		Where(condition).
		Updates(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"data":      data,
			"err":       err,
		}, "daoPermitUpdateMenuByCondition")
		return
	}
	return
}

func (r *DaoPermitMenuImpl) Save(id int64, data map[string]interface{}) (err error) {
	if id == 0 {
		return
	}
	if err = r.Context.Db.
		Model(&models.AdminMenu{}).
		Where("id = ?", id).
		Updates(data).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"id":   id,
			"data": data,
			"err":  err,
		}, "DaoPermitMenuImplSave")
	}
	return
}
func (r *DaoPermitMenuImpl) GetByCondition(condition map[string]interface{}, orderBy []wrappers.DaoOrderBy, limit int) (res []*models.AdminMenu, err error) {

	if len(condition) == 0 {
		return
	}
	var m models.AdminMenu
	db := r.Context.Db.
		Table(m.TableName()).
		Where(condition)
	var orderString string
	for _, item := range orderBy {
		orderString += fmt.Sprintf("%s %s", item.Column, item.SortFormat)
	}
	if orderString != "" {
		db = db.Order(orderString)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	if err = db.
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err,
		}, "daoPermitGetByCondition")
		return
	}
	return
}

func (r *DaoPermitMenuImpl) GetMenuByMenuId(menuId int64) (res models.AdminMenu, err error) {
	var li []models.AdminMenu
	if li, err = r.GetMenu(menuId); err != nil {
		return
	}
	if len(li) > 0 {
		res = li[0]
	}
	return
}

func (r *DaoPermitMenuImpl) GetMenu(menuId ...int64) (res []models.AdminMenu, err error) {
	l := len(menuId)
	res = make([]models.AdminMenu, 0, l)
	if l == 0 {
		return
	}
	var m models.AdminMenu
	if err = r.Context.Db.Table(m.TableName()).
		Where("id IN (?)", menuId).
		Limit(len(menuId)).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"menuId": menuId,
			"err":    err,
		}, "daoPermitGetMenu")
		return
	}
	return
}
func (r *DaoPermitMenuImpl) GetMenuByCondition(condition interface{}) (res []models.AdminMenu, err error) {
	if condition == nil {
		return
	}
	var m models.AdminMenu
	if err = r.Context.Db.Table(m.TableName()).
		Scopes(base.ScopesDeletedAt()).
		Where(condition).
		Limit(1000).
		Find(&res).
		Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"condition": condition,
			"err":       err,
		}, "DaoPermitMenuImplGetMenuByCondition")
		return
	}
	return
}

func NewDaoPermitMenu(ctx ...*base.Context) (res daos.DaoPermitMenu) {
	p := &DaoPermitMenuImpl{}
	p.SetContext(ctx...)
	return p
}
