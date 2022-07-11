package srvs

import "github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"

type SrvHelp interface {
	HelpList(arg *wrapper_admin.ArgHelpList)(res  *wrapper_admin.ResultHelpList,err error)
	HelpDetail(arg *wrapper_admin.ArgHelpDetail)(res  *wrapper_admin.ResultHelpDetail,err error)
	HelpEdit(arg *wrapper_admin.ArgHelpEdit)(res  *wrapper_admin.ResultHelpEdit,err error)
}