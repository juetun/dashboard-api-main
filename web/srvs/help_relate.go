package srvs

import "github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"

type SrvHelpRelate interface {
	HelpTree(arg *wrapper_admin.ArgHelpTree) (res *wrapper_admin.ResultHelpTree, err error)
	TreeEditNode(arg *wrapper_admin.ArgTreeEditNode) (res *wrapper_admin.ResultTreeEditNode, err error)
}
