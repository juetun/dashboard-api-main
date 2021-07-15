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

type SrvPermitGroup interface{
	MenuImportSet(arg * wrappers.ArgMenuImportSet)(res *wrappers.ResultMenuImportSet,err error)
}
