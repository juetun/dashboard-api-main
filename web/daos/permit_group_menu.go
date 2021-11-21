package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroupMenu interface {
	base.DaoBatchAdd
	BatchAddData(tableName string, list []base.ModelBase) (err error)

	DeleteGroupMenuByGroupIdAndIds(groupId int64, ids ...int64) (err error)

	DeleteGroupMenuByGroupIds(groupId ...int64) (err error)

	DeleteGroupPermitByMenuIds(groupId int64, module string, pageMenuId []int64) (err error)

	GetMenuIdsByPermitByGroupIds(module string, groupIds ...int64) (res []models.AdminUserGroupMenu, err error)
}
