// Package srv_impl
// /**
package srv_impl

import (
	"fmt"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"gorm.io/gorm"
)

type (
	SrvPermitImport struct {
		SrvPermitCommon
	}
	userPageImportHandler func(arg *UserPageImportParam, res *wrapper_admin.ResultPageImport) (err error)
	UserPageImportParam   struct {
		wrapper_admin.ArgPageImport
		AdminMenu       models.AdminMenu `json:"admin_menu" form:"-"`
		IsSupperAdmin   bool             `json:"is_supper_admin" form:"-"`
		OperatorGroupId []int64          `json:"operator_group_id" form:"-"`
	}
)

//用户页面具备的接口权限列表
func (r *SrvPermitImport) UserPageImport(arg *wrapper_admin.ArgPageImport) (res *wrapper_admin.ResultPageImport, err error) {
	res = wrapper_admin.NewResultPageImport()
	var (
		liAdminMenu []*models.AdminMenu
		argData     = UserPageImportParam{
			ArgPageImport: *arg,
		}
	)

	//获取用户是否是超级管理员，或用户组ID
	if err = r.userPermit(&argData); err != nil {
		return
	}
	res.IsSuperAdmin = argData.IsSupperAdmin

	if liAdminMenu, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetMenuByPermitKey(arg.Module, arg.PageName); err != nil {
		return
	}

	if len(liAdminMenu) == 0 {
		err = fmt.Errorf("您选择的菜单未配置(%s)", arg.PageName)
		return
	}

	argData.AdminMenu = *(liAdminMenu[0])
	var handlers = []userPageImportHandler{
		r.initUserCommonImport, //页面公共接口权限查询
		r.initUserPageImport,   //页面接口权限查询
		r.initUserSubPage,      //子页面查询
	}

	for _, item := range handlers {
		if err = item(&argData, res); err != nil {
			return
		}
	}
	return
}

//页面公共接口权限查询
func (r *SrvPermitImport) initUserCommonImport(arg *UserPageImportParam, res *wrapper_admin.ResultPageImport) (err error) {
	var (
		adminMenu                = &models.AdminMenu{Module: arg.AdminMenu.Module}
		moduleCommonImportString = adminMenu.GetCommonImportString(arg.AdminMenu.Module)
		e                        error
	)

	if adminMenu, e = dao_impl.NewDaoPermitMenu(r.Context).
		GetAdminMenuByModule(moduleCommonImportString); e != nil {
		res.ShowError = true
		res.ErrorMsg = fmt.Sprintf("系统配置(%s)不存在或已删除", arg.AdminMenu.Module)
		return
	}
	if adminMenu == nil {
		res.ShowError = true
		res.ErrorMsg = fmt.Sprintf("系统配置(%s)不存在或已删除", arg.AdminMenu.Module)
		return
	}
	if arg.IsSupperAdmin { //如果是超管
		if e = r.initUserCommonImportSupperAdmin(adminMenu, res); e != nil {
			res.ShowError = true
			res.ErrorMsg = e.Error()
		}
		return
	}
	if e = r.initUserCommonImportGeneral(adminMenu, arg.OperatorGroupId, res); e != nil {
		res.ShowError = true
		res.ErrorMsg = e.Error()
		return
	}
	return
}

func (r *SrvPermitImport) initUserCommonImportSupperAdmin(adminMenu *models.AdminMenu, res *wrapper_admin.ResultPageImport) (err error) {
	if adminMenu == nil {
		return
	}
	var (
		mapImportList map[int64]models.AdminMenuImportCache
		ok            bool
		dt            models.AdminMenuImportCache
		argNumber     = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(adminMenu.Id))
	)
	if mapImportList, err = dao_impl.NewDaoPermitImport(r.Context).
		GetChildImportByMenuId(argNumber); err != nil {
		return
	}
	if dt, ok = mapImportList[adminMenu.Id]; !ok {
		return
	}

	for _, item := range dt {
		res.CommonImport[item.ImportPermitKey] = 1
	}
	return
}

func (r *SrvPermitImport) initUserCommonImportGeneral(adminMenu *models.AdminMenu, groupId []int64, res *wrapper_admin.ResultPageImport) (err error) {
	if adminMenu == nil {
		return
	}
	var (
		resGroupImport  []*models.AdminUserGroupImport
		importIds       []int64
		adminMenuImport map[int64]*models.AdminMenuImport
	)
	if resGroupImport, err = dao_impl.NewDaoPermitGroupImport(r.Context).
		GetMenuIdsByPermitByGroupIds([]int64{adminMenu.Id},
			groupId); err != nil {
		return
	}
	importIds = make([]int64, 0, len(resGroupImport))
	for _, item := range resGroupImport {
		importIds = append(importIds, item.ImportId)
	}
	argNumber := base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(importIds...))
	if adminMenuImport, err = dao_impl.NewDaoPermitImport(r.Context).
		GetImportMenuByImportIds(argNumber); err != nil {
		return
	}
	for _, item := range adminMenuImport {
		res.CommonImport[item.ImportPermitKey] = 1
	}
	return
}

func (r *SrvPermitImport) initUserPageImportSupperAdmin(arg *UserPageImportParam, res *wrapper_admin.ResultPageImport) (err error) {
	var (
		mapImportList map[int64]models.AdminMenuImportCache
		ok            bool
		dt            models.AdminMenuImportCache
		argNumber     = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(arg.AdminMenu.Id))
	)
	if mapImportList, err = dao_impl.NewDaoPermitImport(r.Context).
		GetChildImportByMenuId(argNumber); err != nil {
		return
	}
	if dt, ok = mapImportList[arg.AdminMenu.Id]; !ok {
		return
	}

	for _, item := range dt {
		res.PageImport[item.ImportPermitKey] = 1
	}
	return
}

func (r *SrvPermitImport) initUserPageImportGeneral(arg *UserPageImportParam, res *wrapper_admin.ResultPageImport) (err error) {
	var (
		resGroupImport  []*models.AdminUserGroupImport
		importIds       []int64
		adminMenuImport map[int64]*models.AdminMenuImport
	)
	if resGroupImport, err = dao_impl.NewDaoPermitGroupImport(r.Context).
		GetMenuIdsByPermitByGroupIds([]int64{arg.AdminMenu.Id},
			arg.OperatorGroupId); err != nil {
		return
	}
	importIds = make([]int64, 0, len(resGroupImport))
	for _, item := range resGroupImport {
		importIds = append(importIds, item.ImportId)
	}
	argNumber := base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(importIds...))
	if adminMenuImport, err = dao_impl.NewDaoPermitImport(r.Context).
		GetImportMenuByImportIds(argNumber); err != nil {
		return
	}
	for _, item := range adminMenuImport {
		res.PageImport[item.ImportPermitKey] = 1
	}
	return
}

//页面公共接口权限查询
func (r *SrvPermitImport) initUserPageImport(arg *UserPageImportParam, res *wrapper_admin.ResultPageImport) (err error) {
	if arg.IsSupperAdmin { //如果是超管
		err = r.initUserPageImportSupperAdmin(arg, res)
		return
	}
	err = r.initUserPageImportGeneral(arg, res)
	return
}

func (r *SrvPermitImport) initUserSubPageSupperAdmin(arg *UserPageImportParam, res *wrapper_admin.ResultPageImport) (err error) {

	return
}
func (r *SrvPermitImport) initUserSubPageGeneral(arg *UserPageImportParam, res *wrapper_admin.ResultPageImport) (err error) {

	return
}

//页面公共接口权限查询
func (r *SrvPermitImport) initUserSubPage(arg *UserPageImportParam, res *wrapper_admin.ResultPageImport) (err error) {
	if arg.IsSupperAdmin { //如果是超管
		err = r.initUserSubPageSupperAdmin(arg, res)
		return
	}
	err = r.initUserSubPageGeneral(arg, res)
	return
}

func (r *SrvPermitImport) UpdateImportValue(arg *wrappers.ArgUpdateImportValue) (res *wrappers.ResultUpdateImportValue, err error) {
	res = &wrappers.ResultUpdateImportValue{}
	var condition = fmt.Sprintf("`id` IN (%s)", strings.Join(arg.Ids, ","))
	dao := dao_impl.NewDaoPermitImport(r.Context)
	var data = make(map[string]interface{}, 1)
	t := time.Now().Format(utils.DateTimeGeneral)
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
		ImportIds: []wrappers.ImportSingle{},
		MenuIds:   []wrappers.MenuSingle{},
	}

	if err = NewSrvPermitMenu(r.Context).
		GetMenuPermitKeyByPath(
			&arg.ArgGetImportByMenuIdSingle,
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

func (r *SrvPermitImport) GetChildImport(nowMenuId int64) (importIds []wrappers.ImportSingle, err error) {
	importIds = []wrappers.ImportSingle{}

	var (
		dao        = dao_impl.NewDaoPermitImport(r.Context)
		importMap  map[int64]models.AdminMenuImportCache
		importList models.AdminMenuImportCache
		ok         bool
		arg        = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(nowMenuId))
	)

	if importMap, err = dao.GetChildImportByMenuId(arg); err != nil {
		return
	}
	if importList, ok = importMap[nowMenuId]; !ok {
		return
	}
	importIds = make([]wrappers.ImportSingle, 0, len(importList))
	for _, value := range importList {
		importIds = append(importIds, wrappers.ImportSingle{
			ImportId:  value.ImportId,
			PermitKey: value.ImportPermitKey,
		})
	}
	return
}

func (r *SrvPermitImport) GetOpList(dao daos.DaoPermit, arg *wrappers.ArgPermitMenu) (opList map[string][]wrappers.OpOne, err error) {
	opList = map[string][]wrappers.OpOne{}
	var (
		t    wrappers.OpOne
		list []wrappers.Op
	)
	if list, err = dao.GetPermitImportByModule(arg); err != nil {
		return
	}
	l := len(list)
	opList = make(map[string][]wrappers.OpOne, l)
	for _, value := range list {
		if _, ok := opList[value.MenuPermitKey]; !ok {
			opList[value.MenuPermitKey] = make([]wrappers.OpOne, 0, l)
		}
		t = wrappers.OpOne(value.PermitKey)
		opList[value.MenuPermitKey] = append(opList[value.MenuPermitKey], t)
	}
	return
}

func (r *SrvPermitImport) GetChildMenu(nowMenuId int64) (menuIds []wrappers.MenuSingle, err error) {
	menuIds = []wrappers.MenuSingle{}
	dao := dao_impl.NewDaoPermit(r.Context)
	var res []*models.AdminMenu
	if res, err = dao.GetAdminMenuList(&wrappers.ArgAdminMenu{
		ParentId: nowMenuId,
	}); err != nil {
		return
	}
	for _, item := range res {
		menuIds = append(menuIds, wrappers.MenuSingle{
			MenuId:    item.Id,
			PermitKey: item.PermitKey,
		})
	}
	return
}
func (r *SrvPermitImport) GetImport(arg *wrappers.ArgGetImport) (res *wrappers.ResultGetImport, err error) {
	res = &wrappers.ResultGetImport{
		Pager: response.NewPager(response.PagerBaseQuery(&arg.PageQuery)),
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
	res.List, err = r.joinChecked(arg, list)

	// []models.AdminImport{}
	return
}

func (r *SrvPermitImport) SetApiPermit(arg *wrappers.ArgAdminSetPermit) (err error) {
	dao := dao_impl.NewDaoPermitGroupImport(r.Context)
	switch arg.Act {
	case models.SetPermitAdd:
		var (
			permitImport map[int64]*models.AdminImport
			t            = time.Now()
			list         []base.ModelBase
			dt           *models.AdminUserGroupImport
			menu         *models.AdminMenu
		)

		if menu, err = dao_impl.NewDaoPermitMenu(r.Context).
			GetMenuByMenuId(arg.MenuId); err != nil {
			return
		}
		if permitImport, err = dao_impl.NewDaoPermitImport(r.Context).
			GetImportFromDbByIds(arg.PermitIds...); err != nil {
			return
		}

		for _, pid := range arg.PermitIds {
			dt = &models.AdminUserGroupImport{
				GroupId:       arg.GroupId,
				ImportId:      pid,
				MenuId:        arg.MenuId,
				Module:        menu.Module,
				MenuPermitKey: menu.PermitKey,
				DefaultOpen:   arg.DefaultOpen,
				CreatedAt:     t,
				UpdatedAt:     t,
				DeletedAt:     nil,
			}
			if dtm, ok := permitImport[pid]; ok {
				dt.AppName = dtm.AppName
				dt.ImportPermitKey = dtm.PermitKey
			}
			list = append(list, dt)
		}
		if err = dao.BatchAddDataGroupImport(list); err != nil {
			err = fmt.Errorf("操作异常")
			return
		}

	case models.SetPermitCancel:
		if err = dao.DeleteGroupImportWithGroupIdAndImportIds(arg.GroupId, arg.PermitIds...); err != nil {
			return
		}
	default:
		err = fmt.Errorf("act格式错误")
		return
	}
	return
}

func (r *SrvPermitImport) userPermit(arg *UserPageImportParam) (err error) {
	// 判断当前用户是否是超级管理员,
	var getUserArgument = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(arg.UUserHid))
	if arg.IsSupperAdmin, err = r.GetAdminUserInfo(getUserArgument, arg.UUserHid); err != nil {
		return
	}
	if !arg.IsSupperAdmin { //如果账号不是超管
		// 判断当前用户是否是超级管理员,如果不是超级管理员，组织所属组权限
		if arg.OperatorGroupId, arg.IsSupperAdmin, err = NewSrvPermitUserImpl(r.Context).
			GetUserAdminGroupIdByUserHid(arg.UUserHid); err != nil {
			return
		}
	}
	return
}

func (r *SrvPermitImport) getImportId(l int, list []models.AdminImport) (importId []int64) {
	importId = make([]int64, 0, l)
	for _, value := range list {
		importId = append(importId, value.Id)
	}
	return
}

func (r *SrvPermitImport) ImportList(arg *wrappers.ArgImportList) (res *wrappers.ResultImportList, err error) {
	if arg.Order == "" {
		arg.Order = "id desc"
	}
	res = &wrappers.ResultImportList{Pager: response.NewPagerAndDefault(&arg.PageQuery)}
	dao := dao_impl.NewDaoPermit(r.Context)

	var actRes *base.ActErrorHandlerResult
	// 获取分页数据
	if err = res.Pager.CallGetPagerData(func(pager *response.Pager) (err error) {
		pager.TotalCount, actRes, err = dao.GetImportListCount(arg)
		return
	}, func(pager *response.Pager) (err error) {
		var list []models.AdminImport
		list, err = dao.GetImportListData(actRes, arg, pager)
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
	var dta map[int64][]wrappers.AdminImportListMenu
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

// 批量添加接口
func (r *SrvPermitImport) batchAddImport(arg *wrappers.ArgEditImport) (err error) {
	dao := dao_impl.NewDaoPermitImport(r.Context)
	var dt models.AdminImport
	dt.SetRequestMethods(arg.RequestMethod)
	var dataList []models.AdminImport
	var t = base.TimeNormal{Time: time.Now()}
	var listImport []models.AdminImport
	for _, item := range arg.UrlPaths {
		dt = r.orgImportData(item, arg, &t)
		dataList = append(dataList, dt)
		if listImport, err = dao.GetImportByCondition(map[string]interface{}{"app_name": arg.AppName, "url_path": item}); err != nil {
			return
		}
		// 验证数据是否重复
		for _, value := range listImport {
			if err = r.editImportParam(arg, &value); err != nil {
				return
			}
		}
	}
	for _, item := range dataList {
		if err = dao.AddData(&item); err != nil {
			return
		}
	}
	return
}
func (r *SrvPermitImport) orgImportData(item string, arg *wrappers.ArgEditImport, t *base.TimeNormal) (dt models.AdminImport) {

	dt = models.AdminImport{
		AppName:       arg.AppName,
		AppVersion:    arg.AppVersion,
		UrlPath:       item,
		SortValue:     arg.SortValue,
		RequestMethod: dt.RequestMethod,
		DefaultOpen:   arg.DefaultOpen,
		NeedLogin:     arg.NeedLogin,
		NeedSign:      arg.NeedSign,
		CreatedAt:     t.Time,
		UpdatedAt:     t.Time,
	}
	return
}
func (r *SrvPermitImport) EditImport(arg *wrappers.ArgEditImport) (res *wrappers.ResultEditImport, err error) {
	res = &wrappers.ResultEditImport{Result: false}
	defer func() {
		if err != nil {
			return
		}
		res.Result = true
	}()
	var (
		dao        = dao_impl.NewDaoPermit(r.Context)
		daoImport  = dao_impl.NewDaoPermitImport(r.Context)
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

	// 如果是添加接口，支持批量添加
	if arg.Id == 0 {
		if err = r.batchAddImport(arg); err != nil {
			return
		}
		return
	}

	if listImport, err = daoImport.GetImportByCondition(map[string]interface{}{"app_name": arg.AppName, "url_path": arg.UrlPath}); err != nil {
		return
	}

	// 验证数据是否重复
	for _, value := range listImport {
		if err = r.editImportParam(arg, &value); err != nil {
			return
		}
	}
	var mAi = &models.AdminImport{}
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

	var dt []models.AdminImport
	if dt, err = dao.GetAdminImportById(arg.Id); err != nil {
		return
	}

	if len(dt) == 0 {
		err = fmt.Errorf("您编辑的接口信息不存在或已删除")
		return
	}
	if dt[0].PermitKey == "" {
		data["permit_key"] = dt[0].GetPathName()
		if err = daoImport.UpdateMenuImport(fmt.Sprintf("import_id=%d", dt[0].Id), map[string]interface{}{"import_permit_key": data["permit_key"]}); err != nil {
			return
		}
	}
	// 如果更新了app_name
	if dt[0].AppName != arg.AppName {
		if err = daoImport.UpdateMenuImport(fmt.Sprintf("import_id = %d ", dt[0].Id),
			map[string]interface{}{"import_app_name": arg.AppName}); err != nil {
			return
		}

	}
	if _, err = dao.UpdateAdminImport(map[string]interface{}{"id": arg.Id}, data); err != nil {
		return
	}
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
	dao := dao_impl.NewDaoPermitImport(r.Context)
	if err = dao.DeleteImportByIds([]int64{arg.ID}...); err != nil {
		return
	}
	daoGroup := dao_impl.NewDaoPermitGroupImport(r.Context)
	if err = daoGroup.DeleteGroupImportWithGroupIdAndImportIds(0, []int64{arg.ID}...); err != nil {
		return
	}
	res.Result = true
	return
}

func (r *SrvPermitImport) joinChecked(arg *wrappers.ArgGetImport, data []models.AdminImport) (res []wrappers.AdminImport, err error) {
	res = make([]wrappers.AdminImport, 0, len(data))
	var dt wrappers.AdminImport
	var importId = make([]int64, 0, len(data))
	for _, value := range data {
		importId = append(importId, value.Id)
	}
	daoGroupImport := dao_impl.NewDaoPermitGroupImport(r.Context)
	var li []models.AdminUserGroupImport
	if li, err = daoGroupImport.GetSelectImportByImportId(arg.GroupId, importId...); err != nil {
		return
	}
	var m = make(map[int64]int64, len(li))
	for _, it := range li {
		m[it.ImportId] = it.ImportId
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

func (r *SrvPermitImport) getImportMenuGroup(dao daos.DaoPermit, l int, data []models.AdminImport) (res map[int64][]wrappers.AdminImportListMenu, err error) {
	importId := r.getImportId(l, data)
	daoImportMenu := dao_impl.NewDaoPermitImport(r.Context)

	var list map[int64]*models.AdminMenuImport
	var mapAdminMenu map[int64]*models.AdminMenu
	var mapAdminMenuGroup map[string]*models.AdminMenu
	var argNumber = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(importId...))
	if list, err = daoImportMenu.GetImportMenuByImportIds(argNumber); err != nil {
		return
	} else if mapAdminMenuGroup, mapAdminMenu, err = r.getImportMenuGroupMap(dao, list); err != nil {
		return
	}
	res = make(map[int64][]wrappers.AdminImportListMenu, l)

	var (
		dt  wrappers.AdminImportListMenu
		dtm *models.AdminMenu
		ok  bool
		ll  = len(list)
	)
	var dtt *models.AdminMenu
	for _, value := range list {
		if _, ok = res[value.ImportId]; !ok {
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

func (r *SrvPermitImport) getImportMenuGroupMap(dao daos.DaoPermit, list map[int64]*models.AdminMenuImport) (mapAdminMenuModule map[string]*models.AdminMenu, mapAdminMenu map[int64]*models.AdminMenu, err error) {
	ll := len(list)
	menuIds := make([]int64, 0, ll)
	var mapMenuIds = make(map[int64]bool, ll)
	for _, value := range list {
		if _, ok := mapMenuIds[value.MenuId]; !ok {
			menuIds = append(menuIds, value.MenuId)
			mapMenuIds[value.MenuId] = true
		}
	}
	var adminMenu []*models.AdminMenu
	if adminMenu, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetMenu(menuIds...); err != nil {
		return
	}

	mapAdminMenu = make(map[int64]*models.AdminMenu, len(adminMenu))

	var (
		m          = make(map[int64]int64, len(adminMenu))
		modules    = make([]string, 0, len(adminMenu))
		modulesMap = make(map[string]string, len(adminMenu))
		dta        []*models.AdminMenu
	)

	for _, value := range adminMenu {
		if _, ok := modulesMap[value.Module]; !ok {
			modules = append(modules, value.Module)
		}
	}

	if dta, err = dao_impl.NewDaoPermitMenu(r.Context).
		GetMenuByPermitKey("", modules...); err != nil {
		return
	}

	mapAdminMenuModule = make(map[string]*models.AdminMenu, len(dta))

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
