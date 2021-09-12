/**
* @Author:changjiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/9/12 11:36 上午
 */
package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitImport interface {
	ImportList(arg *wrappers.ArgImportList) (res *wrappers.ResultImportList, err error)

	GetImport(arg *wrappers.ArgGetImport) (res *wrappers.ResultGetImport, err error)

	EditImport(arg *wrappers.ArgEditImport) (res *wrappers.ResultEditImport, err error)

	UpdateImportValue(arg *wrappers.ArgUpdateImportValue) (res *wrappers.ResultUpdateImportValue, err error)

	DeleteImport(arg *wrappers.ArgDeleteImport) (res *wrappers.ResultDeleteImport, err error)
}
