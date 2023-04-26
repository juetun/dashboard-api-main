package srv_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type SrvHelpImpl struct {
	base.ServiceBase
	dao daos.DaoHelp
}

func (r *SrvHelpImpl) HelpList(arg *wrapper_admin.ArgHelpList) (res *wrapper_admin.ResultHelpList, err error) {
	res = &wrapper_admin.ResultHelpList{Pager: response.NewPager(response.PagerBaseQuery(&arg.PageQuery))}
	switch arg.PageType {
	case response.DefaultPageTypeList:
		err = r.getWithList(arg, res)
	default:
		err = fmt.Errorf("当前不支持你选择的分类类型")
	}
	return
}

func (r *SrvHelpImpl) getWithList(arg *wrapper_admin.ArgHelpList, res *wrapper_admin.ResultHelpList) (err error) {
	var actRes *base.ActErrorHandlerResult
	var actGetCount = func(pagerObject *response.Pager) (err error) {
		actRes, err = r.dao.GetCountByArg(arg, pagerObject)
		return
	}

	var actGetList = func(pagerObject *response.Pager) (err error) {
		var list []*models.HelpDocument
		if list, err = r.dao.GetListByArg(arg, actRes, pagerObject); err != nil {
			return
		}

		pagerObject.List = r.orgHelpList(list, arg)
		return
	}
	err = res.CallGetPagerData(actGetCount, actGetList)
	return
}

func (r *SrvHelpImpl) orgHelpList(documents []*models.HelpDocument, arg *wrapper_admin.ArgHelpList) (res []*wrapper_admin.ResultHelpListItem) {
	res = make([]*wrapper_admin.ResultHelpListItem, 0, len(documents))
	var dt *wrapper_admin.ResultHelpListItem
	for _, item := range documents {
		dt = &wrapper_admin.ResultHelpListItem{}
		dt.SetHelpDocument(item, arg.TimeNow)
		res = append(res, dt)
	}
	return
}

func (r *SrvHelpImpl) HelpDetail(arg *wrapper_admin.ArgHelpDetail) (res *wrapper_admin.ResultHelpDetail, err error) {
	res = &wrapper_admin.ResultHelpDetail{}

	var (
		resHelpDocumentRelate map[string]*models.HelpDocumentRelate
	)
	if resHelpDocumentRelate, err = dao_impl.NewDaoHelpRelate(r.Context).
		GetByDocKeys(base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(arg.PKey))); err != nil {
		return
	}
	if helpDocumentRelate, ok := resHelpDocumentRelate[arg.PKey]; !ok {
		err = fmt.Errorf("您要编辑的帮助文档关系信息不存在,或已删除")
		return
	} else {
		res.Label = helpDocumentRelate.Label
	}

	var help *models.HelpDocument
	if arg.Id > 0 {
		if help, err = r.getById(arg.Id); err != nil {
			return
		}
	} else if arg.PKey != "" {
		if help, err = r.getByKey(arg.PKey); err != nil {
			return
		}
	}
	res.ParseFromHelpDoc(help)
	return
}

func (r *SrvHelpImpl) getByKey(PKey string) (res *models.HelpDocument, err error) {
	var helpMap map[string]*models.HelpDocument
	if helpMap, err = r.dao.GetByPKey(base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(PKey))); err != nil {
		return
	}
	var ok bool
	if res, ok = helpMap[PKey]; !ok {
		res = &models.HelpDocument{
			PKey: PKey,
		}
		res.Default()
		return
	}
	return
}

func (r *SrvHelpImpl) getById(id int64) (res *models.HelpDocument, err error) {
	var helpMap map[int64]*models.HelpDocument
	var argNumber = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(id))
	if helpMap, err = r.dao.GetByIds(argNumber); err != nil {
		return
	}
	var ok bool
	if res, ok = helpMap[id]; !ok {
		err = fmt.Errorf("你要查看(或编辑)的帮助信息不存在或已删除")
		return
	}
	return
}

func (r *SrvHelpImpl) validatePk(arg *wrapper_admin.ArgHelpEdit, helpMap map[string]*models.HelpDocument) (help *models.HelpDocument, err error) {
	var ok bool

	if help, ok = helpMap[arg.PKey]; !ok { //如果没查到数据
		return
	}
	if arg.Id == 0 {
		err = fmt.Errorf("你输入的数据唯一KEY已存在")
		return
	}
	if help.Id != arg.Id {
		err = fmt.Errorf("你输入的数据唯一KEY已存在")
		return
	}
	return
}

func (r *SrvHelpImpl) HelpEdit(arg *wrapper_admin.ArgHelpEdit) (res *wrapper_admin.ResultHelpEdit, err error) {
	res = &wrapper_admin.ResultHelpEdit{}
	var (
		help map[string]*models.HelpDocument
	)
	if help, err = r.dao.GetByPKey(base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(arg.PKey))); err != nil {
		return
	}

	//判断key的唯一性 (去重)
	if _, err = r.validatePk(arg, help); err != nil {
		return
	}

	if arg.Id != 0 {
		if _, err = r.getById(arg.Id); err != nil {
			return
		}
	}

	data := &models.HelpDocument{
		Id:        arg.Id,
		PKey:      arg.PKey,
		Label:     arg.Label,
		Desc:      arg.Desc,
		Content:   arg.Content,
		CreatedAt: arg.TimeNow,
	}
	if err = data.DefaultBeforeAdd(); err != nil {
		return
	}
	if err = r.dao.AddHelpDocument(data); err != nil {
		return
	}
	res.Result = true
	return
}

func NewSrvHelp(context ...*base.Context) (res srvs.SrvHelp) {
	p := &SrvHelpImpl{}
	p.SetContext(context...)
	p.dao = dao_impl.NewDaoHelp(p.Context)
	return p
}
