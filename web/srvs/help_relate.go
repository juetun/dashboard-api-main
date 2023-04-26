package srvs

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_outernet"
)

type SrvHelpRelate interface {
	HelpTree(arg *wrapper_admin.ArgHelpTree) (res *wrapper_admin.ResultHelpTree, err error)

	TreeEditNode(arg *wrapper_admin.ArgTreeEditNode) (res *wrapper_admin.ResultTreeEditNode, err error)

	//获取帮助详情
	GetDescMedia(desc string, argCommon *base.GetDataTypeCommon) (resDesc string, err error)

	Tree(arg *wrapper_outernet.ArgTree) (res *wrapper_outernet.ResultTree, err error)
}
