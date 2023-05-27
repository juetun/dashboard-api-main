package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/outernets/outernet_impl"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet, func(r *gin.Engine, urlPrefix string) {

		router := r.Group(urlPrefix + "/help")
		help := outernet_impl.NewConOuterNetsHelp()
		router.GET("/tree", help.Tree) //帮助文档树形结构
		router.GET("/data", help.Data) //帮助文档
	})
}
