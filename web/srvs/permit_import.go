// Package srvs /**
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type SrvPermitImport interface {
	SetApiPermit(arg *wrappers.ArgAdminSetPermit) (err error)

	ImportList(arg *wrappers.ArgImportList) (res *wrappers.ResultImportList, err error)

	GetImport(arg *wrappers.ArgGetImport) (res *wrappers.ResultGetImport, err error)

	EditImport(arg *wrappers.ArgEditImport) (res *wrappers.ResultEditImport, err error)

	GetImportByMenuId(arg *wrappers.ArgGetImportByMenuId) (res wrappers.ResultGetImportByMenuId, err error)

	GetOpList(dao daos.DaoPermit, arg *wrappers.ArgPermitMenu) (opList map[string][]wrappers.OpOne, err error)

	UpdateImportValue(arg *wrappers.ArgUpdateImportValue) (res *wrappers.ResultUpdateImportValue, err error)

	DeleteImport(arg *wrappers.ArgDeleteImport) (res *wrappers.ResultDeleteImport, err error)

	GetChildMenu(nowMenuId int64) (menuIds []wrappers.MenuSingle, err error)

	GetChildImport(nowMenuId int64) (importIds []wrappers.ImportSingle, err error)

	//用户页面具备的接口权限
	UserPageImport(arg *wrapper_admin.ArgPageImport) (res *wrapper_admin.ResultPageImport, err error)
}
