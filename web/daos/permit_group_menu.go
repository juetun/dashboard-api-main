package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroupMenu interface {
	base.DaoBatchAdd
	BatchAddDataUserGroupMenu(list []base.ModelBase) (err error)

	DeleteGroupMenuByGroupIdAndIds(groupId int64, ids ...int64) (err error)

	// DeleteGroupMenuByGroupIds 删除组对应的界面权限
	DeleteGroupMenuByGroupIds(groupId ...int64) (err error)

	// DeleteGroupImportByGroupIds 删除组对应的接口权限
	DeleteGroupImportByGroupIds(groupId ...int64) (err error)

	DeleteGroupPermitByMenuIds(groupId int64, module string, pageMenuId []int64) (err error)

	GetMenuIdsByPermitByGroupIds(module string, groupIds ...int64) (res []models.AdminUserGroupMenu, err error)
}
