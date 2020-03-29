package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/impl/con_impl"
	"github.com/juetun/app-dashboard/web/controllers/impl/img_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		post := con_impl.NewControllerPost()
		trash := con_impl.NewControllerTrash()
		img := img_impl.NewControllerImg()

		c := r.Group(urlPrefix + "/console")
		p := c.Group("/post")
		postV := validate.NewValidate().NewPostV.MyValidate()

		p.GET("/", post.Index)
		p.GET("/create", post.Create)
		p.POST("/", postV, post.Store)
		p.GET("/edit/:id", post.Edit)
		p.PUT("/:id", postV, post.Update)
		p.DELETE("/:id", post.Destroy)
		p.GET("/trash", trash.TrashIndex)
		p.PUT("/:id/trash", trash.UnTrash)
		p.POST("/imgUpload", img.ImgUpload)
	}, )
}
