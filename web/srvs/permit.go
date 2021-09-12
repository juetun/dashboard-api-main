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
)

type PermitService interface {
	AdminGroupEdit(arg *wrappers.ArgAdminGroupEdit) (res *wrappers.ResultAdminGroupEdit, err error)

	Flag(arg *wrappers.ArgFlag) (res *wrappers.ResultFlag, err error)

	GetAppConfig(arg *wrappers.ArgGetAppConfig) (res *wrappers.ResultGetAppConfig, err error)

	Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error)

	AdminSetPermit(arg *wrappers.ArgAdminSetPermit) (res *wrappers.ResultAdminSetPermit, err error)

	AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error)

	AdminMenuWithCheck(arg *wrappers.ArgAdminMenuWithCheck) (res *wrappers.ResultMenuWithCheck, err error)

	MenuAdd(arg *wrappers.ArgMenuAdd) (res *wrappers.ResultMenuAdd, err error)

	MenuSave(arg *wrappers.ArgMenuSave) (res *wrappers.ResultMenuSave, err error)

	MenuDelete(arg *wrappers.ArgMenuDelete) (res *wrappers.ResultMenuDelete, err error)

	MenuImport(arg *wrappers.ArgMenuImport) (res *wrappers.ResultMenuImport, err error)

	AdminMenuSearch(arg *wrappers.ArgAdminMenu) (res wrappers.ResAdminMenuSearch, err error)

	GetMenu(arg *wrappers.ArgGetMenu) (res wrappers.ResultGetMenu, err error)

	GetImportByMenuId(arg *wrappers.ArgGetImportByMenuId) (res wrappers.ResultGetImportByMenuId, err error)
}
