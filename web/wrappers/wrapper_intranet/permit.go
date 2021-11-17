package wrapper_intranet

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	ArgGetUerImportPermit struct {
		UHid       string   `json:"u_hid" form:"u_hid"`
		Uris       string   `json:"uris" form:"uris"`
		UriStrings []string `json:"-" form:"-"`
	}

	ResultGetUerImportPermit struct {
	}

	ArgGetImportPermit struct {
		AppName string `json:"app_name"`
	}
	ResultGetImportPermit struct {
		RouterNotNeedSign  map[string]*RouterNotNeedItem `json:"not_sign"`  // 不需要签名验证的路由权限
		RouterNotNeedLogin map[string]*RouterNotNeedItem `json:"not_login"` // 不需要登录的路由权限
	}
	RouterNotNeedItem struct {
		GeneralPath map[string]ItemGateway `json:"general,omitempty"` // 普通路径
		RegexpPath  []ItemGateway          `json:"regexp,omitempty"`  // 按照正则匹配的路径
	}
	ItemGateway struct {
		Uri     string           `json:"url,omitempty"`
		Methods map[string]uint8 `json:"method,omitempty"`
	}
)

func (r *ArgGetUerImportPermit) Default(c *gin.Context) (err error) {
	if r.Uris == "" {
		err = fmt.Errorf("请选择你要判断的接口路径")
		return
	}
	uString := strings.Split(r.Uris, ",")
	r.UriStrings = make([]string, 0, len(uString))
	for _, item := range uString {
		if item == "" {
			continue
		}
		r.UriStrings = append(r.UriStrings, item)
	}

	if len(r.UriStrings) == 0 {
		err = fmt.Errorf("请选择你要判断的接口路径")
		return
	}

	if r.UHid == "" {
		err = fmt.Errorf("请选择你要查看权限的用户")
		return
	}
	return
}

func (r *ArgGetImportPermit) Default(c *gin.Context) (err error) {
	return
}
