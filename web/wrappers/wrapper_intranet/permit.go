package wrapper_intranet

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type (
	ArgGetUerImportPermit struct {
		UHid    string          `json:"u_hid" form:"u_hid"`
		Uris    string          `json:"uris" form:"uris"`
		UrlInfo []UerImportItem `json:"url_info" form:"url_info"`
	}
	UerImportItem struct {
		Uk     string `json:"uk,omitempty"` // 唯一ID
		Uri    string `json:"uri,omitempty"`
		App    string `json:"app,omitempty"`
		Method string `json:"method"`
	}
	ResultGetUerImportPermit struct {
		MapHavePermit map[string]bool `json:"map_have_permit"`
		IsSuper       bool            `json:"is_super"`
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

func NewResultGetUerImportPermit(arg *ArgGetUerImportPermit) (res *ResultGetUerImportPermit) {
	res = &ResultGetUerImportPermit{IsSuper: false, MapHavePermit: make(map[string]bool, len(arg.UrlInfo))}
	for _, i2 := range arg.UrlInfo {
		res.MapHavePermit[i2.ToUk()] = false
	}
	return
}
func (r *UerImportItem) ToUk() (res string) {
	if r.Uk != "" {
		res = r.Uk
		return
	}
	res = fmt.Sprintf("%s%s%s", r.Method, r.App, r.Uri)
	return
}
func (r *ArgGetUerImportPermit) Default(c *gin.Context) (err error) {
	if r.UHid == "" {
		err = fmt.Errorf("请选择你要查看权限的用户")
		return
	}

	if r.Uris != "" {
		if err = json.Unmarshal([]byte(r.Uris), &r.UrlInfo); err != nil {
			err = fmt.Errorf("您传递的参数uris格式不正确")
			return
		}
	}

	if len(r.UrlInfo) == 0 {
		err = fmt.Errorf("请选择你要判断的接口路径")
		return
	}

	return
}

func (r *ArgGetImportPermit) Default(c *gin.Context) (err error) {
	return
}
