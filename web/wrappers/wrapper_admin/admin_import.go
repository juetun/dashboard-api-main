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
	ArgMenuImport struct {
		app_param.RequestUser
		response.PageQuery
		HaveSelect string `json:"have_select" form:"have_select"`
		MenuId     int    `json:"menu_id" form:"menu_id"`
		AppName    string `json:"app_name" form:"app_name"`
		UrlPath    string `json:"url_path" form:"url_path"`
	}
	ResultMenuImport struct {
		*response.Pager
	}
)

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
