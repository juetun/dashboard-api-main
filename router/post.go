package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/study/app-dashboard/lib/middlewares"
	"github.com/juetun/study/app-dashboard/web/controllers/con_impl"
	"github.com/juetun/study/app-dashboard/web/controllers/img_impl"
	"github.com/juetun/study/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		post := con_impl.NewControllerPost()
		trash := con_impl.NewControllerTrash()
		img := img_impl.NewControllerImg()

		c := r.Group(urlPrefix + "/console")
		p := c.Group("/post")
		postV := validate.NewValidate().NewPostV.MyValidate()

		p.GET("/", middlewares.Permission("console.post.index"), post.Index)
		p.GET("/create", middlewares.Permission("console.post.create"), post.Create)
		p.POST("/", middlewares.Permission("console.post.store"), postV, post.Store)
		p.GET("/edit/:id", middlewares.Permission("console.post.edit"), post.Edit)
		p.PUT("/:id", middlewares.Permission("console.post.update"), postV, post.Update)
		p.DELETE("/:id", middlewares.Permission("console.post.destroy"), post.Destroy)
		p.GET("/trash", middlewares.Permission("console.post.trash"), trash.TrashIndex)
		p.PUT("/:id/trash", middlewares.Permission("console.post.unTrash"), trash.UnTrash)
		p.POST("/imgUpload", middlewares.Permission("console.post.imgUpload"), img.ImgUpload)
	}, )
}
