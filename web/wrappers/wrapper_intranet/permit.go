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
		NeedLogin  bool `json:"need_login"`  // 是否需要登录
		NeedSign   bool `json:"need_sign"`   // 是否需要签名验证
		HavePermit bool `json:"have_permit"` // 是否有权限
	}
)

func (r *ArgGetImportPermit) Default(c *gin.Context) (err error) {

	return
}
