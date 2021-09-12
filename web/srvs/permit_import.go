/**
* @Author:changjiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/9/12 11:36 上午
 */
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitImport interface {

	SetApiPermit(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (err error)

	ImportList(arg *wrappers.ArgImportList) (res *wrappers.ResultImportList, err error)

	GetImport(arg *wrappers.ArgGetImport) (res *wrappers.ResultGetImport, err error)

	EditImport(arg *wrappers.ArgEditImport) (res *wrappers.ResultEditImport, err error)

	GetImportByMenuId(arg *wrappers.ArgGetImportByMenuId) (res wrappers.ResultGetImportByMenuId, err error)

	GetOpList(dao daos.DaoPermit, arg *wrappers.ArgPermitMenu) (opList map[string][]wrappers.OpOne, err error)

	UpdateImportValue(arg *wrappers.ArgUpdateImportValue) (res *wrappers.ResultUpdateImportValue, err error)

	DeleteImport(arg *wrappers.ArgDeleteImport) (res *wrappers.ResultDeleteImport, err error)

	GetChildMenu(nowMenuId int) (menuIds []int, err error)

	GetChildImport(nowMenuId int) (importIds []int, err error)
}
