package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/admins/admin_impl"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		controller := admin_impl.NewConAdminUser()
		h := r.Group(urlPrefix + "/permit")
		h.POST("/admin_user", controller.AdminUser)
		h.POST("/admin_user_edit", controller.AdminUserEdit)
		h.POST("/admin_user_update_with_column", controller.AdminUserUpdateWithColumn)
		h.POST("/admin_user_delete", controller.AdminUserDelete)
		return
	})
}
