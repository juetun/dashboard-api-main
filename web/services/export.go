/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:34 上午
 */
package services

import (
	"time"

	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/basic/utils"
	"github.com/juetun/dashboard-api-main/web"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/pojos"
	"github.com/juetun/dashboard-api-main/web/services/export"
)

type ServiceExport struct {
	base.ServiceBase
}

func NewServiceExport(context ...*base.Context) (p *ServiceExport) {
	p = &ServiceExport{}
	p.SetContext(context)
	p.Context.CacheClient = app_obj.GetRedisClient()
	return

}

func (r *ServiceExport) List(args *pojos.ArgumentsExportList) (res pojos.ResultExportList, err error) {
	res = pojos.ResultExportList{List: []pojos.ExportShowObject{}}
	dao := daos.NewDaoExport(r.Context)
	list, err := dao.GetListByUser(args.User.UserId, args.Limit)
	if err != nil {
		return
	}

	// 更新超时数据
	err = r.UpdateExpireData(dao, list)
	if err != nil {
		return
	}
	var dt pojos.ExportShowObject
	for _, value := range *list {
		dt = pojos.ExportShowObject{
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
func (r *ServiceExport) Cancel(args *pojos.ArgumentsExportCancel) (res pojos.ResultExportCancel, err error) {

	return
}
func (r *ServiceExport) Init(args *pojos.ArgumentsExportInit) (res pojos.ResultExportInit, err error) {
	ex := export.NewServiceActExport(r.Context)
	res, err = ex.SetArgument(args).Run()
	return
}

func (r *ServiceExport) Progress(args *pojos.ArgumentsExportProgress) (res pojos.ResultExportProgress, err error) {
	res = pojos.ResultExportProgress{Data: map[string]int{}}
	dao := daos.NewDaoExport(r.Context)
	list, err := dao.Progress(args)
	if err != nil {
		return
	}
	for _, value := range *list {
		progress, _ := r.Context.CacheClient.Get(value.GetCacheKey()).Int()
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
func (r *ServiceExport) UpdateExpireData(dao *daos.DaoExport, list *[]models.ZExportData) (err error) {
	hIds := &[]string{}
	for _, value := range *list {
		// 如果导出任务 如果一天都未结束 那么就判定超时 ，退出任务
		if time.Now().Unix()-value.UpdatedAt.Unix() > 1*86400 {
			*hIds = append(*hIds, value.Hid)
		}
	}

	if dao == nil {
		dao = daos.NewDaoExport(r.Context)
	}
	err = dao.UpdateByHIds(map[string]interface{}{
		"status": web.ExportExpire,
	}, hIds)
	return
}
