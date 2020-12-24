package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	con_impl2 "github.com/juetun/dashboard-api-main/web/cons/con_impl"
	"github.com/juetun/dashboard-api-main/web/validate"
)

func init() {
	app_start.HandleFunc = append(app_start.HandleFunc, func(r *gin.Engine, urlPrefix string) {
		post := con_impl2.NewControllerPost()
		trash := con_impl2.NewControllerTrash()
		img := con_impl2.NewControllerImg()

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
