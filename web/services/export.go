/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:34 上午
 */
package services

import (
	"github.com/juetun/base-wrapper/lib/base"
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
	return

}

func (r *ServiceExport) List(args *pojos.ArgumentsExportList) (res pojos.ResultExportList, err error) {
	res = pojos.ResultExportList{List: []models.ZExportData{}}
	dao := daos.NewDaoExport(r.Context)
	list, err := dao.GetListByUser(args.User.UserId, args.Limit)
	if err != nil {
		return
	}
	res.List = *list
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
		res.Data[value.Hid] = value.Progress
	}
	return
}
