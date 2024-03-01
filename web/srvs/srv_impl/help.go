package srv_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"github.com/juetun/library/common/app_param/upload_operate"
	"github.com/juetun/library/common/recommend"
	"strings"
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
	res = &wrapper_admin.ResultHelpDetail{UploadDataType: recommend.AdDataDataTypeHelpDocument}
	defer func() {
		res.UploadDataId = fmt.Sprintf("%v", res.Id)
	}()
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
		if helpDocumentRelate == nil {
			helpDocumentRelate = &models.HelpDocumentRelate{}
		}
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

	if res.Content, err = NewSrvHelpRelate(r.Context).
		GetDescMedia(help.Content, &arg.GetDataTypeCommon); err != nil {
		return
	}

	res.ParseFromHelpDoc(help)

	return
}

func (r *SrvHelpImpl) getByKey(PKey string) (res *models.HelpDocument, err error) {
	var helpMap map[string]*models.HelpDocument
	if helpMap, err = r.dao.GetByPKey(base.NewArgGetByStringIds(
		base.ArgGetByStringIdsOptionIds(PKey),
	)); err != nil {
		return
	}
	var ok bool
	if res, ok = helpMap[PKey]; !ok {
		res = &models.HelpDocument{
			PKey: PKey,
		}
		err = res.Default()
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
	go func(helpDoc *models.HelpDocument) {
		_, _ = dao_impl.NewDaoHelpRelate(r.Context).
			GetByDocKeys(base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(helpDoc.PKey), base.ArgGetByStringIdsOptionRefreshCache(base.RefreshCacheYes)))
		_, _ = r.dao.GetByIds(base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(helpDoc.Id), base.ArgGetByNumberIdsOptionRefreshCache(base.RefreshCacheYes)));
	}(data)
	res.Result = true
	return
}

func (r *SrvHelpImpl) description(uploadInfo *upload_operate.ResultMapUploadInfo, helpDocument *models.HelpDocument, mapKeysReplace *wrappers.MapDetailReplace, ) {
	var (
		uploadImage *upload_operate.UploadFile
		uploadVideo *upload_operate.UploadVideo
		ok          bool
		html        string
	)

	for key, value := range mapKeysReplace.Img {
		if uploadImage, ok = uploadInfo.File[key]; !ok { //图片地址替换
			continue
		}
		html, _ = uploadImage.GetEditorHtml(value)
		//html = strings.ReplaceAll(html, "%", "%%")
		helpDocument.Content = strings.Replace(helpDocument.Content, value, html, -1)
	}
	for key, value := range mapKeysReplace.Video {
		if uploadVideo, ok = uploadInfo.Video[key]; ok { //视频文件替换
			continue
		}
		html, _ = uploadVideo.GetEditorHtml(key)
		helpDocument.Content = strings.Replace(helpDocument.Content, value, html, -1)
	}
	return
}
func NewSrvHelp(context ...*base.Context) (res srvs.SrvHelp) {
	p := &SrvHelpImpl{}
	p.SetContext(context...)
	p.dao = dao_impl.NewDaoHelp(p.Context)
	return p
}
