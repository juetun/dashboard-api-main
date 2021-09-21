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
	Menu(arg *wrappers.ArgPermitMenu) (res *wrappers.ResultPermitMenuReturn, err error)

	GetMenuPermitKeyByPath(arg *wrappers.ArgGetImportByMenuIdSingle, dao daos.DaoPermit) (err error)
}
