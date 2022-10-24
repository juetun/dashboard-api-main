package wrapper_admin

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"

	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/library/common/app_param"
)

const (
	ArgMenuImportHaveSelectYes = "yes"
)

type (
	ArgPageImport struct {
		app_param.RequestUser
	}
	ResultPageImport struct {
		CommonImport map[string]uint8 `json:"common_import"` //公共接口列表
		PageImport   map[string]uint8 `json:"page_import"`   //页面接口列表
		SubPageMenu  map[string]uint8 `json:"sub_page_menu"` //子页面KEY
	}
	ArgMenuImport struct {
		app_param.RequestUser
		response.PageQuery
		HaveSelect string `json:"have_select" form:"have_select"`
		MenuId     int    `json:"menu_id" form:"menu_id"`
		AppName    string `json:"app_name" form:"app_name"`
		UrlPath    string `json:"url_path" form:"url_path"`
	}
	ResultMenuImport struct {
		response.Pager
	}
)

func NewResultPageImport() (res *ResultPageImport) {
	res = &ResultPageImport{
		CommonImport: map[string]uint8{},
		PageImport:   map[string]uint8{},
		SubPageMenu:  map[string]uint8{},
	}
	return
}

func (r *ArgPageImport) Default(ctx *base.Context) (err error) {
	if err = r.InitRequestUser(ctx); err != nil {
		return
	}
	return
}

func (r *ArgMenuImport) Default(c *base.Context) (err error) {
	_ = c
	if r.HaveSelect != "" {
		if r.HaveSelect != ArgMenuImportHaveSelectYes {
			err = fmt.Errorf("have_select 格式不正确")
			return
		}
	}

	return
}
