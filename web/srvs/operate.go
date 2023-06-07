package srvs

import (
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"github.com/juetun/library/common/app_param/operate_admin"
)

type SrvOperate interface {
	AddLog(list []*models.OperateLog) (err error)

	OperateLog(arg *wrapper_admin.ArgOperateLog) (res *wrapper_admin.ResultOperateLog, err error)

	AddLogMessage(args *operate_admin.ArgAddAdminLog) (res *operate_admin.ResultAddAdminLog, err error)
}
