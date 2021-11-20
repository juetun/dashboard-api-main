package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type SrvGatewayImportPermit interface {

	GetImportPermit(w *wrapper_intranet.ArgGetImportPermit) (res wrapper_intranet.ResultGetImportPermit, err error)

	GetUerImportPermit(arg *wrapper_intranet.ArgGetUerImportPermit) (res *wrapper_intranet.ResultGetUerImportPermit, err error)
}
