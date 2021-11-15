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

type SrvGatewayImportPermitImpl struct {
	base.ServiceBase
}

func (r *SrvGatewayImportPermitImpl) GetImportPermit(arg *wrapper_intranet.ArgGetImportPermit) (res wrapper_intranet.ResultGetImportPermit, err error) {
	res = wrapper_intranet.ResultGetImportPermit{
		RouterNotNeedSign:  map[string]*wrapper_intranet.RouterNotNeedItem{},
		RouterNotNeedLogin: map[string]*wrapper_intranet.RouterNotNeedItem{},
	}

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

func (r *SrvGatewayImportPermitImpl) needLoginNot(res *wrapper_intranet.ResultGetImportPermit, item *models.AdminImport) (err error) {
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
		res.RouterNotNeedLogin[item.AppName].GeneralPath[item.UrlPath] = wrapper_intranet.ItemGateway{
			Methods: item.GetRequestMethods(),
		}
		return
	}
	res.RouterNotNeedLogin[item.AppName].RegexpPath = append(res.RouterNotNeedLogin[item.AppName].RegexpPath, wrapper_intranet.ItemGateway{
		Uri:     reGxp,
		Methods: item.GetRequestMethods(),
	})

	return
}

func (r *SrvGatewayImportPermitImpl) needSignNot(res *wrapper_intranet.ResultGetImportPermit, item *models.AdminImport) (err error) {

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
			Methods: item.GetRequestMethods(),
		}
		return
	}

	res.RouterNotNeedSign[item.AppName].RegexpPath = append(res.RouterNotNeedSign[item.AppName].RegexpPath, wrapper_intranet.ItemGateway{
		Uri:     reGxp,
		Methods: item.GetRequestMethods(),
	})

	return
}

func (r *SrvGatewayImportPermitImpl) RoutePathMath(regexpString, path string) (matched bool, err error) {
	matched, err = regexp.Match(regexpString, []byte(path))
	return
}

func (r *SrvGatewayImportPermitImpl) RoutePathToRegexp(path string) (regexpString string, notHaveRegexp bool, err error) {
	var mat *regexp.Regexp
	mat, err = regexp.Compile(":[^/]+")
	regexpString = fmt.Sprintf("^%s$", mat.ReplaceAllString(path, "([^/]+)"))
	if regexpString == fmt.Sprintf("^%s$", path) {
		notHaveRegexp = true
	}
	return
}

// 获取指定appName下的接口列表
func (r *SrvGatewayImportPermitImpl) getAppImportList(appName, userHid string) (list []models.AdminImport, err error) {
	var (
		mapAdminImport map[int]models.AdminImport
		listMenuImport []models.AdminMenuImport
	)

	mapAdminImport, err = dao_impl.NewDaoGatewayPermit(r.Context).
		GetImportListByAppName(appName)

	if userHid != "" {
		if listMenuImport, err = r.getUserGroupImportList(appName, userHid); err != nil {
			return
		}
		for _, item := range listMenuImport {
			if dt, ok := mapAdminImport[item.ImportId]; ok {
				list = append(list, dt)
			}
		}
		return
	}
	return
}

// TODO 获取用户的接口权限
func (r *SrvGatewayImportPermitImpl) getUserGroupImportList(appName, userHid string) (list []models.AdminMenuImport, err error) {

	return
}

func (r *SrvGatewayImportPermitImpl) inSlice(val string, slice []string) (ok bool) {
	for _, item := range slice {
		if val == item {
			ok = true
			return
		}
	}
	return
}

func NewSrvGatewayImportPermit(context ...*base.Context) (res srvs.SrvGatewayImportPermit) {
	p := &SrvGatewayImportPermitImpl{}
	p.SetContext(context...)
	return p
}
