package srvs

import (
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_outernet"
)

type SrvHelpRelate interface {
	HelpTree(arg *wrapper_admin.ArgHelpTree) (res  *wrapper_admin.ResultHelpTree, err error)

	TreeEditNode(arg *wrapper_admin.ArgTreeEditNode) (res *wrapper_admin.ResultTreeEditNode, err error)

	Tree(arg *wrapper_outernet.ArgTree) (res *wrapper_outernet.ResultTree, err error)
}
