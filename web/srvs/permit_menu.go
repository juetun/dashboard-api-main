// Package srvs
/**
* @Author:ChangJiang
* @Description:
* @File:permit_menu
* @Version: 1.0.0
* @Date 2021/9/12 1:52 下午
 */
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitMenu interface {

	// Menu 菜单
	Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error)

	AdminMenu(arg *wrappers.ArgAdminMenu) (res *wrappers.ResultAdminMenu, err error)

	GetMenuPermitKeyByPath(arg *wrappers.ArgGetImportByMenuIdSingle, dao daos.DaoPermit) (err error)

	AdminMenuWithCheck(arg *wrappers.ArgAdminMenuWithCheck) (res *wrappers.ResultMenuWithCheck, err error)
}
