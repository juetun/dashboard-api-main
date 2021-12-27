package daos

import (
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type DaoPermitMenu interface {
	GetMenuByCondition(condition interface{}) (res []models.AdminMenu, err error)

	UpdateMenuByCondition(condition interface{}, data map[string]interface{}) (err error)

	Save(id int64, data map[string]interface{}) (err error)

	GetByCondition(condition map[string]interface{}, orderBy []wrappers.DaoOrderBy, limit int) (res []*models.AdminMenu, err error)

	GetMenu(menuId ...int64) (res []models.AdminMenu, err error)

	GetMenuByMenuId(menuId int64) (res models.AdminMenu, err error)

	//
	GetAllSystemList() (res []*models.AdminMenu, err error)

	Add(data *models.AdminMenu) (err error)


	GetMenuByPermitKey(module string, permitKey ...string) (res []models.AdminMenu, err error)

	GetAdminMenuByModule(module string)(res models.AdminMenu,err error)
}
