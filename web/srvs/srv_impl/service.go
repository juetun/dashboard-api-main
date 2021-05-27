/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 11:30 下午
 */
package srv_impl

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type SrvServiceImpl struct {
	base.ServiceBase
}

func (r *SrvServiceImpl) add(dao daos.DaoService, arg *wrappers.ArgServiceEdit) (res bool, err error) {
	t := base.TimeNormal{Time: time.Now()}
	dt := &models.AdminApp{
		UniqueKey:  arg.UniqueKey,
		Port:       arg.Port,
		Name:       arg.Name,
		HostConfig: arg.HostConfig,
		Desc:       arg.Desc,
		IsStop:     arg.IsStop,
		CreatedAt:  t,
		UpdatedAt:  t,
	}
	dt.MarshalHosts()
	if err = dao.Create(dt); err != nil {
		return
	}
	res = true
	return
}
func (r *SrvServiceImpl) update(dao daos.DaoService, arg *wrappers.ArgServiceEdit) (res bool, err error) {
	var apps []models.AdminApp
	if apps, err = dao.GetByIds(arg.Id); err != nil {
		return
	}

	if len(apps) == 0 {
		err = fmt.Errorf("您要编辑的服务不存在或已删除")
		r.Context.Error(map[string]interface{}{
			"arg": arg,
			"err": err.Error(),
		}, "SrvServiceImplEdit")
		return
	}
	t := base.TimeNormal{Time: time.Now()}
	apps[0].HostConfig = arg.HostConfig
	apps[0].MarshalHosts()
	data := map[string]interface{}{
		"unique_key": arg.UniqueKey,
		"port":       arg.Port,
		"hosts":      apps[0].Hosts,
		"name":       arg.Name,
		"desc":       arg.Desc,
		"is_stop":    arg.IsStop,
		"updated_at": t.Format("2006-01-02 16:04:05"),
	}
	condition := map[string]interface{}{"id": arg.Id}
	if err = dao.Update(condition, data); err != nil {
		return
	}
	res = true
	return
}
func (r *SrvServiceImpl) Edit(arg *wrappers.ArgServiceEdit) (res *wrappers.ResultServiceEdit, err error) {
	res = &wrappers.ResultServiceEdit{}
	dao := dao_impl.NewDaoServiceImpl(r.Context)
	if arg.Id == 0 { // 如果是添加
		res.Result, err = r.add(dao, arg)
		return
	}
	res.Result, err = r.update(dao, arg)
	return
}

func (r *SrvServiceImpl) Detail(arg *wrappers.ArgDetail) (res *wrappers.ResultDetail, err error) {
	res = &wrappers.ResultDetail{}
	if arg.Id == 0 {
		return
	}
	dao := dao_impl.NewDaoServiceImpl(r.Context)
	var dt []models.AdminApp
	if dt, err = dao.GetByIds(arg.Id); err != nil {
		return
	}
	if len(dt) == 0 {
		err = fmt.Errorf("您要查看的服务信息不存在或已删除")
		return
	}
	dt[0].UnmarshalHosts()
	res.AdminApp = dt[0]
	return
}
func (r *SrvServiceImpl) List(arg *wrappers.ArgServiceList) (res *wrappers.ResultServiceList, err error) {
	res = &wrappers.ResultServiceList{
		Pager: response.NewPager(),
	}
	var db *gorm.DB
	dao := dao_impl.NewDaoServiceImpl(r.Context)
	// 获取分页数据
	if err = res.Pager.CallGetPagerData(func(pagerObject *response.Pager) (err error) {
		pagerObject.TotalCount, db, err = dao.GetCount(db, arg)
		return
	}, func(pagerObject *response.Pager) (err error) {
		var list []models.AdminApp
		list, err = dao.GetList(db, arg, pagerObject)
		pagerObject.List, err = r.orgList(dao, list)
		return
	}); err != nil {
		return
	}
	return
}

func (r *SrvServiceImpl) orgList(dao daos.DaoService, list []models.AdminApp) (res []wrappers.AdminApp, err error) {
	res = make([]wrappers.AdminApp, 0, len(list))

	var dt = wrappers.AdminApp{}
	for _, it := range list {
		dt.Id = it.Id
		dt.UniqueKey = it.UniqueKey
		dt.Port = it.Port
		dt.Name = it.Name
		dt.Desc = it.Desc
		dt.IsStop = it.IsStop
		dt.CreatedAt = it.CreatedAt
		dt.UpdatedAt = it.UpdatedAt
		res = append(res, dt)
	}
	return
}
func NewSrvServiceImpl(context ...*base.Context) srvs.SrvService {
	p := &SrvServiceImpl{}
	p.SetContext(context...)
	return p
}
