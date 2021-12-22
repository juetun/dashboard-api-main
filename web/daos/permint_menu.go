package daos

import (
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitMenu interface {
	GetMenuByCondition(condition interface{}) (res []models.AdminMenu, err error)

	GetMenu(menuId ...int64) (res []models.AdminMenu, err error)

	GetMenuByMenuId(menuId int64) (res models.AdminMenu, err error)
}
