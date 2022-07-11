package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/admins/admin_impl"
)

//系统文档接口，如用户协议、帮助功能等
func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		router := r.Group(urlPrefix + "/help")
		impl := admin_impl.NewConsoleHelp()
		router.GET("/list", impl.HelpList)     //帮助文档列表
		router.GET("/detail", impl.HelpDetail) //帮助内容详情
		router.POST("/edit", impl.HelpEdit)    //帮助内容详情

		router.GET("/trees", impl.HelpTree)    //帮助文档的树结构
		router.GET("/edit_node", impl.TreeEditNode) //编辑树节点
	})
}
