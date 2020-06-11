/**
* 异步导出实现
* @Author:changjiang
* @Description:
* @File:export_act
* @Version: 1.0.0
* @Date 2020/6/10 4:21 下午
 */
package export

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web"
	"github.com/juetun/dashboard-api-main/web/daos"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/pojos"
)

type AsyncExport struct {
	IsFinish bool
	Context  *base.Context
	model    models.ZExportData
	args     *pojos.ArgumentsExportInit
}

func NewAsyncExport(context *base.Context) (r *AsyncExport) {
	r = &AsyncExport{
		Context: context,
	}
	return
}
func (r *AsyncExport) SetExportData(args *pojos.ArgumentsExportInit, model models.ZExportData) (p *AsyncExport) {
	r.model = model
	r.args = args
	return r
}

func (r *AsyncExport) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	go r.Act(ctx)

	// 如果是外部取消任务，任务结束状态实现
	if !r.IsFinish {
		cancel()
	}
}

// 将导出生成的文件上传到阿里云
func (r *AsyncExport) uploadFileToTarget(excel *ExcelOperate, exportData *models.ZExportData) (err error) {
	fileUpload := NewNewFileUpload()
	fileUpload.SetFile(excel.FileName).
		Run()
	err = fileUpload.Err
	if err != nil {
		r.Context.Log.Errorln(
			"message", fmt.Sprintf("文件（ %s）上传文件到阿里云失败 ", excel.FileName),
			"content:", err.Error())
		return
	}
	exportData.DownloadLink = fileUpload.DownloadUrl
	exportData.Domain = fileUpload.Endpoint
	exportData.FilePath = fileUpload.ObjectName
	return
}

func (r *AsyncExport) Act(ctx context.Context) {

	var err error
	defer func() {
		<-ctx.Done()
	}()

	excel := NewExcelOperate(r.Context)
	defer func() {
		excel.Close() // excel文件生成

		// 上传EXCEL文件到指定路径
		err = r.uploadFileToTarget(excel, &r.model)
		if err != nil {
			r.Context.Log.Errorln("message", fmt.Sprintf("上传文件(%s)到指定路径错误 ", excel.FileName), "content:", err.Error())
			r.model.Status = web.ExportFailure
		} else {
			r.model.Status = web.ExportSuccess
		}
		r.model.Progress = 100
		err := daos.NewDaoExport(r.Context).Update(&r.model)
		if err != nil {
			r.Context.Log.Errorln("message", "update export progress to database is error ", "content:", err.Error())
		}
	}()

	excel.FileName = r.args.FileName + ".xlsx"
	sheetNames := make([]string, 0)
	for key, value := range r.args.Program {
		if value.SheetName == "" {
			value.SheetName = "sheet" + strconv.Itoa(key+1)
			r.args.Program[key] = value
		}
		sheetNames = append(sheetNames, value.SheetName)
	}
	excel = excel.SetSheet(&sheetNames) // 生成Excel的sheet

	httpHeader := r.orgRequestHttpHeader(r.args.HttpHeader)
	for _, value := range r.args.Program {
		value.HttpHeader = httpHeader
		// 设置每个sheet的第一行数据及字段映射信息
		excel = excel.SetHeader(value.SheetName, value.Header).
			Init()

		// 处理第一页数据，目的为获取数据的格式 总条数或判断什么时候 获取数据可以结束
		firstPage, _ := r.doOnePage(excel, &value)

		// 获取除第一页外的数据
		r.doOtherPage(excel, &firstPage, &value)

	}
}

// 将本次HTTP请求的header传递到后台 调用获取数据接口
func (r *AsyncExport) orgRequestHttpHeader(header http.Header) (resHeader http.Header) {
	resHeader = header
	return
}
func (r *AsyncExport) doOtherPage(excel *ExcelOperate, page *Pager, artSheet *pojos.ArgumentExportSheet) {
	if !page.IsNotFinish { // 如果获取数据未结束
		return
	}

	// 处理第一页数据，目的为获取数据的格式 总条数或判断什么时候 获取数据可以结束
	r.doOnePage(excel, artSheet)
}

func (r *AsyncExport) doOnePage(excel *ExcelOperate, artSheet *pojos.ArgumentExportSheet) (firstPageData Pager, err error) {

	// 获取第一页数据
	firstPageData, err = r.getData(artSheet)
	r.Context.Log.Infoln("first page return:", firstPageData)
	if err != nil {
		r.Context.Log.Errorln("message", "export get first page data exception ", "content:", err.Error())
		return
	}
	r.WritePageData(&firstPageData, excel)
	return

}

// 生成请求接口的JSON字符串
func (r *AsyncExport) getArgString(params map[string]interface{}) (res *[]byte) {
	res = &[]byte{}
	*res, _ = json.Marshal(params)
	return
}

// 将获取的到的数据写入excel
func (r *AsyncExport) WritePageData(page *Pager, excel *ExcelOperate) {
	for _, item := range page.List {
		excel.WriterData(excel.NowSheetName, item)
	}
}

// 支持导出的分页数据最终生成的结构
type Pager struct {
	IsNotFinish bool                     `json:"is_not_finish"` // 用于特殊数据判断
	PageNo      int                      `json:"page_no"`
	PageSize    int                      `json:"page_size"`
	List        []map[string]interface{} `json:"list"`
	TotalCount  int                      `json:"total_count"`
}

type ServerConfig struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Port   int    `json:"port"`
	Domain string `json:"domain"`
}

func (r *AsyncExport) getAccessUrl(argSheet *pojos.ArgumentExportSheet) (res string, err error) {
	var serverMessage ServerConfig

	// 获取应用的服务器位置
	apps, err := daos.NewDaoAppPath().GetAllApp()
	if err != nil {
		return
	}

	var domain string
	for _, value := range *apps {
		if value.Key == argSheet.AppName {
			domain, err = value.GetRandomDomain()
			if err != nil {
				r.Context.Log.Errorln("message", fmt.Sprintf("the app_name %s(%s) config domain is err ", value.Name, value.Key), "content:", err.Error())
				return
			}
			serverMessage = ServerConfig{
				Name:   value.Name,
				Key:    value.Key,
				Port:   value.Port,
				Domain: domain,
			}
			if serverMessage.Port == 0 {
				res = fmt.Sprintf("%s/admin_car/v1%s", serverMessage.Domain, argSheet.Uri)
			} else {
				res = fmt.Sprintf("%s:%d/admin_car/v1%s", serverMessage.Domain, serverMessage.Port, argSheet.Uri)
			}
		}
	}

	return
}

// 获取数据
func (r *AsyncExport) getData(artSheet *pojos.ArgumentExportSheet) (res Pager, err error) {
	res = Pager{}
	httpRequest := NewHttpRequest(r.Context)

	var request RequestObject
	Uri, err := r.getAccessUrl(artSheet)
	switch strings.ToUpper(artSheet.HttpMethod) {
	case "POST":
		args := r.getArgString(r.dealDefaultQuery(artSheet.Query))
		r.Context.Log.Errorln("message", fmt.Sprintf("the params: %s ", string(*args)))
		request = RequestObject{
			Uri:                "",
			ContentType:        "application/json",
			Body:               bytes.NewBuffer(*args),
			RequestMethod:      "POST",
			HttpRequestContent: artSheet.HttpRequestContent,
		}
		request.Uri = Uri
		if err != nil {
			r.Context.Log.Errorln("message", " getData", "content:", err.Error())
			return
		}
		break
	case "GET":
		request = RequestObject{
			Uri:                "",
			ContentType:        "application/json",
			Body:               nil,
			RequestMethod:      "GET",
			Query:              artSheet.Query,
			HttpRequestContent: artSheet.HttpRequestContent,
		}
		request.Uri = Uri

		break
	}

	result, err := httpRequest.Send(&request)
	if err != nil {
		r.Context.Log.Errorln("message", " getData", "content:", err.Error())
		return
	}
	var dt struct {
		Code int    `json:"code"`
		Data Pager  `json:"data"`
		Msg  string `json:"msg"`
	}
	err = json.Unmarshal(*result, &dt)
	resString := string(*result)
	if err != nil {
		r.Context.Log.Errorln("message", " Unmarshal JSON  Exception", "content:", resString)
	} else {
		r.Context.Log.Infoln("message", "request return ", "content:", resString)
	}
	res = dt.Data
	return
}

func (r *AsyncExport) changeInt(v interface{}) (res int, err error) {
	switch v.(type) {
	case string:
		tp := v.(string)
		res, _ = strconv.Atoi(tp)
		break
	case int:
		res, _ = v.(int)
		break
	}

	return
}
func (r *AsyncExport) dealDefaultQuery(param map[string]interface{}) (res map[string]interface{}) {
	res = map[string]interface{}{}
	for key, value := range param {
		switch key {
		case "page_no":
			pageNO, _ := r.changeInt(value)
			if pageNO < 1 {
				pageNO = 1
			}
			res["page_no"] = strconv.Itoa(pageNO)
			break
		case "page_size":
			pageSize, _ := r.changeInt(value)
			if pageSize < 1 {
				pageSize = 500
			}
			res["page_size"] = strconv.Itoa(pageSize)
			break
		default:
			res[key] = value
		}
	}
	return
}
