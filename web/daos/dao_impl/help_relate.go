package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type DaoHelpRelateImpl struct {
	base.ServiceDao
}

func (r *DaoHelpRelateImpl) GetByTopId(bizCode string, topIds ...int64) (res map[int64][]*wrapper_admin.ResultHelpTreeItem, err error) {
	if len(topIds) == 0 || bizCode == "" {
		return
	}
	var list []*models.HelpDocumentRelate
	var m *models.HelpDocumentRelate
	var actResParam = r.GetDefaultActErrorHandlerResult(m)

	var actHandler = func() (actRes *base.ActErrorHandlerResult) {
		actRes = actResParam
		actRes.Err = actRes.Db.
			Table(actRes.TableName).
			Scopes(base.ScopesDeletedAt()).
			Where("id IN (?) AND biz_code=?", topIds, bizCode).
			Find(&list).
			Error
		if actRes.Err != nil {
			return
		}
		return
	}
	if err = r.ActErrorHandler(actHandler); err != nil {
		return
	}
	if res, err = r.orgGroupRelate(bizCode, actResParam, topIds...); err != nil {
		return
	}
	return
}

func (r *DaoHelpRelateImpl) getChildByTopIds(bizCode string, actRes *base.ActErrorHandlerResult, topIds ...int64) (dataMap map[int64][]*wrapper_admin.ResultHelpTreeItem, childTopIds []int64, err error) {
	if len(topIds) == 0 {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"topIds": topIds,
			"err":    err.Error(),
		}, "DaoHelpRelateImplGetChildByTopIds")
	}()
	var list []*models.HelpDocumentRelate
	if err = actRes.Db.Table(actRes.TableName).
		Where("parent_id IN(?) AND biz_code = ?", topIds, bizCode).
		Scopes(base.ScopesDeletedAt()).
		Find(&list).
		Error; err != nil {
		return
	}
	dataMap, childTopIds = r.groupRelate(list)
	return
}

func (r *DaoHelpRelateImpl) orgGroupRelate(bizCode string, actRes *base.ActErrorHandlerResult, topIds ...int64) (res map[int64][]*wrapper_admin.ResultHelpTreeItem, err error) {
	var (
		childTopIds []int64
	)
	if res, childTopIds, err = r.getChildByTopIds(bizCode, actRes, topIds...); err != nil {
		return
	}
	if len(childTopIds) > 0 {
		var resChildren map[int64][]*wrapper_admin.ResultHelpTreeItem
		//获得儿子的数据
		if resChildren, err = r.orgGroupRelate(bizCode, actRes, childTopIds...); err != nil {
			return
		}
		for _, item := range res {
			for _, it := range item {
				if dt, ok := resChildren[it.Id]; ok {
					it.Child = dt
				}
			}
		}
		return
	}

	return
}

//组织数据使用
func (r *DaoHelpRelateImpl) groupRelate(list []*models.HelpDocumentRelate) (dataMap map[int64][]*wrapper_admin.ResultHelpTreeItem, childTopIds []int64) {
	dataMap = map[int64][]*wrapper_admin.ResultHelpTreeItem{}
	var (
		l     = len(list)
		dtRes *wrapper_admin.ResultHelpTreeItem
	)
	for _, item := range list {
		if _, ok := dataMap[item.ParentId]; !ok {
			dataMap[item.ParentId] = make([]*wrapper_admin.ResultHelpTreeItem, 0, l)
		}
		dtRes = &wrapper_admin.ResultHelpTreeItem{}
		dtRes.HelpDocumentRelate = *item
		//如果不是叶子节点
		if item.IsLeafNode == models.HelpDocumentRelateIsLeafNodeNo && item.Id != 0 {
			childTopIds = append(childTopIds, item.Id)
		}
		dataMap[item.ParentId] = append(dataMap[item.ParentId], dtRes)
	}
	return
}

func (r *DaoHelpRelateImpl) AddOneHelpRelate(relate *models.HelpDocumentRelate) (err error) {
	if relate == nil {
		return
	}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"relate": relate,
			"err":    err.Error(),
		}, "DaoHelpRelateImplAddOneHelpRelate")
		err = base.NewErrorRuntime(err, base.ErrorSqlCode)
	}()
	err = r.ActErrorHandler(func() (actRes *base.ActErrorHandlerResult) {
		actRes = r.GetDefaultActErrorHandlerResult(relate)
		actRes.Err = r.AddOneData(actRes.ParseAddOneDataParameter(base.AddOneDataParameterModel(relate)))
		return
	})
	return
}

func NewDaoHelpRelate(c ...*base.Context) daos.DaoHelpRelate {
	p := &DaoHelpRelateImpl{}
	p.SetContext(c...)
	return p
}