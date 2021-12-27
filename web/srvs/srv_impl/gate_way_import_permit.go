package srv_impl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type SrvGatewayImportImpl struct {
	base.ServiceBase
}

func (r *SrvGatewayImportImpl) GetUerImportPermit(arg *wrapper_intranet.ArgGetUerImportPermit) (res *wrapper_intranet.ResultGetUerImportPermit, err error) {
	res = wrapper_intranet.NewResultGetUerImportPermit(arg)
	var (
		isAdmin      bool
		isSuperAdmin bool
		groupIds     []int64
	)
	if isAdmin, isSuperAdmin, groupIds, err = NewSrvPermitGroupImpl(r.Context).GetUserGroup(arg.UHid); err != nil {
		return
	}
	if len(groupIds) == 0 {
		err = fmt.Errorf("您没有权限访问本系统")
		return
	}
	if !isAdmin {
		err = fmt.Errorf("您不是系统管理员,无权限访问")
		return
	}
	if isSuperAdmin { // 如果是超级管理员
		res.IsSuper = true
		for _, item := range arg.UrlInfo {
			res.MapHavePermit[item.ToUk()] = true
		}
		return
	}

	// 不是超级管理员
	if err = r.getImportPermitAsGeneralAdmin(arg, groupIds, res); err != nil {
		return
	}
	return
}

// 普通用户判断是否有接口权限
func (r *SrvGatewayImportImpl) getImportPermitAsGeneralAdmin(arg *wrapper_intranet.ArgGetUerImportPermit, groupIds []int64, res *wrapper_intranet.ResultGetUerImportPermit) (err error) {

	var (
		apps       = arg.GetUrlApps()
		permitList []wrapper_intranet.AdminUserGroupPermit
	)

	if permitList, err = r.GetUserGroupAppPermit(groupIds, apps...); err != nil {
		return
	}

	for _, permit := range permitList {
		for _, item := range arg.UrlInfo {
			if item.App != "" && permit.AppName == item.App {
				res.MapHavePermit[item.ToUk()] = permit.AdminImport.MatchPath(item.Uri, item.Method)
			}
		}
	}
	return
}

// GetUserGroupAppPermit 获取用户组的每个APP的权限
func (r *SrvGatewayImportImpl) GetUserGroupAppPermit(groupIds []int64, apps ...string) (res []wrapper_intranet.AdminUserGroupPermit, err error) {

	var (
		dtm []wrapper_intranet.AdminUserGroupPermit
		dao = dao_impl.NewDaoPermitGroup(r.Context)
	)

	for _, groupId := range groupIds {
		for _, app := range apps {
			if dtm, err = dao.GetGroupAppPermitImport(groupId, app); err != nil {
				return
			}
			res = append(res, dtm...)
		}
	}

	return
}

func (r *SrvGatewayImportImpl) GetImportPermit(arg *wrapper_intranet.ArgGetImportPermit) (res wrapper_intranet.ResultGetImportPermit, err error) {

	res = wrapper_intranet.ResultGetImportPermit{RouterNotNeedSign: map[string]*wrapper_intranet.RouterNotNeedItem{}, RouterNotNeedLogin: map[string]*wrapper_intranet.RouterNotNeedItem{}}

	dao := dao_impl.NewDaoPermitImport(r.Context)

	var importList []models.AdminImport

	if importList, err = dao.GetImportForGateway(arg); err != nil {
		return
	}
	for _, item := range importList {

		if item.AppName == "" {
			continue
		}
		item.UrlPath = fmt.Sprintf("/%s", strings.TrimLeft(item.UrlPath, "/"))
		if err = r.needSignNot(&res, &item); err != nil {
			return
		}

		if err = r.needLoginNot(&res, &item); err != nil {
			return
		}

	}
	return

}

func (r *SrvGatewayImportImpl) needLoginNot(res *wrapper_intranet.ResultGetImportPermit, item *models.AdminImport) (err error) {
	if item.NeedLogin == models.NeedLoginTrue {
		return
	}
	var reGxp string
	if res.RouterNotNeedLogin == nil {
		res.RouterNotNeedLogin = map[string]*wrapper_intranet.RouterNotNeedItem{}
	}
	if _, ok := res.RouterNotNeedLogin[item.AppName]; !ok {
		res.RouterNotNeedLogin[item.AppName] = &wrapper_intranet.RouterNotNeedItem{
			GeneralPath: map[string]wrapper_intranet.ItemGateway{},
			RegexpPath:  []wrapper_intranet.ItemGateway{},
		}
	}

	var notHaveRegexp bool

	if reGxp, notHaveRegexp, err = r.RoutePathToRegexp(item.UrlPath); err != nil {
		return
	}

	if notHaveRegexp {
		res.RouterNotNeedLogin[item.AppName].GeneralPath[item.UrlPath] = wrapper_intranet.ItemGateway{Methods: item.GetRequestMethodMap()}
		return
	}
	res.RouterNotNeedLogin[item.AppName].RegexpPath = append(res.RouterNotNeedLogin[item.AppName].RegexpPath, wrapper_intranet.ItemGateway{
		Uri:     reGxp,
		Methods: item.GetRequestMethodMap(),
	})

	return
}

func (r *SrvGatewayImportImpl) needSignNot(res *wrapper_intranet.ResultGetImportPermit, item *models.AdminImport) (err error) {

	if item.NeedSign == models.NeedSignTrue {
		return
	}

	var reGxp string
	if res.RouterNotNeedSign == nil {
		res.RouterNotNeedSign = map[string]*wrapper_intranet.RouterNotNeedItem{}
	}
	if _, ok := res.RouterNotNeedSign[item.AppName]; !ok {
		res.RouterNotNeedSign[item.AppName] = &wrapper_intranet.RouterNotNeedItem{
			GeneralPath: map[string]wrapper_intranet.ItemGateway{},
			RegexpPath:  []wrapper_intranet.ItemGateway{},
		}
	}
	var notHaveRegexp bool
	if reGxp, notHaveRegexp, err = r.RoutePathToRegexp(item.UrlPath); err != nil {
		return
	}

	if notHaveRegexp {
		res.RouterNotNeedSign[item.AppName].GeneralPath[item.UrlPath] = wrapper_intranet.ItemGateway{
			Methods: item.GetRequestMethodMap(),
		}
		return
	}

	res.RouterNotNeedSign[item.AppName].RegexpPath = append(res.RouterNotNeedSign[item.AppName].RegexpPath, wrapper_intranet.ItemGateway{
		Uri:     reGxp,
		Methods: item.GetRequestMethodMap(),
	})

	return
}

func (r *SrvGatewayImportImpl) RoutePathMath(regexpString, path string) (matched bool, err error) {
	matched, err = regexp.Match(regexpString, []byte(path))
	return
}

func (r *SrvGatewayImportImpl) RoutePathToRegexp(path string) (regexpString string, notHaveRegexp bool, err error) {
	var mat *regexp.Regexp
	mat, err = regexp.Compile(":[^/]+")
	regexpString = fmt.Sprintf("^%s$", mat.ReplaceAllString(path, "([^/]+)"))
	if regexpString == fmt.Sprintf("^%s$", path) {
		notHaveRegexp = true
	}
	return
}

// 获取指定appName下的接口列表
// func (r *SrvGatewayImportPermitImpl) getAppImportList(appName, userHid string) (list []models.AdminImport, err error) {
// 	var (
// 		mapAdminImport map[int]models.AdminImport
// 		listMenuImport []models.AdminMenuImport
// 	)
//
// 	mapAdminImport, err = dao_impl.NewDaoGatewayPermit(r.Context).
// 		GetImportListByAppName(appName)
//
// 	if userHid != "" {
// 		if listMenuImport, err = r.getUserGroupImportList(appName, userHid); err != nil {
// 			return
// 		}
// 		for _, item := range listMenuImport {
// 			if dt, ok := mapAdminImport[item.ImportId]; ok {
// 				list = append(list, dt)
// 			}
// 		}
// 		return
// 	}
// 	return
// }

// TODO 获取用户的接口权限
// func (r *SrvGatewayImportPermitImpl) getUserGroupImportList(appName, userHid string) (list []models.AdminMenuImport, err error) {
// 	_, _ = appName, userHid
// 	return
// }

func (r *SrvGatewayImportImpl) inSlice(val string, slice []string) (ok bool) {
	for _, item := range slice {
		if val == item {
			ok = true
			return
		}
	}
	return
}

func NewSrvGatewayImport(context ...*base.Context) (res srvs.SrvGatewayImportPermit) {
	p := &SrvGatewayImportImpl{}
	p.SetContext(context...)
	return p
}
