/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:35 上午
 */
package pojos

import (
	"net/http"
	"strings"

	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/dashboard-api-main/web/models"
)

type ArgumentsExportList struct {
	User  app_obj.JwtUserMessage `form:"-" json:"user"`
	Limit int                    `json:"limit"`
}
type ResultExportList struct {
	List []models.ZExportData `json:"list"`
}

type ArgumentsExportCancel struct {
	User app_obj.JwtUserMessage `form:"-" json:"user"`
}
type ResultExportCancel struct {
}

type ExcelHeader struct {
	Label     string `json:"label"`      // 表头的第一行显示的文字
	ColumnKey string `json:"column_key"` // 每列对应的 接口返回数据的KEY （如:  [{name:"a",email:""},{name:"b",email:""}] ColumnKey="name" )
}
type ArgumentExportSheet struct {
	AppName    string                 `json:"app_name"`
	Uri        string                 `json:"uri"`
	SheetName  string                 `json:"sheet_name"`
	Header     []ExcelHeader          `json:"header"`
	Query      map[string]interface{} `json:"query"` // URL或POST 请求表单参数
	HttpMethod string                 `json:"http_method"`
	HttpRequestContent
}
type HttpRequestContent struct {
	HttpHeader http.Header `json:"http_header"`
}
type ArgumentsExportInit struct {
	User       app_obj.JwtUserMessage `form:"-" json:"user"`
	FileName   string                 `json:"file_name"` // 生成文件的名称
	Program    []ArgumentExportSheet  `json:"program"`
	HttpHeader http.Header            `form:"-"  json:"http_header"`
}

type ResultExportInit struct {
	ExportHid string `json:"export_hid"`
}

type ArgumentsExportProgress struct {
	IDS      string                 `form:"ids" json:"ids"`
	UserId   string                 `form:"-" json:"user_id"`
	User     app_obj.JwtUserMessage `form:"-" json:"user"`
	IdString []string               `form:"-" json:"id_string"`
}

func (r *ArgumentsExportProgress) InitIds() {
	if r.IdString == nil {
		r.IdString = make([]string, 0)
	}
	if r.IDS != "" {
		ids := strings.Split(r.IDS, ",")
		for _, value := range ids {
			if value == "" {
				continue
			}
			r.IdString = append(r.IdString, value)
		}
	}
}

type ResultExportProgress struct {
	Data map[string]int `json:"data"`
}
