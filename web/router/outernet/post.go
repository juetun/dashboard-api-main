package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons_outernet/con_impl"
	"github.com/juetun/dashboard-api-main/web/validate"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet, func(r *gin.Engine, urlPrefix string) {
		post := con_impl.NewControllerPost()
		trash := con_impl.NewControllerTrash()
		img := con_impl.NewControllerImg()

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
	}, )
}
