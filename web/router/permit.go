/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:04 上午
 */
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_start"
	"github.com/juetun/dashboard-api-main/web/controllers/impl/permit_impl"
)

func init() {
	app_start.HandleFunc = append(app_start.HandleFunc, func(r *gin.Engine, urlPrefix string) {
		controller := permit_impl.NewControllerPermit()
		h := r.Group(urlPrefix + "/permit")

		h.POST("/admin_user", controller.AdminUser)
		h.POST("/admin_group", controller.AdminGroup)
		h.POST("/admin_menu", controller.AdminMenu)         // 菜单列表
		h.POST("/admin_menu_add", controller.MenuAdd)       // 添加菜单
		h.POST("/admin_menu_save", controller.MenuSave)     // 修改菜单
		h.POST("/admin_menu_delete", controller.MenuDelete) // 删除菜单

		// 权限菜单列表
		h.GET("/menu", controller.Menu)
		h.POST("/flag", controller.Flag) // 指定链接是否有权限
	})
}
