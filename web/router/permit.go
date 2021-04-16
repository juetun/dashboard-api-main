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
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/con_impl"
)

func init() {
	app_start.HandleFunc = append(app_start.HandleFunc, func(r *gin.Engine, urlPrefix string) {
		controller := con_impl.NewControllerPermit()
		h := r.Group(urlPrefix + "/permit")

		h.POST("/admin_user", controller.AdminUser)
		h.POST("/admin_user_add", controller.AdminUserAdd)
		h.POST("/admin_user_delete", controller.AdminUserDelete)
		h.POST("/admin_user_group_release", controller.AdminUserGroupRelease)
		h.POST("/admin_user_group_add", controller.AdminUserGroupAdd)

		h.POST("/admin_group", controller.AdminGroup)
		h.POST("/edit_admin_group", controller.AdminGroupEdit)
		h.POST("/admin_group_delete", controller.AdminGroupDelete)
		h.POST("/admin_menu", controller.AdminMenu) // 菜单列表
		h.GET("/get_menu", controller.GetMenu)      // 获取菜单信息

		h.POST("/admin_menu_search", controller.AdminMenuSearch)
		h.POST("/admin_menu_add", controller.MenuAdd)       // 添加菜单
		h.POST("/admin_menu_save", controller.MenuSave)     // 修改菜单
		h.POST("/admin_menu_delete", controller.MenuDelete) // 删除菜单
		h.POST("/get_import", controller.GetImport)         // 获取接口列表
		// 权限菜单列表
		h.GET("/menu", controller.Menu)
		h.POST("/flag", controller.Flag) // 指定链接是否有权限
	})
}
