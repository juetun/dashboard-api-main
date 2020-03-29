package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/impl/con_impl"
	"github.com/juetun/app-dashboard/web/validate"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		cate := r.Group(urlPrefix + "/console")
		consoleCate := con_impl.NewControllerCategory()
		cateV := validate.NewValidate().NewCateV.MyValidate()
		cate.GET("/cate", consoleCate.Index)
		cate.GET("/cate/edit/:id", consoleCate.Edit)
		cate.PUT("/cate/:id", cateV, consoleCate.Update)
		cate.POST("/cate", cateV, consoleCate.Store)
		cate.DELETE("/cate/:id", consoleCate.Destroy)
	})
}
