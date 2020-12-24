/**
* 导出任务单独逻辑处理
* @Author:changjiang
* @Description:
* @File:init
* @Version: 1.0.0
* @Date 2020/6/10 3:58 下午
 */
package export

import (
	"encoding/json"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web"
	"github.com/juetun/dashboard-api-main/web/daos/dao_impl"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ServiceActExport struct {
	base.ServiceBase
	args       *wrappers.ArgumentsExportInit
	argsString string
}

func (r *ServiceActExport) Run() (res wrappers.ResultExportInit, err error) {
	res = wrappers.ResultExportInit{}

	// 在数据库中添加任务数据
	dt, err := r.InsertDataToDb()
	if err != nil {
		return
	}

	// 异步执行导出任务动作,后续修改成rpc生成实现
	NewAsyncExport(r.Context).
		SetExportData(r.args, *dt).
		Run()

	res.ExportHid = dt.Hid
	return
}

// 将数据内容添加到数据库
func (r *ServiceActExport) InsertDataToDb() (dt *models.ZExportData, err error) {
	extFileName := "xlsx"
	dt = &models.ZExportData{
		Name:          r.args.FileName + "." + extFileName,
		Progress:      web.ExportLoading,
		Type:          extFileName,
		Arguments:     r.argsString,
		DownloadLink:  "",
		CreateUserHid: r.args.User.UserId,
	}
	dao := dao_impl.NewDaoExport(r.Context)
	_, err = dao.Create(dt)
	if err != nil {
		return
	}
	return
}

func (r *ServiceActExport) SetArgument(args *wrappers.ArgumentsExportInit) (p *ServiceActExport) {
	arg, _ := json.Marshal(args)
	r.argsString = string(arg)
	r.Context.Log.Logger.Errorln("message", "export argument", "content:", r.argsString)
	r.args = args
	return r
}

func NewServiceActExport(context ...*base.Context) (p *ServiceActExport) {
	p = &ServiceActExport{}
	p.SetContext(context...)
	return
}
