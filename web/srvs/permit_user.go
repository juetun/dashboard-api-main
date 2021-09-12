/**
* @Author:changjiang
* @Description:
* @File:permit_user
* @Version: 1.0.0
* @Date 2021/9/12 12:10 下午
 */
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitUser interface {

	AdminUserAdd(arg *wrappers.ArgAdminUserAdd) (res wrappers.ResultAdminUserAdd, err error)
	AdminUser(arg *wrappers.ArgAdminUser) (res *wrappers.ResultAdminUser, err1 error)

	AdminUserDelete(arg *wrappers.ArgAdminUserDelete) (res wrappers.ResultAdminUserDelete, err error)
}
