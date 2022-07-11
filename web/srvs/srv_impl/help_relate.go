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
	res = []*wrapper_admin.ResultFormPage{}
	if arg.BizCode == "" {
		return
	}
	var (
		data    map[int64][]*wrapper_admin.ResultHelpTreeItem
		ok      bool
		dataRes []*wrapper_admin.ResultHelpTreeItem
	)
	if data, err = r.dao.GetByTopId(arg.BizCode, arg.TopId); err != nil {
		return
	}
	if dataRes, ok = data[arg.TopId]; !ok {
		res = []*wrapper_admin.ResultFormPage{}
	}
	var mapParent = make(map[int64]int64)

	res = r.setExpand(dataRes, arg.CurrentId, &mapParent)
	var currentId = arg.CurrentId
	for {
		if len(mapParent) == 0 {
			break
		}
		if currentId, ok = mapParent[currentId]; !ok || currentId == 0 {
			break
		}
		delete(mapParent, currentId)
		r.setAllParentExpand(res, currentId)

	}
	return
}

func (r *SrvHelpRelateImpl) setAllParentExpand(res []*wrapper_admin.ResultFormPage, currentId int64) {
	for _, value := range res {
		if value.Id == currentId {
			value.Expand = true
		}
		if len(value.Children) > 0 {
			r.setAllParentExpand(value.Children, currentId)
		}
	}

	return
}

func (r *SrvHelpRelateImpl) setExpand(data []*wrapper_admin.ResultHelpTreeItem, currentId int64, mapParent *map[int64]int64) (res []*wrapper_admin.ResultFormPage) {
	var dataItem *wrapper_admin.ResultFormPage
	res = make([]*wrapper_admin.ResultFormPage, 0, len(data))
	for _, value := range data {
		(*mapParent)[value.Id] = value.ParentId
		dataItem = &wrapper_admin.ResultFormPage{
			Title:      value.Label,
			Id:         value.Id,
			DocKey:     value.DocKey,
			Display:    value.Display,
			IsLeafNode: value.IsLeafNode,
		}
		if value.Id == currentId {
			dataItem.Expand = true
		}
		if len(value.Child) > 0 {
			dataItem.Children = r.setExpand(value.Child, currentId, mapParent)
		}
		res = append(res, dataItem)
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
