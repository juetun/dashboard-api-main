/**
* @Author:changjiang
* @Description:
* @File:excel_operate
* @Version: 1.0.0
* @Date 2020/6/10 4:45 下午
 */
package export

import (
	"math"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ExcelOperate struct {
	Context      *base.Context
	FileHandler  *excelize.File
	NowSheetName string                `json:"now_sheet_name"`
	PathFileName string                `json:"path_file_name"` // excel文件的名称(linux的绝对路径)
	FileNameBase string                `json:"file_name_base"` // 标准文件名称
	SheetNames   []string              `json:"sheet_name"`
	SheetNameMap map[string]ExcelSheet `json:"sheet_name_map"`
	TotalLine    int                   `json:"total_line"` // 当前文件已写入多少行

}
type ExcelSheet struct {
	Header []wrappers.ExcelHeader
	Index  int
}

func NewExcelOperate(context *base.Context) (r *ExcelOperate) {
	r = &ExcelOperate{
		Context: context,
		SheetNames: []string{
			"sheet1",
		},
		SheetNameMap: map[string]ExcelSheet{
			"sheet1": {
				Header: []wrappers.ExcelHeader{},
				Index:  0,
			},
		},
		PathFileName:  "export.xlsx",
		TotalLine: 0,
	}
	r.FileHandler = excelize.NewFile()
	return
}

func (r *ExcelOperate) SetHeader(sheetName string, header []wrappers.ExcelHeader) (p *ExcelOperate) {
	r.NowSheetName = sheetName
	sheet := r.SheetNameMap[r.NowSheetName]
	sheet.Header = header
	r.SheetNameMap[r.NowSheetName] = sheet
	if r.judgeNeedWriterHeader(&header) {
		var data = map[string]interface{}{}
		for _, value := range header {
			data[value.ColumnKey] = value.Label
		}
		r.WriterData(r.NowSheetName, data)
	}
	return r
}

func (r *ExcelOperate) Init() (p *ExcelOperate) {
	return r
}
func (r *ExcelOperate) Close() (err error) {
	for _, sheetName := range r.SheetNames {
		r.FileHandler.SetActiveSheet(r.SheetNameMap[sheetName].Index)
	}
	err = r.FileHandler.SaveAs(r.PathFileName)
	return
}

// 获取Excel需要的Key 如 A1 A2 B1 B2
func (r *ExcelOperate) getLocKey(line, columnIndex int) (res string) {
	var collection = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := len(collection)
	cCount := int(math.Floor(float64(columnIndex) / float64(length)))
	y := columnIndex % length
	if cCount > 0 { // 当前导出最多支持 26*26列 （26个英文字符）
		res += string(collection[cCount]) + string(collection[y])
	} else {
		res += string(collection[y])
	}
	res += strconv.Itoa(line)
	return
}

func (r *ExcelOperate) WriterData(sheetName string, data map[string]interface{}) {
	r.NowSheetName = sheetName
	r.TotalLine++
	var locKey string
	for k, header := range r.SheetNameMap[r.NowSheetName].Header {
		var v interface{}
		if _, ok := data[header.ColumnKey]; ok {
			v = data[header.ColumnKey]
		} else { // 如果数据不存在，则用默认数据占位
			v = ""
		}
		locKey = r.getLocKey(r.TotalLine, k)
		r.FileHandler.SetCellValue(r.NowSheetName, locKey, v)
	}
}

func (r *ExcelOperate) SetSheet(sheetNames *[]string) (p *ExcelOperate) {
	r.SheetNames = *sheetNames
	for key, value := range r.SheetNames {
		if value == "" {
			value = "sheet" + strconv.Itoa(key)
			r.SheetNames[key] = value
		}
		if _, ok := r.SheetNameMap[value]; !ok {
			sheet := r.SheetNameMap[value]
			sheet.Index = r.FileHandler.NewSheet(value)
			r.SheetNameMap[value] = sheet
		}
	}
	return r
}
func (r *ExcelOperate) judgeNeedWriterHeader(data *[]wrappers.ExcelHeader) (res bool) {
	res = false
	for _, value := range *data {
		if value.Label != "" {
			res = true
			return
		}
	}
	return
}
