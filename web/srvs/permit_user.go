// Package srvs /**
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type SrvPermitUser interface {
	AdminUserEdit(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error)

	AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err1 error)

	AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error)

	AdminUserUpdateWithColumn(arg *wrapper_admin.ArgAdminUserUpdateWithColumn) (res *wrapper_admin.ResultAdminUserUpdateWithColumn, err error)

	// GetUserAdminGroupIdByUserHid 获取用户的后台权限组
	GetUserAdminGroupIdByUserHid(uid int64) (groupId []int64, isSuperAdmin bool, err error)

	ValidateUserHavePermit(args *wrapper_intranet.ArgValidateUserHavePermit) (res *wrapper_intranet.ResultValidateUserHavePermit,err error)
}
