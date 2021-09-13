package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/admins/admin_impl"
	"github.com/juetun/dashboard-api-main/web/validate"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
		post := admin_impl.NewControllerPost()
		trash := admin_impl.NewControllerTrash()
		img := admin_impl.NewControllerImg()

		c := r.Group(urlPrefix + "/console")
		postV := validate.NewValidate().NewPostV.MyValidate()

		c.GET("/post", post.Index)
		c.GET("/post/create", post.Create)
		c.POST("/post", postV, post.Store)
		c.GET("/post/edit/:id", post.Edit)
		c.PUT("/post/:id", postV, post.Update)
		c.DELETE("/post/:id", post.Destroy)
		c.GET("/post/trash", trash.TrashIndex)
		c.PUT("/post/:id/trash", trash.UnTrash)
		c.POST("/post/imgUpload", img.ImgUpload)
	})
}
