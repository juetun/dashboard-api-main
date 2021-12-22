// Package srvs
/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2021/6/6 10:27 下午
 */
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type PermitService interface {
	AdminSetPermit(arg *wrappers.ArgAdminSetPermit) (res *wrappers.ResultAdminSetPermit, err error)

	MenuAdd(arg *wrappers.ArgMenuAdd) (res *wrappers.ResultMenuAdd, err error)

	MenuSave(arg *wrappers.ArgMenuSave) (res *wrappers.ResultMenuSave, err error)

	MenuDelete(arg *wrappers.ArgMenuDelete) (res *wrappers.ResultMenuDelete, err error)

	MenuImport(arg *wrapper_admin.ArgMenuImport) (res *wrapper_admin.ResultMenuImport, err error)

	AdminMenuSearch(arg *wrappers.ArgAdminMenu) (res wrappers.ResAdminMenuSearch, err error)

	GetMenu(arg *wrappers.ArgGetMenu) (res wrappers.ResultGetMenu, err error)

	// GetImportByMenuId(arg *wrappers.ArgGetImportByMenuId) (res wrappers.ResultGetImportByMenuId, err error)
}
