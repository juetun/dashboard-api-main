package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/study/app-dashboard/lib/middlewares"
	"github.com/juetun/study/app-dashboard/web/controllers/con_impl"
	"github.com/juetun/study/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		cate := r.Group(urlPrefix + "/console/cate")
		consoleCate := con_impl.NewControllerCategory()
		cateV := validate.NewValidate().NewCateV.MyValidate()
		cate.GET("/", middlewares.Permission("console.cate.index"), consoleCate.Index)
		cate.GET("/edit/:id", middlewares.Permission("console.cate.edit"), consoleCate.Edit)
		cate.PUT("/:id", middlewares.Permission("console.cate.update"), cateV, consoleCate.Update)
		cate.POST("/", middlewares.Permission("console.cate.store"), cateV, consoleCate.Store)
		cate.DELETE("/:id", middlewares.Permission("console.cate.destroy"), consoleCate.Destroy)
	})
}
