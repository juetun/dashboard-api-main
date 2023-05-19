package srvs

import "github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"

type SrvOperate interface {
	OperateLog(arg *wrapper_admin.ArgOperateLog) (res *wrapper_admin.ResultOperateLog,err error)
}
