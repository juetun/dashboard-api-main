// Package srv_impl
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

type SrvServiceImpl struct {
	base.ServiceBase
}

func (r *SrvServiceImpl) add(dao daos.DaoService, arg *wrappers.ArgServiceEdit) (res bool, err error) {
	t := time.Now()
	dt := &models.AdminApp{
		UniqueKey:    arg.UniqueKey,
		Port:         arg.Port,
		Name:         arg.Name,
		HostConfig:   arg.HostConfig,
		Desc:         arg.Desc,
		IsStop:       uint8(arg.IsStop),
		SupportCache: arg.SupportCache,
		CreatedAt:    t,
		UpdatedAt:    t,
	}
	if err = dt.MarshalHosts(); err != nil {
		return
	}
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
	if err = apps[0].MarshalHosts(); err != nil {
		return
	}
	data := map[string]interface{}{
		"unique_key":    arg.UniqueKey,
		"port":          arg.Port,
		"hosts":         apps[0].Hosts,
		"name":          arg.Name,
		"desc":          arg.Desc,
		"is_stop":       arg.IsStop,
		"support_cache": arg.SupportCache,
		"updated_at":    t.Format(utils.DateGeneral),
	}

	condition := map[string]interface{}{"id": arg.Id}
	if err = dao.Update(condition, data); err != nil {
		return
	}
	res = true
	return
}

func (r *SrvServiceImpl) validateUniqueKey(dao daos.DaoService, arg *wrappers.ArgServiceEdit) (err error) {
	var (
		argNumber = base.NewArgGetByStringIds(base.ArgGetByStringIdsOptionIds(arg.UniqueKey))
		mapApp    map[string]*models.AdminApp
		dt        *models.AdminApp
	)
	if mapApp, err = dao.GetByUniqueKey(argNumber); err != nil {
		return
	}
	dt, _ = mapApp[arg.UniqueKey]
	if dt != nil {
		if arg.Id == 0 {
			err = fmt.Errorf("唯一KEY(%s)已被使用(%s)", dt.UniqueKey, dt.Name)
			return
		} else if arg.Id != dt.Id {
			err = fmt.Errorf("唯一KEY(%s)已被使用(%s)", dt.UniqueKey, dt.Name)
			return
		}
	}
	return
}

func (r *SrvServiceImpl) validatePortUse(dao daos.DaoService, arg *wrappers.ArgServiceEdit) (err error) {
	var (
		argNumber  = base.NewArgGetByNumberIds(base.ArgGetByNumberIdsOptionIds(int64(arg.Port)))
		mapPortApp map[int64]*models.AdminApp
		dt         *models.AdminApp
	)
	if mapPortApp, err = dao.GetByPort(argNumber); err != nil {
		return
	}
	dt, _ = mapPortApp[int64(arg.Port)]
	if dt != nil {
		if arg.Id == 0 {
			err = fmt.Errorf("端口(%d)已被使用(%s)", dt.Port, dt.Name)
			return
		} else if arg.Id != dt.Id {
			err = fmt.Errorf("端口(%d)已被使用(%s)", dt.Port, dt.Name)
			return
		}
	}
	return
}

func (r *SrvServiceImpl) Edit(arg *wrappers.ArgServiceEdit) (res *wrappers.ResultServiceEdit, err error) {
	res = &wrappers.ResultServiceEdit{}
	var (
		dao = dao_impl.NewDaoServiceImpl(r.Context)
	)
	if err = r.validatePortUse(dao, arg); err != nil {
		return
	}

	if err = r.validateUniqueKey(dao, arg); err != nil {
		return
	}
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
	if err = dt[0].UnmarshalHosts(); err != nil {
		return
	}
	if err = dt[0].Default(); err != nil {
		return
	}
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

func (r *SrvServiceImpl) orgList(dao daos.DaoService, list []models.AdminApp) (res []*wrappers.AdminApp, err error) {
	res = make([]*wrappers.AdminApp, 0, len(list))

	var dt *wrappers.AdminApp
	for _, it := range list {
		dt = &wrappers.AdminApp{}
		dt.Id = it.Id
		dt.UniqueKey = it.UniqueKey
		dt.Port = it.Port
		dt.Name = it.Name
		dt.Desc = it.Desc
		dt.SupportCache = it.SupportCache
		if it.Desc == "" {
			dt.DisableExpand = true
		}
		dt.IsStopName = it.ParseIsStop()

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
