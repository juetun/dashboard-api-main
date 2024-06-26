// Package admin
/**
* @Author:ChangJiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:04 上午
 */
package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/admins/admin_impl"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		permitRoute(r, urlPrefix) //permit
		importRoute(r, urlPrefix) //import
	})
}

func importRoute(r *gin.Engine, urlPrefix string) {
	h := r.Group(urlPrefix + "/permit")
	controller := admin_impl.NewAdminConPermitImportImpl()
	h.POST("/user_page_import", controller.UserPageImport) //页面的接口权限列表
	h.POST("/get_import", controller.GetImport)            // 获取接口列表
	h.POST("/menu_import", controller.MenuImport)          // 获取接口列表
	h.POST("/edit_import", controller.EditImport)          //接口编辑和添加
	h.POST("/menu_import_set", controller.MenuImportSet)   // 界面接口设置
	h.POST("/import_list", controller.ImportList)
	h.POST("/update_import_value", controller.UpdateImportValue)
	h.DELETE("/delete_import/:id", controller.DeleteImport)

	// 根据菜单号 获取页面的接口ID
	h.GET("/get_import_by_menu_id", controller.GetImportByMenuId)
}

//menu
func permitRoute(r *gin.Engine, urlPrefix string) {
	controller := admin_impl.NewControllerPermit()
	h := r.Group(urlPrefix + "/permit")
	// 权限菜单列表
	h.GET("/menu", controller.Menu)
	h.GET("get_system",controller.GetSystem)


	h.POST("/admin_user_group_release", controller.AdminUserGroupRelease)
	h.POST("/admin_user_group_add", controller.AdminUserGroupAdd) // 用户组添加管理员

	h.POST("/admin_group", controller.AdminGroup) // 管理员用户组列表查询
	h.POST("/admin_set_permit", controller.AdminSetPermit)
	h.POST("/edit_admin_group", controller.AdminGroupEdit) // 编辑用户组
	h.POST("/admin_group_delete", controller.AdminGroupDelete)

	h.POST("/admin_menu", controller.AdminMenu)               // 菜单列表
	h.POST("/menu_with_check", controller.AdminMenuWithCheck) // 菜单列表 用于设置管理员组权限使用
	h.GET("/get_menu", controller.GetMenu)                    // 获取菜单信息

	h.POST("/admin_menu_search", controller.AdminMenuSearch)
	h.POST("/admin_menu_add", controller.MenuAdd)       // 添加菜单
	h.POST("/admin_menu_save", controller.MenuSave)     // 修改菜单
	h.POST("/admin_menu_delete", controller.MenuDelete) // 删除菜单

	// 获取服务列表
	h.GET("/get_app_config", controller.GetAppConfig)
	return
}
