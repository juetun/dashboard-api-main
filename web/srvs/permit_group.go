// Package srvs
/**
* @Author:changjiang
* @Description:
* @File:group
* @Version: 1.0.0
* @Date 2021/6/20 10:01 下午
 */
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitGroup interface {
	AdminUserGroupRelease(arg *wrappers.ArgAdminUserGroupRelease) (res wrappers.ResultAdminUserGroupRelease, err error)

	AdminUserGroupAdd(arg *wrappers.ArgAdminUserGroupAdd) (res wrappers.ResultAdminUserGroupAdd, err error)

	AdminGroup(arg *wrappers.ArgAdminGroup) (res *wrappers.ResultAdminGroup, err error)

	AdminGroupDelete(arg *wrappers.ArgAdminGroupDelete) (res wrappers.ResultAdminGroupDelete, err error)

	AdminGroupEdit(arg *wrappers.ArgAdminGroupEdit) (res *wrappers.ResultAdminGroupEdit, err error)

	MenuImportSet(arg *wrappers.ArgMenuImportSet) (res *wrappers.ResultMenuImportSet, err error)
}
