package srv_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_outernet"
)

type SrvHelpRelateImpl struct {
	base.ServiceBase
	dao daos.DaoHelpRelate
}

func (r *SrvHelpRelateImpl) Tree(arg *wrapper_outernet.ArgTree) (res *wrapper_outernet.ResultTree, err error) {
	res = &wrapper_outernet.ResultTree{}
	var (
		dataRes   []*wrapper_admin.ResultHelpTreeItem
		ok        bool
		mapParent = make(map[int64]int64)
		currentId = arg.CurrentId
		docKey    string
	)

	if dataRes, err = r.getDataWithBizAndTopId(arg.BizCode, arg.TopId); err != nil {
		return
	}

	res.Data, docKey = r.setExpandAsTree(dataRes, arg.CurrentId, &mapParent)

	if docKey != "" {
		var mapDoc map[string]*models.HelpDocument
		argFetchDoc := base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(docKey))
		if mapDoc, err = dao_impl.NewDaoHelp(r.Context).GetByPKey(argFetchDoc); err != nil {
			return
		}
		if tmp, ok := mapDoc[docKey]; ok {
			res.DocContent = tmp.Content
		}
	}
	for {
		if len(mapParent) == 0 {
			break
		}
		if currentId, ok = mapParent[currentId]; !ok || currentId == 0 {
			break
		}
		delete(mapParent, currentId)
		r.setAllParentExpandAsTree(res.Data, currentId)
	}
	return
}

func (r *SrvHelpRelateImpl) setAllParentExpandAsTree(res []*wrapper_outernet.ResultFormPage, currentId int64) {
	for _, value := range res {
		if value.Id == currentId {
			value.Expand = true
		}
		if len(value.Children) > 0 {
			r.setAllParentExpandAsTree(value.Children, currentId)
		}
	}

	return
}
func (r *SrvHelpRelateImpl) setExpandAsTree(data []*wrapper_admin.ResultHelpTreeItem, currentId int64, mapParent *map[int64]int64) (res []*wrapper_outernet.ResultFormPage, docKey string) {
	var dataItem *wrapper_outernet.ResultFormPage
	res = make([]*wrapper_outernet.ResultFormPage, 0, len(data))
	for _, value := range data {
		if value.Display == models.HelpDocumentRelateIsLeafNodeNo {
			continue
		}
		(*mapParent)[value.Id] = value.ParentId
		dataItem = &wrapper_outernet.ResultFormPage{
			Title:      value.Label,
			Id:         value.Id,
			DocKey:     value.DocKey,
			IsLeafNode: value.IsLeafNode,
		}
		if value.Id == currentId {
			dataItem.Expand = true
		}
		if value.IsLeafNode == models.HelpDocumentRelateIsLeafNodeYes {
			docKey = value.DocKey
		}
		if len(value.Child) > 0 {
			var docKeyTmp string
			dataItem.Children, docKeyTmp = r.setExpandAsTree(value.Child, currentId, mapParent)
			if docKeyTmp != "" {
				docKey = docKeyTmp
			}
		}
		res = append(res, dataItem)
	}
	return
}

func (r *SrvHelpRelateImpl) getDataWithBizAndTopId(bizCode string, topId int64) (dataRes []*wrapper_admin.ResultHelpTreeItem, err error) {
	var (
		ok   bool
		data map[int64][]*wrapper_admin.ResultHelpTreeItem
	)

	if data, err = r.dao.GetByTopId(bizCode, topId); err != nil {
		return
	}
	if dataRes, ok = data[topId]; !ok {
		return
	}
	return
}

func (r *SrvHelpRelateImpl) getTopHelp() (res []*models.HelpDocumentRelate, err error) {
	res, err = r.dao.GetByTopHelp()
	return
}

func (r *SrvHelpRelateImpl) orgAdminHelpTree(dataRes []*wrapper_admin.ResultHelpTreeItem, arg *wrapper_admin.ArgHelpTree, res *wrapper_admin.ResultHelpTree) (err error) {
	var ok bool
	var mapParent = make(map[int64]int64)
	var haveData bool

	for _, item := range dataRes {
		if item.Id == arg.CurrentId {
			dataRes = item.Child
			haveData = true
			break
		}
	}
	if !haveData {
		return
	}
	res.List = r.setExpand(dataRes, arg, &mapParent)
	var currentId = arg.CurrentId
	for {
		if len(mapParent) == 0 {
			break
		}
		if currentId, ok = mapParent[currentId]; !ok || currentId == 0 {
			break
		}
		delete(mapParent, currentId)
		r.setAllParentExpand(res.List, currentId)
	}
	return
}

func (r *SrvHelpRelateImpl) defaultHelpTreeActive(arg *wrapper_admin.ArgHelpTree, menu []*models.HelpDocumentRelate) (res []*wrapper_admin.ResultHelpTreeItemMenu) {
	res = make([]*wrapper_admin.ResultHelpTreeItemMenu, 0, len(menu))
	var dt *wrapper_admin.ResultHelpTreeItemMenu
	if arg.BizCode == "" {
		for k, item := range menu {
			dt = &wrapper_admin.ResultHelpTreeItemMenu{}
			dt.SetHelpDocumentRelate(item)
			if k == 0 {
				dt.Active = true
				arg.BizCode = item.BizCode
				arg.CurrentId = item.Id
			}
			res = append(res, dt)
		}
		return
	}
	for _, item := range menu {
		dt = &wrapper_admin.ResultHelpTreeItemMenu{}
		dt.SetHelpDocumentRelate(item)
		if item.BizCode == arg.BizCode {
			dt.Active = true
			arg.BizCode = item.BizCode
			arg.CurrentId = item.Id
		}
		res = append(res, dt)
	}
	return
}

func (r *SrvHelpRelateImpl) HelpTree(arg *wrapper_admin.ArgHelpTree) (res *wrapper_admin.ResultHelpTree, err error) {
	res = wrapper_admin.NewResultHelpTree()
	var menu []*models.HelpDocumentRelate
	if menu, err = r.getTopHelp(); err != nil {
		return
	} else if len(menu) == 0 {
		return
	}
	//初始化默认值
	res.Menu = r.defaultHelpTreeActive(arg, menu)

	var (
		dataRes []*wrapper_admin.ResultHelpTreeItem
	)

	if dataRes, err = r.getDataWithBizAndTopId(arg.BizCode, arg.TopId); err != nil {
		return
	} else if len(dataRes) == 0 {
		return
	}

	err = r.orgAdminHelpTree(dataRes, arg, res)
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

func (r *SrvHelpRelateImpl) setExpand(data []*wrapper_admin.ResultHelpTreeItem, arg *wrapper_admin.ArgHelpTree, mapParent *map[int64]int64) (res []*wrapper_admin.ResultFormPage) {
	var dataItem *wrapper_admin.ResultFormPage
	res = make([]*wrapper_admin.ResultFormPage, 0, len(data))
	for _, value := range data {
		(*mapParent)[value.Id] = value.ParentId

		dataItem = wrapper_admin.NewResultFormPage().
			SetResultHelpTreeItem(value, arg.CurrentId)
		if len(value.Child) > 0 {
			dataItem.Children = r.setExpand(value.Child, arg, mapParent)
		}
		res = append(res, dataItem)
	}
	return
}

func (r *SrvHelpRelateImpl) TreeEditNode(arg *wrapper_admin.ArgTreeEditNode) (res *wrapper_admin.ResultTreeEditNode, err error) {

	res = &wrapper_admin.ResultTreeEditNode{}

	if arg.Id == 0 {
		var data = &models.HelpDocumentRelate{
			Id:         arg.Id,
			Display:    arg.Display,
			ParentId:   arg.ParentId,
			IsLeafNode: arg.IsLeafNode,
			DocKey:     arg.DocKey,
			Label:      arg.Label,
			BizCode:    arg.BizCode,
			CreatedAt:  arg.TimeNow,
			UpdatedAt:  arg.TimeNow,
		}
		if err = r.dao.AddOneHelpRelate(data); err != nil {
			return
		}
		res.Result = true
		return
	}
	var data = map[string]interface{}{
		"biz_code":     arg.BizCode,
		"display":      arg.Display,
		"parent_id":    arg.ParentId,
		"label":        arg.Label,
		"is_leaf_node": arg.IsLeafNode,
		"doc_key":      arg.DocKey,
		"updated_at":   arg.TimeNow.Format(utils.DateTimeGeneral),
	}
	if err = r.dao.UpdateById(arg.Id, data); err != nil {
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
