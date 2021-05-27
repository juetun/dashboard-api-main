/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:34 上午
 */
package srv_impl

import (
	"time"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/basic/utils"
	"github.com/juetun/dashboard-api-main/web"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl/export"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ServiceExport struct {
	base.ServiceBase
}

func NewServiceExport(context ...*base.Context) (p *ServiceExport) {
	p = &ServiceExport{}
	p.SetContext(context...)
	p.Context.CacheClient = app_obj.GetRedisClient()
	return

}

func (r *ServiceExport) List(args *wrappers.ArgumentsExportList) (res wrappers.ResultExportList, err error) {
	res = wrappers.ResultExportList{List: []wrappers.ExportShowObject{}}
	dao := dao_impl.NewDaoExport(r.Context)
	list, err := dao.GetListByUser(args.User.UserId, args.Limit)
	if err != nil {
		return
	}

	// 更新超时数据
	err = r.UpdateExpireData(dao, list)
	if err != nil {
		return
	}
	var dt wrappers.ExportShowObject
	for _, value := range *list {
		dt = wrappers.ExportShowObject{
			Hid:            value.Hid,
			Name:           value.Name,
			Progress:       value.Progress,
			Status:         value.Status,
			Type:           value.Type,
			DownloadLink:   value.DownloadLink,
			CreateAtString: utils.ShowDateTime(value.CreatedAt.Time),
		}
		res.List = append(res.List, dt)
	}
	return
}
func (r *ServiceExport) Cancel(args *wrappers.ArgumentsExportCancel) (res wrappers.ResultExportCancel, err error) {

	return
}
func (r *ServiceExport) Init(args *wrappers.ArgumentsExportInit) (res wrappers.ResultExportInit, err error) {
	ex := export.NewServiceActExport(r.Context)
	res, err = ex.SetArgument(args).Run()
	return
}

func (r *ServiceExport) Progress(args *wrappers.ArgumentsExportProgress) (res wrappers.ResultExportProgress, err error) {
	res = wrappers.ResultExportProgress{Data: map[string]int{}}
	dao := dao_impl.NewDaoExport(r.Context)
	list, err := dao.Progress(args)
	if err != nil {
		return
	}
	for _, value := range *list {
		progress, _ := r.Context.CacheClient.Get(r.Context.GinContext.Request.Context(),value.GetCacheKey()).Int()
		if progress <= 0 {
			res.Data[value.Hid] = value.Progress
		} else {
			res.Data[value.Hid] = progress
		}
	}
	// 更新超时数据
	err = r.UpdateExpireData(dao, list)
	return
}

// 更新超时数据
func (r *ServiceExport) UpdateExpireData(dao *dao_impl.DaoExport, list *[]models.ZExportData) (err error) {
	hIds := &[]string{}
	for _, value := range *list {
		// 如果导出任务 如果一天都未结束 那么就判定超时 ，退出任务
		if time.Now().Unix()-value.UpdatedAt.Unix() > 1*86400 {
			*hIds = append(*hIds, value.Hid)
		}
	}

	if dao == nil {
		dao = dao_impl.NewDaoExport(r.Context)
	}
	err = dao.UpdateByHIds(map[string]interface{}{
		"status": web.ExportExpire,
	}, hIds)
	return
}
