/**
* @Author:changjiang
* @Description:
* @File:permit_import
* @Version: 1.0.0
* @Date 2021/9/12 11:37 上午
 */
package srv_impl

import (
	"fmt"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type SrvPermitImport struct {
	base.ServiceBase
}

func (r *SrvPermitImport) UpdateImportValue(arg *wrappers.ArgUpdateImportValue) (res *wrappers.ResultUpdateImportValue, err error) {
	res = &wrappers.ResultUpdateImportValue{}
	var condition = fmt.Sprintf("id IN (%s)", strings.Join(arg.Ids, ","))
	dao := dao_impl.NewPermitImportImpl(r.Context)
	var data = make(map[string]interface{}, 1)
	t := time.Now().Format("2006-01-02 15:04:05")
	data["updated_at"] = t
	switch arg.Column {
	case "is_del":
		data[arg.Column] = t
		if arg.Val == "1" {
			res.Result, err = dao.DeleteByCondition(condition)
		} else {
			res.Result, err = dao.UpdateByCondition(condition, data)
		}

	case "need_login", "default_open", "need_sign":
		if arg.Val != "1" && arg.Val != "2" {
			err = fmt.Errorf("您设置的值不正确")
			return
		}
		data[arg.Column] = arg.Val
		res.Result, err = dao.UpdateByCondition(condition, data)
	case "":
	default:
		err = fmt.Errorf("您选择的数据值(column:%s)不正确", arg.Column)
		return
	}
	if err != nil {
		return
	}
	return
}

func (r *SrvPermitImport) GetImportByMenuId(arg *wrappers.ArgGetImportByMenuId) (res wrappers.ResultGetImportByMenuId, err error) {
	res = wrappers.ResultGetImportByMenuId{
		ImportIds: []int{},
		MenuIds:   []int{},
	}

	permitMenu := NewSrvPermitMenuImpl(r.Context)

	if err = permitMenu.GetMenuPermitKeyByPath(
		&arg.ArgGetImportByMenuIdSingle,
		dao_impl.NewDaoPermit(r.Context),
	); err != nil {
		return
	}

	if arg.NowMenuId == 0 {
		return
	}

	if res.MenuIds, err = r.GetChildMenu(arg.NowMenuId); err != nil {
		return
	}

	if res.ImportIds, err = r.GetChildImport(arg.NowMenuId); err != nil {
		return
	}
	return
}

func (r *SrvPermitImport) GetChildImport(nowMenuId int) (importIds []int, err error) {
	importIds = []int{}
	dao := dao_impl.NewPermitImportImpl(r.Context)
	var importList []models.AdminMenuImport
	if importList, err = dao.GetChildImportByMenuId(nowMenuId); err != nil {
		return
	}
	importIds = make([]int, 0, len(importList))
	for _, value := range importList {
		importIds = append(importIds, value.ImportId)
	}
	return
}

func (r *SrvPermitImport) GetOpList(dao daos.DaoPermit, arg *wrappers.ArgPermitMenu) (opList map[string][]wrappers.OpOne, err error) {
	var list []wrappers.Op
	if list, err = dao.GetPermitImportByModule(arg); err != nil {
		return
	}
	l := len(list)
	opList = make(map[string][]wrappers.OpOne, l)
	var t wrappers.OpOne
	for _, value := range list {
		if _, ok := opList[value.MenuPermitKey]; !ok {
			opList[value.MenuPermitKey] = make([]wrappers.OpOne, 0, l)
		}
		t = wrappers.OpOne(value.PermitKey)
		opList[value.MenuPermitKey] = append(opList[value.MenuPermitKey], t)
	}
	return
}
func (r *SrvPermitImport) GetChildMenu(nowMenuId int) (menuIds []int, err error) {
	menuIds = []int{}
	dao := dao_impl.NewDaoPermit(r.Context)
	var res []models.AdminMenu
	if res, err = dao.GetAdminMenuList(&wrappers.ArgAdminMenu{
		ParentId: nowMenuId,
	}); err != nil {
		return
	}
	for _, item := range res {
		menuIds = append(menuIds, item.Id)
	}
	return
}
func (r *SrvPermitImport) GetImport(arg *wrappers.ArgGetImport) (res *wrappers.ResultGetImport, err error) {
	res = &wrappers.ResultGetImport{
		Pager: response.NewPager(response.PagerBaseQuery(arg.PageQuery)),
	}
	dao := dao_impl.NewDaoPermit(r.Context)
	var db *gorm.DB
	if db, err = dao.GetImportCount(arg, &res.TotalCount); err != nil {
		return
	}
	if res.TotalCount == 0 {
		return
	}
	var list []models.AdminImport
	if list, err = dao.GetImportList(db, arg); err != nil {
		return
	}
	if !arg.Checked {
		res.List = list
		return
	}
	res.List, err = r.joinChecked(dao, arg, list)

	// []models.AdminImport{}
	return
}

func (r *SrvPermitImport) SetApiPermit(dao daos.DaoPermit, arg *wrappers.ArgAdminSetPermit) (err error) {
	switch arg.Act {
	case models.SetPermitAdd:
		var dt models.AdminUserGroupPermit
		var list []models.AdminUserGroupPermit
		var t = time.Now()
		for _, pid := range arg.PermitIds {
			dt = models.AdminUserGroupPermit{
				GroupId:   arg.GroupId,
				MenuId:    pid,
				PathType:  models.PathTypeApi,
				CreatedAt: t,
				UpdatedAt: t,
				DeletedAt: nil,
			}
			list = append(list, dt)
		}
		var m models.AdminUserGroupPermit
		if err = dao.BatchGroupPermit(m.TableName(), list); err != nil {
			err = fmt.Errorf("操作异常")
			return
		}

	case models.SetPermitCancel:
		if err = dao.DeleteGroupPermit(arg.GroupId, models.PathTypeApi, arg.PermitIds...); err != nil {
			return
		}
	default:
		err = fmt.Errorf("act格式错误")
		return
	}
	return
}
func (r *SrvPermitImport) getImportId(l int, list []models.AdminImport) (importId []int) {
	importId = make([]int, 0, l)
	for _, value := range list {
		importId = append(importId, value.Id)
	}
	return
}
func (r *SrvPermitImport) ImportList(arg *wrappers.ArgImportList) (res *wrappers.ResultImportList, err error) {
	var db *gorm.DB
	if arg.Order == "" {
		arg.Order = "id desc"
	}
	res = &wrappers.ResultImportList{Pager: response.NewPagerAndDefault(&arg.PageQuery)}
	dao := dao_impl.NewDaoPermit(r.Context)

	// 获取分页数据
	if err = res.Pager.CallGetPagerData(func(pager *response.Pager) (err error) {
		pager.TotalCount, db, err = dao.GetImportListCount(db, arg)
		return
	}, func(pager *response.Pager) (err error) {
		var list []models.AdminImport
		list, err = dao.GetImportListData(db, arg, pager)
		pager.List, err = r.orgImportList(dao, list)
		return
	}); err != nil {
		return
	}
	return
}

func (r *SrvPermitImport) orgImportList(dao daos.DaoPermit, list []models.AdminImport) (res []wrappers.AdminImportList, err error) {
	l := len(list)
	res = make([]wrappers.AdminImportList, 0, l)
	var dt wrappers.AdminImportList
	var dta map[int][]wrappers.AdminImportListMenu
	if dta, err = r.getImportMenuGroup(dao, l, list); err != nil {
		return
	}
	for _, value := range list {
		dt = wrappers.AdminImportList{AdminImport: value, Menu: []wrappers.AdminImportListMenu{}}
		dt.RequestMethods = value.GetRequestMethods()
		if _, ok := dta[value.Id]; ok {
			dt.Menu = dta[value.Id]
		}
		res = append(res, dt)
	}
	return
}

func (r *SrvPermitImport) EditImport(arg *wrappers.ArgEditImport) (res *wrappers.ResultEditImport, err error) {
	res = &wrappers.ResultEditImport{Result: false}
	var (
		dao        = dao_impl.NewDaoPermit(r.Context)
		daoService = dao_impl.NewDaoServiceImpl(r.Context)
		listImport []models.AdminImport
		apps       []models.AdminApp
	)

	if apps, err = daoService.GetByKeys(strings.TrimSpace(arg.AppName)); err != nil {
		return
	}
	if len(apps) == 0 {
		err = fmt.Errorf("您输入的应用(%s)不存在或已删除", arg.AppName)
		return
	}

	if listImport, err = dao.GetImportByCondition(map[string]interface{}{"app_name": arg.AppName, "url_path": arg.UrlPath}); err != nil {
		return
	}

	// 验证数据是否重复
	for _, value := range listImport {
		if err = r.editImportParam(arg, &value); err != nil {
			return
		}
	}
	var mAi models.AdminImport
	mAi.SetRequestMethods(arg.RequestMethod)
	data := map[string]interface{}{
		`app_name`:       arg.AppName,
		`app_version`:    arg.AppVersion,
		`url_path`:       arg.UrlPath,
		`request_method`: mAi.RequestMethod,
		`sort_value`:     arg.SortValue,
		`updated_at`:     arg.RequestTime,
		"need_login":     arg.NeedLogin,
		"need_sign":      arg.NeedSign,
	}
	if arg.Id == 0 { // 如果是添加接口
		res.Result, err = r.createImport(dao, arg)
		return
	}

	var m = models.AdminImport{Id: arg.Id}
	var dt []models.AdminImport
	if dt, err = dao.GetAdminImportById(arg.Id); err != nil {
		return
	}

	if len(dt) == 0 {
		err = fmt.Errorf("您编辑的接口信息不存在或已删除")
		return
	}

	if dt[0].PermitKey == "" {
		data["permit_key"] = m.GetPathName()
	}
	// 如果更新了app_name
	if dt[0].AppName != arg.AppName {
		if err = dao_impl.NewPermitImportImpl(r.Context).
			UpdateMenuImport(fmt.Sprintf("import_id = %d ", dt[0].Id),
				map[string]interface{}{"import_app_name": arg.AppName}); err != nil {
			return
		}

	}
	if _, err = dao.UpdateAdminImport(map[string]interface{}{"id": arg.Id}, data); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitImport) createImport(dao daos.DaoPermit, arg *wrappers.ArgEditImport) (res bool, err error) {
	t := time.Now()
	data := models.AdminImport{
		AppName:    arg.AppName,
		AppVersion: arg.AppVersion,
		UrlPath:    arg.UrlPath,
		SortValue:  arg.SortValue,
		UpdatedAt:  t,
		CreatedAt:  t,
	}
	data.SetRequestMethods(arg.RequestMethod)
	if _, err = dao.CreateImport(&data); err != nil {
		return
	}
	res = true
	return
}

func (r *SrvPermitImport) editImportParam(arg *wrappers.ArgEditImport, value *models.AdminImport) (err error) {
	if arg.Id == 0 {
		if value.AppName == arg.AppName && value.UrlPath == arg.UrlPath {
			err = fmt.Errorf("您输入的接口信息已存在")
			return
		}
		return
	}
	if value.AppName == arg.AppName && value.UrlPath == arg.UrlPath {
		if arg.Id != value.Id {
			err = fmt.Errorf("您输入的接口信息已存在")
			return
		}
	}

	return
}

func (r *SrvPermitImport) DeleteImport(arg *wrappers.ArgDeleteImport) (res *wrappers.ResultDeleteImport, err error) {
	res = &wrappers.ResultDeleteImport{}
	dao := dao_impl.NewPermitImportImpl(r.Context)
	if err = dao.DeleteImportByIds([]int{arg.ID}...); err != nil {
		return
	}
	daoGroup := dao_impl.NewDaoPermitGroupImpl(r.Context)
	if err = daoGroup.DeleteUserGroupPermit(models.PathTypeApi, []int{arg.ID}...); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitImport) joinChecked(dao daos.DaoPermit, arg *wrappers.ArgGetImport, data []models.AdminImport) (res []wrappers.AdminImport, err error) {
	res = make([]wrappers.AdminImport, 0, len(data))
	var dt wrappers.AdminImport
	var importId = make([]int, 0, len(data))
	for _, value := range data {
		importId = append(importId, value.Id)
	}
	var li []models.AdminUserGroupPermit
	if li, err = dao.GetSelectImportByImportId(arg.GroupId, importId...); err != nil {
		return
	}
	var m = make(map[int]int, len(li))
	for _, it := range li {
		m[it.MenuId] = it.MenuId
	}
	for _, value := range data {
		dt = wrappers.AdminImport{
			AdminImport: value,
		}
		if _, ok := m[value.Id]; ok {
			dt.Checked = true
		}
		res = append(res, dt)
	}
	return
}

func (r *SrvPermitImport) getImportMenuGroup(dao daos.DaoPermit, l int, data []models.AdminImport) (res map[int][]wrappers.AdminImportListMenu, err error) {
	importId := r.getImportId(l, data)
	daoImportMenu := dao_impl.NewPermitImportImpl(r.Context)

	var list []models.AdminMenuImport
	var mapAdminMenu map[int]models.AdminMenu
	var mapAdminMenuGroup map[string]models.AdminMenu

	if list, err = daoImportMenu.GetImportMenuByImportIds(importId...); err != nil {
		return
	} else if mapAdminMenuGroup, mapAdminMenu, err = r.getImportMenuGroupMap(dao, list); err != nil {
		return
	}
	res = make(map[int][]wrappers.AdminImportListMenu, l)

	var (
		dt  wrappers.AdminImportListMenu
		dtm models.AdminMenu
		ok  bool
		ll  = len(list)
	)
	var dtt models.AdminMenu
	for _, value := range list {
		if _, ok := res[value.ImportId]; !ok {
			res[value.ImportId] = make([]wrappers.AdminImportListMenu, 0, ll)
		}
		dt = wrappers.AdminImportListMenu{
			ImportId: value.Id,
			MenuId:   value.MenuId,
		}

		if dtm, ok = mapAdminMenu[value.MenuId]; ok {
			dt.MenuName = dtm.Label
			dt.Id = dtm.Id
			if dtt, ok = mapAdminMenuGroup[dtm.Module]; ok {
				dt.SystemName = dtt.Label
				dt.SystemModuleId = dtt.Id
				dt.SystemMenuKey = dtt.PermitKey
				dt.SystemIcon = dtt.Icon
			}
		}
		res[value.ImportId] = append(res[value.ImportId], dt)
	}
	return
}

func (r *SrvPermitImport) getImportMenuGroupMap(dao daos.DaoPermit, list []models.AdminMenuImport) (mapAdminMenuModule map[string]models.AdminMenu, mapAdminMenu map[int]models.AdminMenu, err error) {
	ll := len(list)
	menuIds := make([]int, 0, ll)
	for _, value := range list {
		menuIds = append(menuIds, value.MenuId)
	}
	var adminMenu []models.AdminMenu
	if adminMenu, err = dao.GetMenu(menuIds...); err != nil {
		return
	}
	mapAdminMenu = make(map[int]models.AdminMenu, len(adminMenu))
	var (
		m          = make(map[int]int, len(adminMenu))
		modules    = make([]string, 0, len(adminMenu))
		modulesMap = make(map[string]string, len(adminMenu))
	)
	for _, value := range adminMenu {
		if _, ok := modulesMap[value.Module]; !ok {
			modules = append(modules, value.Module)
		}
	}

	var dta []models.AdminMenu
	if dta, err = dao.GetMenuByPermitKey("", modules...); err != nil {
		return
	}

	mapAdminMenuModule = make(map[string]models.AdminMenu, len(dta))
	for _, value := range dta {
		mapAdminMenuModule[value.PermitKey] = value
	}

	for _, value := range adminMenu {
		if _, ok := m[value.Id]; !ok {
			mapAdminMenu[value.Id] = value
			m[value.Id] = value.Id
		}
	}

	return
}

func NewSrvPermitImport(ctx ...*base.Context) (res srvs.SrvPermitImport) {
	p := &SrvPermitImport{}
	p.SetContext(ctx...)
	return p
}
