/**
* @Author:changjiang
* @Description:
* @File:permit_app
* @Version: 1.0.0
* @Date 2021/9/12 1:33 下午
 */
package srv_impl

import (
	"fmt"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvPermitAppImpl struct {
	base.ServiceDao
}

func (r *SrvPermitAppImpl) GetAppConfig(arg *wrappers.ArgGetAppConfig) (res *wrappers.ResultGetAppConfig, err error) {
	res = &wrappers.ResultGetAppConfig{}
	dao := dao_impl.NewDaoServiceImpl(r.Context)
	if arg.Module == "" {
		err = r.getAppConfigList(dao, arg, res)
		return
	}
	err = r.getAppConfigListByModule(dao, arg, res)
	return
}

func (r *SrvPermitAppImpl) getImportMenu(dao daos.DaoService, arg *wrappers.ArgGetAppConfig) (uniqueKeys []string, err error) {
	var importMenus []wrappers.ImportMenu
	if importMenus, err = dao.GetImportMenuByModule(arg.Module); err != nil {
		return
	}
	uniqueKeys = make([]string, 0, len(importMenus))
	for _, it := range importMenus {
		uniqueKeys = append(uniqueKeys, it.AppName)
	}
	return
}

func (r *SrvPermitAppImpl) getAppConfigListByModule(dao daos.DaoService, arg *wrappers.ArgGetAppConfig, res *wrappers.ResultGetAppConfig) (err error) {
	var (
		list       []models.AdminApp
		uniqueKeys []string
	)
	if uniqueKeys, err = r.getImportMenu(dao, arg); err != nil {
		return
	}

	if list, err = dao.GetList(nil, &wrappers.ArgServiceList{
		UniqueKeys: uniqueKeys,
	}, nil); err != nil {
		return
	}
	if err = r.orgResAppConfig(list, res, arg); err != nil {
		return
	}
	return
}

func (r *SrvPermitAppImpl) getAppConfigList(dao daos.DaoService, arg *wrappers.ArgGetAppConfig, res *wrappers.ResultGetAppConfig) (err error) {
	var list []models.AdminApp
	if list, err = dao.GetList(nil, nil, nil); err != nil {
		return
	}
	err = r.orgResAppConfig(list, res, arg)
	return
}
func (r *SrvPermitAppImpl) orgResAppConfig(list []models.AdminApp, res *wrappers.ResultGetAppConfig, arg *wrappers.ArgGetAppConfig) (err error) {
	*res = make(map[string]string, len(list))
	for _, it := range list {
		if err = it.UnmarshalHosts(); err != nil {
			r.Context.Error(map[string]interface{}{
				"err": err.Error(),
				"it":  it,
				"arg": arg,
			}, "permitServiceGetAppConfigList")
			return
		}
		(*res)[it.UniqueKey] = fmt.Sprintf("http://%s:%d", it.HostConfig[arg.Env], it.Port)
	}
	return
}

func NewSrvPermitAppImpl(ctx ...*base.Context) (res srvs.SrvPermitApp) {
	p := &SrvPermitAppImpl{}
	p.SetContext(ctx...)
	return p
}
