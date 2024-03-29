// Package wrappers
/**
* @Author:changjiang
* @Description:
* @File:export
* @Version: 1.0.0
* @Date 2020/6/10 10:35 上午
 */
package wrappers

import (
	"github.com/juetun/base-wrapper/lib/base"
	"net/http"
	"strings"

	"github.com/juetun/library/common/app_param"
)

type ArgumentsExportList struct {
	app_param.RequestUser
	Limit int `json:"limit"`
}

func (r *ArgumentsExportList) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
}

type ResultExportList struct {
	List []ExportShowObject `json:"list"`
}
type ExportShowObject struct {
	Hid            string `gorm:"column:hid;" json:"hid"`
	Name           string `gorm:"column:name;" json:"name"`
	Progress       int    `gorm:"column:progress;" json:"progress"`
	Status         int    `gorm:"column:status;" json:"status"`
	Type           string `gorm:"column:type;" json:"type"`
	DownloadLink   string `gorm:"column:download_link;" json:"download_link"`
	CreateAtString string `json:"create_at"`
}

type ArgumentsExportCancel struct {
	app_param.RequestUser
}

func (r *ArgumentsExportCancel) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	return
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
	app_param.RequestUser
	FileName   string                `json:"file_name"` // 生成文件的名称
	Program    []ArgumentExportSheet `json:"program"`
	HttpHeader http.Header           `form:"-"  json:"http_header"`
}

func (r *ArgumentsExportInit) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.HttpHeader = c.GinContext.Request.Header
	return
}

type ResultExportInit struct {
	ExportHid string `json:"export_hid"`
}

type ArgumentsExportProgress struct {
	IDS    string `form:"ids" json:"ids"`
	UserId string `form:"-" json:"user_id"`
	app_param.RequestUser
	IdString []string `form:"-" json:"id_string"`
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

func (r *ArgumentsExportProgress) Default(c *base.Context) (err error) {
	if err = r.InitRequestUser(c); err != nil {
		return
	}
	r.InitIds()
	return
}

type ResultExportProgress struct {
	Data map[string]int `json:"data"`
}
