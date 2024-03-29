package srv_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_outernet"
	"github.com/juetun/library/common/app_param/upload_operate"
	"strings"
)

type SrvHelpRelateImpl struct {
	base.ServiceBase
	dao     daos.DaoHelpRelate
	daoHelp daos.DaoHelp
}

func (r *SrvHelpRelateImpl) Data(arg *wrapper_outernet.ArgData) (res *wrapper_outernet.ResultData, err error) {
	res = &wrapper_outernet.ResultData{}
	var (
		ok           bool
		resData      map[string]*models.HelpDocument
		helpDocument *models.HelpDocument
	)
	if resData, err = r.daoHelp.
		GetByPKey(base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(arg.DocKey))); err != nil {
		return
	}
	if helpDocument, ok = resData[arg.DocKey]; !ok {
		res.HavError = true
		res.ErrorMsg = "帮助信息不存在或已移除"
		return
	}
	res.DocName = helpDocument.Label
	if res.DocContent, err = r.GetDescMedia(helpDocument.Content, &arg.GetDataTypeCommon); err != nil {
		return
	}

	return
}

func (r *SrvHelpRelateImpl) getMapParentAndList(arg *wrapper_outernet.ArgTree) (haveDocKey string, mapParent map[int64][]*wrapper_outernet.ResultFormPage, dataList []*wrapper_outernet.ResultFormPage, err error) {
	dataList = []*wrapper_outernet.ResultFormPage{}
	var (
		dataRes   models.HelpDocumentRelateCaches
		argString = base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionGetDataTypeCommon(arg.GetDataTypeCommon))
	)

	if dataRes, err = r.dao.GetAllHelpRelate(argString); err != nil {
		return
	}
	var (
		l        = len(dataRes)
		dataItem *wrapper_outernet.ResultFormPage
	)
	mapParent = make(map[int64][]*wrapper_outernet.ResultFormPage, l)
	dataList = make([]*wrapper_outernet.ResultFormPage, 0, len(dataRes))
	for _, item := range dataRes {
		dataItem = wrapper_outernet.NewResultFormPage()
		dataItem.ParseFromHelpDocumentRelateCache(item)
		if item.DocKey == arg.DocKey && arg.DocKey != "" {
			haveDocKey = item.DocKey
			dataItem.IsActive = true
		}
		if item.ParentId == 0 {
			dataList = append(dataList, dataItem)
			continue
		}
		if _, ok := mapParent[item.ParentId]; !ok {
			mapParent[item.ParentId] = make([]*wrapper_outernet.ResultFormPage, 0, l)
		}
		mapParent[item.ParentId] = append(mapParent[item.ParentId], dataItem)
	}
	return
}

func (r *SrvHelpRelateImpl) orgResultHelpTreeItem(arg *wrapper_outernet.ArgTree, res *wrapper_outernet.ResultTree) (err error) {
	var (
		dataList []*wrapper_outernet.ResultFormPage

		mapParent       map[int64][]*wrapper_outernet.ResultFormPage
		docKey          string
		haveSelectIdMap = make(map[int64]bool, 20)
		resData         map[string]*models.HelpDocument
		ok              bool
		helpDocument    *models.HelpDocument
	)

	if docKey, mapParent, dataList, err = r.getMapParentAndList(arg); err != nil {
		return
	}
	//如果没有设置默认的
	r.orgHelpChild("", dataList, mapParent)
	r.defaultDataList(&haveSelectIdMap, docKey, dataList)

	//默认界面展示
	r.defaultPageHelpTreeActive(dataList, &res.Breadcrumbs)
	if resData, err = r.daoHelp.
		GetByPKey(base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(docKey))); err != nil {
		return
	}
	if helpDocument, ok = resData[docKey]; !ok {
		res.NotExists = true
		res.ErrorMsg = "帮助信息不存在或已移除"
	} else {
		if res.DocContent, err = r.GetDescMedia(helpDocument.Content, &arg.GetDataTypeCommon); err != nil {
			return
		}
	}
	res.Menu = dataList
	return
}

type (
	MapMediaDetailReplace struct { //替换
		Video map[string]string `json:"video"`
		Img   map[string]string `json:"img"`
		//File  map[string]string `json:"file"`
	}
)

func (r *SrvHelpRelateImpl) GetDescMedia(desc string, argCommon *base.GetDataTypeCommon) (resDesc string, err error) {
	var (
		keysDescImg, keysDescVideo map[string]string
		argUploadGetInfo           = &upload_operate.ArgUploadGetInfo{}
		mapKeysReplace             = &MapMediaDetailReplace{
			Video: make(map[string]string, 100),
			Img:   make(map[string]string, 100),
		}
		resUpload *upload_operate.ResultMapUploadInfo
	)
	if keysDescImg, keysDescVideo, err = upload_operate.ParseTextEditorContent(desc); err != nil {
		return
	}
	r.collectFileAndVideo(mapKeysReplace, argUploadGetInfo, keysDescImg, keysDescVideo)
	if resUpload, err = upload_operate.NewDaoUpload(r.Context).
		GetUploadByKeys(argUploadGetInfo, argCommon); err != nil {
		return
	}
	resDesc = r.description(resUpload, desc, mapKeysReplace)

	return
}

func (r *SrvHelpRelateImpl) description(uploadInfo *upload_operate.ResultMapUploadInfo, description string, mapKeysReplace *MapMediaDetailReplace, ) (res string) {
	var (
		uploadImage *upload_operate.UploadFile
		uploadVideo *upload_operate.UploadVideo
		ok          bool
		html        string
	)
	res = description
	for key, value := range mapKeysReplace.Img {
		if uploadImage, ok = uploadInfo.File[key]; !ok { //图片地址替换
			continue
		}
		html, _ = uploadImage.GetEditorHtml(value)
		//html = strings.ReplaceAll(html, "%", "%%")
		description = strings.Replace(description, value, html, -1)
	}
	for key, value := range mapKeysReplace.Video {
		if uploadVideo, ok = uploadInfo.Video[key]; ok { //视频文件替换
			continue
		}
		html, _ = uploadVideo.GetEditorHtml(key)
		description = strings.Replace(description, value, html, -1)
	}
	res = description
	return
}

func (r *SrvHelpRelateImpl) collectFileAndVideo(mapKeysReplace *MapMediaDetailReplace, argUploadGetInfo *upload_operate.ArgUploadGetInfo, keysDescImg, keysDescVideo map[string]string) {

	for key, value := range keysDescImg {
		if _, ok := mapKeysReplace.Img[key]; !ok {
			mapKeysReplace.Img[key] = value
			argUploadGetInfo.File = append(argUploadGetInfo.File, key)
		}
	}
	for key, value := range keysDescVideo {
		if _, ok := mapKeysReplace.Video[key]; !ok {
			mapKeysReplace.Video[key] = value
			argUploadGetInfo.VideoKeys = append(argUploadGetInfo.VideoKeys, key)
		}
	}
	return
}
func (r *SrvHelpRelateImpl) defaultDataList(haveSelectIdMap *map[int64]bool, haveDocKey string, dataList []*wrapper_outernet.ResultFormPage) (parentIsActive bool) {
	for key, item := range dataList {
		if haveDocKey == "" { //如果默认选定
			if key == 0 {
				item.IsActive = true
				(*haveSelectIdMap)[item.Id] = true
				if len(item.Children) > 0 {
					r.defaultDataList(haveSelectIdMap, haveDocKey, item.Children)
				}
			}
			continue
		}
		if len(item.Children) > 0 {
			item.IsActive = r.defaultDataList(haveSelectIdMap, haveDocKey, item.Children)
		}
		//如果手工选定了
		if item.IsActive {
			parentIsActive = true
			(*haveSelectIdMap)[item.Id] = true
		}
	}
	return
}

func (r *SrvHelpRelateImpl) defaultPageHelpTreeActive(dataList []*wrapper_outernet.ResultFormPage, Breadcrumbs *[]wrapper_outernet.BreadcrumbItem) {
	for _, item := range dataList {
		item.OpenNames = make([]string, 0, 20)
		if item.IsActive {
			*Breadcrumbs = append(*Breadcrumbs, wrapper_outernet.BreadcrumbItem{
				Label: item.Label,
			})
			r.getOpenNames(item, &item.OpenNames, &item.Active.ActiveName, Breadcrumbs)
		}
	}
	return
}

func (r *SrvHelpRelateImpl) getOpenNames(dataItem *wrapper_outernet.ResultFormPage, OpenNames *[]string, ActiveName *string, Breadcrumbs *[]wrapper_outernet.BreadcrumbItem) {
	for _, item := range dataItem.Children {
		if item.IsActive {
			if item.IsLeafNode == models.HelpDocumentRelateIsLeafNodeYes { //如果是叶子节点
				*ActiveName = item.TreeName
			} else {
				*OpenNames = append(*OpenNames, item.TreeName)
			}
			*Breadcrumbs = append(*Breadcrumbs, wrapper_outernet.BreadcrumbItem{
				Label: item.Label,
			})
			if len(item.Children) > 0 {
				r.getOpenNames(item, OpenNames, ActiveName, Breadcrumbs)
			}
		}
	}
	return
}

func (r *SrvHelpRelateImpl) orgHelpChild(treeName string, dataList []*wrapper_outernet.ResultFormPage, mapParent map[int64][]*wrapper_outernet.ResultFormPage) {
	for _, item := range dataList {
		if treeName != "" {
			item.TreeName = fmt.Sprintf("%v-%v", treeName, item.Name)
		} else {
			item.TreeName = item.Name
		}

		if child, ok := mapParent[item.Id]; ok {
			r.orgHelpChild(item.TreeName, child, mapParent)
			item.Children = child
		}
	}
	return
}
func (r *SrvHelpRelateImpl) Tree(arg *wrapper_outernet.ArgTree) (res *wrapper_outernet.ResultTree, err error) {

	res = wrapper_outernet.NewResultTree()

	if err = r.orgResultHelpTreeItem(arg, res); err != nil {
		return
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
	//var haveData bool

	for _, item := range dataRes {
		if item.Id == arg.CurrentId {
			dataRes = item.Child
			//haveData = true
			break
		}
	}
	//if !haveData {
	//	return
	//}
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
				//arg.TopId = item.Id
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
			//arg.TopId = item.Id
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

func (r *SrvHelpRelateImpl) validateTreeEditNode(arg *wrapper_admin.ArgTreeEditNode) (err error) {

	//验证KEY是否唯一
	var (
		data    *models.HelpDocumentRelate
		dataMap map[string]*models.HelpDocumentRelate
		ok      bool
	)
	var argString = base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(arg.DocKey))
	if dataMap, err = r.dao.GetByDocKeys(argString); err != nil {
		return
	} else if data, ok = dataMap[arg.DocKey]; !ok {
		return
	}

	if data == nil || data.Id == 0 {
		return
	}
	if arg.Id > 0 && arg.Id == data.Id {
		return
	}
	err = fmt.Errorf("KEY(%s)已被(ID:%d)使用", arg.DocKey, data.Id)
	return
}

func (r *SrvHelpRelateImpl) validateParentAllParentId(arg *wrapper_admin.ArgTreeEditNode) (err error) {
	var parentIdMap = make(map[int64]bool, 20)
	if arg.ParentId != 0 {
		if err = r.getRelateById(arg.ParentId, &parentIdMap); err != nil {
			return
		}
		//如果当前参数ID和当前数据的所有父级关系ID存在交集，则直接报参数异常
		if _, ok := parentIdMap[arg.Id]; ok {
			err = fmt.Errorf("您传递的参数存在cricle。参数异常")
			return
		}
	}
	return
}

func (r *SrvHelpRelateImpl) getRelateById(id int64, parentIdMap *map[int64]bool) (err error) {
	var (
		relates []*models.HelpDocumentRelate
		relate  *models.HelpDocumentRelate
	)
	(*parentIdMap)[id] = true
	if relates, err = r.dao.GetByIdFromDb([]int64{id}...); err != nil {
		return
	} else if len(relates) > 0 {
		relate = relates[0]
	}

	if relate != nil && relate.Id != 0 && relate.ParentId != 0 {
		if err = r.getRelateById(relate.ParentId, parentIdMap); err != nil {
			return
		}
	}
	return
}
func (r *SrvHelpRelateImpl) TreeEditNode(arg *wrapper_admin.ArgTreeEditNode) (res *wrapper_admin.ResultTreeEditNode, err error) {

	res = &wrapper_admin.ResultTreeEditNode{}

	if err = r.validateTreeEditNode(arg); err != nil {
		return
	}

	defer func() {

		if res.Result { //刷新缓存
			go func() {
				//刷新关系缓存
				var arg = base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionRefreshCache(base.RefreshCacheYes))
				_, _ = r.dao.GetAllHelpRelate(arg)
			}()
		}
	}()

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

	//验证当前菜单的所有
	if err = r.validateParentAllParentId(arg); err != nil {
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
	p.daoHelp = dao_impl.NewDaoHelp(p.Context)
	return p
}
