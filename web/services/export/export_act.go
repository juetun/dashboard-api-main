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
	"math"
	"os"
	"strconv"
	"strings"
	"time"

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

	// 设置导出到时的时间为3天
	ctx, cancel := context.WithTimeout(context.Background(), 3*86400*time.Second)

	go r.Act(ctx)

	// 如果是外部取消任务，任务结束状态实现
	if r.IsFinish {
		r.Context.Log.Error(
			map[string]string{
				"message": fmt.Sprintf("外部命令中止导出任务(ID:%d) ", r.model.Id),
			},
		)
		cancel()
	}

}

// 将导出生成的文件上传到阿里云
func (r *AsyncExport) uploadFileToTarget(excel *ExcelOperate, exportData *models.ZExportData) (err error) {
	fileUpload := NewNewFileUpload()
	fileUpload.SetFile(excel.PathFileName).
		Run()
	err = fileUpload.Err
	if err != nil {
		r.Context.Log.Logger.Errorln(
			"message", fmt.Sprintf("文件（%s）上传文件到阿里云失败 ", excel.PathFileName),
			"content:", err.Error())
		return
	}
	exportData.DownloadLink = fileUpload.DownloadUrl
	exportData.Domain = fileUpload.BucketUrl
	exportData.FilePath = fileUpload.ObjectName
	exportData.Name = fileUpload.FileName
	return
}

// 更新导出的进度
func (r *AsyncExport) UpdateProgress(progress int) {

	// 正常进度更新不会超过98 超过98默认就是完成 ，完成的逻辑单独处理
	if progress >= web.RunProgressMax {
		progress = web.RunProgressMax
	}
	r.model.Progress = progress

	// 更新Redis进度
	r.UpdateRedisProgressValue()
}

// 更新Redis进度
func (r *AsyncExport) UpdateRedisProgressValue() {
	r.Context.CacheClient.Set(r.model.GetCacheKey(), r.model.Progress, 86400*time.Second)
}

// 获取生成数据完成后的动作，如：将Excel生成文件，更新导出进度等
func (r *AsyncExport) getDataFinishAct(excel *ExcelOperate) {

	r.UpdateProgress(web.RunProgressMax) // 获取数据完成后更新进度

	var err error
	r.Context.Log.Info(
		map[string]string{
			"message": fmt.Sprintf("文件（ %s）生成成功 ", excel.PathFileName),
		},
	)
	excel.Close() // excel文件生成

	// 上传EXCEL文件到指定路径
	err = r.uploadFileToTarget(excel, &r.model)
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message":  fmt.Sprintf("上传文件(%s)到指定路径错误 ", excel.PathFileName),
			"content:": err.Error(),
		})
		r.model.Status = web.ExportFailure
	} else {
		r.model.Status = web.ExportSuccess
	}

	// 修改任务进度
	r.model.Progress = 100

	// 更新Redis任务进度
	r.UpdateRedisProgressValue()

	err = daos.NewDaoExport(r.Context).Update(&r.model)
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message":  "update export progress to database is error ",
			"content:": err.Error(),
		})
	}
}

func (r *AsyncExport) PathFileName() (pathFileName string) {
	pathFileName = r.getSysTmp() + "/" + r.args.FileName + "_" + strconv.Itoa(r.model.Id) + ".xlsx"
	return
}

// 获得操作系统的临时目录文件夹
func (r *AsyncExport) getSysTmp() (res string) {
	res = strings.TrimSuffix(os.TempDir(), "/")
	r.Context.Log.Info(map[string]string{
		"临时目录:": res,
	})
	return
}

// 生成excel sheet名称,没有就使用默认
func (r *AsyncExport) defaultSheetName() (sheetNames *[]string) {
	sheetNames = &[]string{}
	for key, value := range r.args.Program {
		if value.SheetName == "" {
			value.SheetName = "sheet" + strconv.Itoa(key+1)
			r.args.Program[key] = value
		}
		*sheetNames = append(*sheetNames, value.SheetName)
	}
	return
}
func (r *AsyncExport) Act(ctx context.Context) {

	defer func() {
		<-ctx.Done()
	}()

	r.work()
}

func (r *AsyncExport) work() {

	// 设置导出任务开始的进度
	r.UpdateProgress(web.RunProgressInit)

	excel := NewExcelOperate(r.Context)

	defer r.getDataFinishAct(excel)

	// 获取生成文件的名字
	excel.PathFileName = r.PathFileName()

	sheetNames := r.defaultSheetName()

	excel = excel.SetSheet(sheetNames) // 生成Excel的sheet对象

	for _, value := range r.args.Program {

		value.HttpHeader = r.args.HttpHeader

		// 设置每个sheet的第一行数据及字段映射信息
		excel = excel.SetHeader(value.SheetName, value.Header).
			Init()
		r.UpdateProgress(web.RunProgressStart)

		// 处理第一页数据，目的为获取数据的格式 总条数或判断什么时候 获取数据可以结束
		firstPage, err := r.doOnePage(excel, &value)

		r.UpdateProgress(web.RunProgressFistPage)

		if err != nil {
			return
		}
		// 获取除第一页外的数据
		err = r.doOtherPage(excel, &firstPage, &value)
		if err != nil {
			return
		}
	}
}

// 将本次HTTP请求的header传递到后台 调用获取数据接口
// func (r *AsyncExport) orgRequestHttpHeader(header http.Header) (resHeader http.Header) {
// 	resHeader = header
// 	return
// }
func (r *AsyncExport) doOtherPage(excel *ExcelOperate, page *Pager, artSheet *pojos.ArgumentExportSheet) (err error) {
	var pageIndex = 1
	var progressStep float64 = 1
	for {
		if page.IsFinish { // 如果获取数据未结束
			return
		}
		pageIndex++
		// 如果没有获取到数据了或者取到10万页了，则说明取完了
		if page.List == nil || len(page.List) < 1 || pageIndex > 100000 {
			break
		}

		artSheet.Query["page_no"] = r.getPageNoFormat(artSheet.Query["page_no"], page.PageNo+1)

		// 更新导出进度
		r.UpdateProgress(int(math.Floor(float64(r.model.Progress) + progressStep)))

		// 处理第一页数据，目的为获取数据的格式 总条数或判断什么时候 获取数据可以结束
		*page, err = r.doOnePage(excel, artSheet)
	}

	return
}

func (r *AsyncExport) getPageNoFormat(pageInter interface{}, pageNo int) (res interface{}) {
	switch pageInter.(type) {
	case int:
		res = pageNo
		break
	case string:
		res = strconv.Itoa(pageNo)
		break
	default:
		res = pageNo
	}
	return
}
func (r *AsyncExport) doOnePage(excel *ExcelOperate, artSheet *pojos.ArgumentExportSheet) (pageData Pager, err error) {
	// 获取第一页数据
	pageData, err = r.getData(artSheet)
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message":  "数据获取异常",
			"content:": err.Error()})
		return
	}
	r.Context.Log.Info(map[string]string{"desc": fmt.Sprintf("first page return:%v", pageData)})
	r.WritePageData(&pageData, excel)
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
	IsFinish   bool                     `json:"is_not_finish"` // 用于特殊数据判断
	PageNo     int                      `json:"page_no"`
	PageSize   int                      `json:"page_size"`
	List       []map[string]interface{} `json:"list"`
	TotalCount int                      `json:"total_count"`
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
				r.Context.Log.Logger.Errorln("message", fmt.Sprintf("the app_name %s(%s) config domain is err ", value.Name, value.Key), "content:", err.Error())
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
		r.Context.Log.Logger.Errorln("message", fmt.Sprintf("the params: %s ", string(*args)))
		request = RequestObject{
			ReSendTimes:        3,
			Uri:                "",
			ContentType:        "application/json",
			Body:               bytes.NewBuffer(*args),
			RequestMethod:      "POST",
			HttpRequestContent: artSheet.HttpRequestContent,
		}
		request.Uri = Uri
		if err != nil {
			r.Context.Log.Logger.Errorln("message", " getData", "content:", err.Error())
			return
		}
		break
	case "GET":
		request = RequestObject{
			ReSendTimes:        3,
			Uri:                "",
			ContentType:        "application/json",
			Body:               nil,
			RequestMethod:      "GET",
			Query:              artSheet.Query,
			HttpRequestContent: artSheet.HttpRequestContent,
		}
		request.Uri = Uri
		break
	default:
		err = fmt.Errorf("当前不支持%s方法", artSheet.HttpMethod)
		r.Context.Log.Error(map[string]string{
			"desc": err.Error(),
		}, )
		return
	}

	result, err := httpRequest.Send(&request)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", " getData", "content:", err.Error())
		return
	}
	var dt struct {
		Code int    `json:"code"`
		Data Pager  `json:"data"`
		Msg  string `json:"msg"`
	}

	// 如果没有获取到数据
	if len(*result) == 0 {
		r.Context.Log.Error(map[string]string{
			"message": "请求接口返回数据为空",
			"request": fmt.Sprintf("%v", request),
		})
		return
	}
	err = json.Unmarshal(*result, &dt)
	resString := string(*result)
	if err != nil {
		r.Context.Log.Logger.Errorln("message", " Unmarshal JSON  Exception", "content:", resString)
	} else {
		r.Context.Log.Logger.Infoln("message", "request return ", "content:", resString)
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
