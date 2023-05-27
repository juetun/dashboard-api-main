package wrapper_outernet

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
)

type (
	ArgData struct {
		DocKey            string                 `json:"doc_key" form:"doc_key"`
		GetDataTypeCommon base.GetDataTypeCommon `json:"-" form:"-"`
	}

	ResultData struct {
		HavError   bool   `json:"hav_error"`
		ErrorMsg   string `json:"error_msg"`
		DocName    string `json:"doc_name"`
		DocContent string `json:"doc_content"`
	}

	ArgTree struct {
		TopId             int64                  `json:"top_id" form:"top_id"`
		BizCode           string                 `json:"biz_code" form:"biz_code"`
		CurrentId         int64                  `json:"current_id" form:"current_id"`
		DocKey            string                 `json:"doc_key" form:"doc_key"`
		GetDataTypeCommon base.GetDataTypeCommon `json:"-" form:"-"`
	}
	ResultTree struct {
		Menu        []*ResultFormPage `json:"menu"`
		Breadcrumbs []BreadcrumbItem  `json:"breadcrumbs"`
		DocContent  string            `json:"doc_content"`
		NotExists   bool              `json:"not_exists"`
		ErrorMsg    string            `json:"error_msg"`
	}
	ResultHelpTreeItemMenu struct {
		models.HelpDocumentRelate
		Active bool `json:"active"`
	}
	ResultFormPage struct {
		Id         int64                `json:"id"`
		Label      string               `json:"label"`
		Name       string               `json:"name,omitempty"`
		BizCode    string               `json:"biz_code"`
		TreeName   string               `json:"treeName"`
		IsLeafNode uint8                `json:"is_leaf_node,omitempty"`
		Active     ResultFormPageActive `json:"active"`
		IsActive   bool                 `json:"is_active"`
		OpenNames  []string             `json:"openNames"`
		Children   []*ResultFormPage    `json:"showChildren"`
	}
	BreadcrumbItem struct {
		Label string `json:"label"`
	}
	ResultFormPageActive struct {
		ActiveName string `json:"activeName"`
	}
)

func (r *ArgData) Default(ctx *base.Context) (err error) {

	return
}

func (r *ResultFormPage) ParseFromHelpDocumentRelateCache(data *models.HelpDocumentRelateCache) {
	r.Id = data.Id
	r.Label = data.Label
	r.Name = data.DocKey
	r.IsLeafNode = data.IsLeafNode
	r.BizCode = data.BizCode
}

func NewResultTree() (res *ResultTree) {
	res = &ResultTree{
		Menu:        []*ResultFormPage{},
		Breadcrumbs: make([]BreadcrumbItem, 0, 20),
		DocContent:  "",
	}
	return
}

func NewResultFormPage() (res *ResultFormPage) {
	res = &ResultFormPage{
		OpenNames: []string{},
		Children:  []*ResultFormPage{},
	}
	return
}

func (r *ArgTree) Default(context *base.Context) (err error) {

	return
}

func (r *ResultHelpTreeItemMenu) SetHelpDocumentRelate(relate *models.HelpDocumentRelate) {
	r.HelpDocumentRelate = *relate
}
