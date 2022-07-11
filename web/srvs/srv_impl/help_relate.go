package srv_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type SrvHelpRelateImpl struct {
	base.ServiceBase
	dao daos.DaoHelpRelate
}

func (r *SrvHelpRelateImpl) HelpTree(arg *wrapper_admin.ArgHelpTree) (res wrapper_admin.ResultHelpTree, err error) {
	res = []*wrapper_admin.ResultHelpTreeItem{}
	if arg.BizCode == "" {
		return
	}
	var data map[int64][]*wrapper_admin.ResultHelpTreeItem
	if data, err = r.dao.GetByTopId(arg.BizCode, arg.TopId); err != nil {
		return
	}
	var ok bool
	if res, ok = data[arg.TopId]; !ok {
		res = []*wrapper_admin.ResultHelpTreeItem{}
	}
	return
}

func (r *SrvHelpRelateImpl) TreeEditNode(arg *wrapper_admin.ArgTreeEditNode) (res *wrapper_admin.ResultTreeEditNode, err error) {

	res = &wrapper_admin.ResultTreeEditNode{}

	var data = &models.HelpDocumentRelate{}
	data.Id = arg.Id
	data.Display = arg.Display
	data.ParentId = arg.ParentId
	data.IsLeafNode = arg.IsLeafNode
	data.DocKey = arg.DocKey
	data.Label = arg.Label
	data.BizCode = arg.BizCode
	data.CreatedAt = arg.TimeNow
	data.UpdatedAt = arg.TimeNow
	if err = r.dao.AddOneHelpRelate(data); err != nil {
		return
	}
	res.Result = true
	return
}

func NewSrvHelpRelate(context ...*base.Context) (res srvs.SrvHelpRelate) {
	p := &SrvHelpRelateImpl{}
	p.SetContext(context...)
	p.dao = dao_impl.NewDaoHelpRelate(p.Context)
	return p
}
