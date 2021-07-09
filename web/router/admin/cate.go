package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons_admin/con_impl"
	"github.com/juetun/dashboard-api-main/web/validate"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {
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
