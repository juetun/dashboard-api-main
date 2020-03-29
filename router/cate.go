package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/impl/con_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		cate := r.Group(urlPrefix + "/console/cate")
		consoleCate := con_impl.NewControllerCategory()
		cateV := validate.NewValidate().NewCateV.MyValidate()
		cate.GET("/", consoleCate.Index)
		cate.GET("/edit/:id", consoleCate.Edit)
		cate.PUT("/:id", cateV, consoleCate.Update)
		cate.POST("/", cateV, consoleCate.Store)
		cate.DELETE("/:id", consoleCate.Destroy)
	})
}
