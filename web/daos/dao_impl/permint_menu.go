package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitMenuImpl struct {
	base.ServiceDao
}

func (r *DaoPermitMenuImpl) GetMenu(menuId ...int64) (res []models.AdminMenu, err error) {
	var m models.AdminMenu
	l := len(menuId)
	res = make([]models.AdminMenu, 0, l)
	if l == 0 {
		return
	}
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
