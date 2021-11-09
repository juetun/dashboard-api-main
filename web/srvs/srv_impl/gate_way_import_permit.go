package srv_impl

import (
	"fmt"
	"regexp"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_intranet"
)

type SrvGatewayImportPermitImpl struct {
	base.ServiceBase
}

func (r *SrvGatewayImportPermitImpl) GetImportPermit(arg *wrapper_intranet.ArgGetImportPermit) (res *wrapper_intranet.ResultGetImportPermit, err error) {

	var (
		ok         bool
		importList []models.AdminImport
	)
	res = &wrapper_intranet.ResultGetImportPermit{}

	importList, err = r.getAppImportList(arg.AppName, arg.UserHid)

	// 先路径用等于匹配
	if ok = r.equalFlag(arg, res, importList); ok {
		return
	}

	// 如果路径匹配不到则用正则表达式匹配
	if ok, err = r.regexpFlag(arg, res, importList); ok {
		return
	}
	return

}
func (r *SrvGatewayImportPermitImpl) regexpFlag(arg *wrapper_intranet.ArgGetImportPermit, res *wrapper_intranet.ResultGetImportPermit, importList []models.AdminImport) (ok bool, err error) {

	var (
		okFlag       bool
		regexpString string
	)

	for _, adminImport := range importList {

		if !r.inSlice(arg.Method, adminImport.GetRequestMethods()) { // 如果requestMethod不匹配，则接跳过
			continue
		}
		if regexpString, err = r.RoutePathToRegexp(adminImport.UrlPath); err != nil {
			err = fmt.Errorf("系统配置的接口路径不正确(%d)", adminImport.Id)
			return
		}

		if okFlag, err = r.RoutePathMath(regexpString, fmt.Sprintf("%s/%s", arg.PathType, arg.Uri)); err != nil {
			err = fmt.Errorf("匹配路由参数异常%s(system:%s)", regexpString, fmt.Sprintf("%s/%s", arg.PathType, arg.Uri))
			return
		}

		if okFlag {
			res.NeedLogin = adminImport.NeedLogin > 0
			res.NeedSign = adminImport.NeedSign > 0
			ok = true
			return
		}
	}
	return
}
func (r *SrvGatewayImportPermitImpl) RoutePathMath(regexpString, path string) (matched bool, err error) {
	matched, err = regexp.Match(regexpString, []byte(path))
	return
}

func (r *SrvGatewayImportPermitImpl) RoutePathToRegexp(path string) (regexpString string, err error) {
	var mat *regexp.Regexp
	mat, err = regexp.Compile(":[^/]+")
	regexpString = fmt.Sprintf("^%s$", mat.ReplaceAllString(path, "([^/]+)"))
	return
}

// 获取指定appName下的接口列表
func (r *SrvGatewayImportPermitImpl) getAppImportList(appName, userHid string) (list []models.AdminImport, err error) {
	var mapAdminImport map[int]models.AdminImport
	var listMenuImport []models.AdminMenuImport
	mapAdminImport, err = dao_impl.NewDaoGatewayPermit(r.Context).
		GetImportListByAppName(appName, false)
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

// TODO
func (r *SrvGatewayImportPermitImpl) getUserGroupImportList(appName, userHid string) (list []models.AdminMenuImport, err error) {

	return
}

// 先路径用等于匹配
func (r *SrvGatewayImportPermitImpl) equalFlag(arg *wrapper_intranet.ArgGetImportPermit, res *wrapper_intranet.ResultGetImportPermit, importList []models.AdminImport) (ok bool) {

	for _, adminImport := range importList {
		if adminImport.UrlPath == fmt.Sprintf("%s/%s", arg.PathType, arg.Uri) && r.inSlice(arg.Method, adminImport.GetRequestMethods()) {
			res.NeedLogin = adminImport.NeedLogin > 0
			res.NeedSign = adminImport.NeedSign > 0
			ok = true
			return
		}
	}
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
