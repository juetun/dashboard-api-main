// Package srvs
/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 11:29 下午
 */
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvService interface {
	Edit(arg *wrappers.ArgServiceEdit) (res *wrappers.ResultServiceEdit, err error)

	List(arg *wrappers.ArgServiceList) (res *wrappers.ResultServiceList, err error)

	Detail(arg *wrappers.ArgDetail) (res *wrappers.ResultDetail, err error)
}
