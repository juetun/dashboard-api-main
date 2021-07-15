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
	DeleteImport(arg *wrappers.ArgDeleteImport) (res *wrappers.ResultDeleteImport, err error)

	AdminGroupDelete(arg *wrappers.ArgAdminGroupDelete) (res wrappers.ResultAdminGroupDelete, err error)

	AdminUserGroupRelease(arg *wrappers.ArgAdminUserGroupRelease) (res wrappers.ResultAdminUserGroupRelease, err error)

	AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error)

	AdminUserAdd(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error)

	AdminGroupEdit(arg *wrappers.ArgAdminGroupEdit) (res *wrappers.ResultAdminGroupEdit, err error)

	AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err1 error)

	Flag(arg *wrappers.ArgFlag) (res *wrappers.ResultFlag, err error)

	GetAppConfig(arg *wrappers.ArgGetAppConfig) (res *wrappers.ResultGetAppConfig, err error)

	Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error)

	AdminGroup(arg *wrappers.ArgAdminGroup) (res *wrappers.ResultAdminGroup, err error)

	AdminSetPermit(arg *wrappers.ArgAdminSetPermit) (res *wrappers.ResultAdminSetPermit, err error)

	AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error)

	AdminMenuWithCheck(arg *wrappers.ArgAdminMenuWithCheck) (res *wrappers.ResultMenuWithCheck, err error)

	MenuAdd(arg *wrappers.ArgMenuAdd) (res *wrappers.ResultMenuAdd, err error)

	MenuSave(arg *wrappers.ArgMenuSave) (res *wrappers.ResultMenuSave, err error)

	MenuDelete(arg *wrappers.ArgMenuDelete) (res *wrappers.ResultMenuDelete, err error)

	GetImport(arg *wrappers.ArgGetImport) (res *wrappers.ResultGetImport, err error)

	MenuImport(arg *wrappers.ArgMenuImport) (res *wrappers.ResultMenuImport, err error)

	EditImport(arg *wrappers.ArgEditImport) (res *wrappers.ResultEditImport, err error)

	ImportList(arg *wrappers.ArgImportList) (res *wrappers.ResultImportList, err error)

	AdminUserGroupAdd(arg *wrappers.ArgAdminUserGroupAdd) (res wrappers.ResultAdminUserGroupAdd, err error)

	AdminMenuSearch(arg *wrappers.ArgAdminMenu) (res wrappers.ResAdminMenuSearch, err error)

	GetMenu(arg *wrappers.ArgGetMenu) (res wrappers.ResultGetMenu, err error)

	UpdateImportValue(arg *wrappers.ArgUpdateImportValue) (res *wrappers.ResultUpdateImportValue, err error)
}
