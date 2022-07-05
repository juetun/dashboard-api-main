package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoPermitGroupImport interface {
	base.DaoBatchAdd

	GetSelectImportByImportId(groupId int64, importId ...int64) (res []models.AdminUserGroupImport, err error)

	//菜单下的有权限的接口列表
	GetMenuIdsByPermitByGroupIds(menuIds, groupIds []int64, ) (res []*models.AdminUserGroupImport, err error)

	BatchAddDataGroupImport(list []base.ModelBase) (err error)

	DeleteGroupImportWithMenuId(menuId ...int64) (err error)

	DeleteGroupImportWithGroupIdAndImportIds(groupId int64, importIds ...int64) (err error)

	DeleteGroupMenuPermitByGroupIdAndMenuIds(groupId int64, menuIds ...int64) (err error)
}
