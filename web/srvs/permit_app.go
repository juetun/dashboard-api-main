/**
* @Author:changjiang
* @Description:
* @File:permit_app
* @Version: 1.0.0
* @Date 2021/9/12 1:33 下午
 */
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitApp interface {
	GetAppConfig(arg *wrappers.ArgGetAppConfig) (res *wrappers.ResultGetAppConfig, err error)

}