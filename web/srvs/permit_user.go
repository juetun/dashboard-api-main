// Package srvs /**
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type SrvPermitUser interface {
	AdminUserEdit(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error)

	AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err1 error)

	AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error)

	AdminUserUpdateWithColumn(arg *wrapper_admin.ArgAdminUserUpdateWithColumn) (res *wrapper_admin.ResultAdminUserUpdateWithColumn, err error)

	// GetUserAdminGroupIdByUserHid 获取用户的后台权限组
	GetUserAdminGroupIdByUserHid(uid string) (groupId []int64, isSuperAdmin bool, err error)
}
