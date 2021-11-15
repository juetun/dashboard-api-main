package wrapper_intranet

import (
	"github.com/gin-gonic/gin"
)

type (
	ArgGetImportPermit struct {
		UserHid  string `json:"user_hid"` // 当前登录用户ID
		Token    string `json:"token"`
		Method   string `json:"method"`
		AppName  string `json:"app_name"`
		PathType string `json:"path_type"`
		Uri      string `json:"uri"`
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
		Uri     string     `json:"url,omitempty"`
		Methods map[string]uint8 `json:"method,omitempty"`
	}
)

func (r *ArgGetImportPermit) Default(c *gin.Context) (err error) {
	return
}
