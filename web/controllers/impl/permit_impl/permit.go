/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:06 上午
 */
package permit_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/controllers/inter"
	"github.com/juetun/dashboard-api-main/web/pojos"
	"github.com/juetun/dashboard-api-main/web/services"
)

type ControllerPermit struct {
	base.ControllerBase
}

func NewControllerPermit() inter.Permit {
	controller := &ControllerPermit{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerPermit) Menu(c *gin.Context) {
	var arg pojos.ArgPermitMenu
	err := c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000000, nil)
		return
	}

	arg.JwtUserMessage = r.GetUser(c)

	srv := services.NewPermitService(base.GetControllerBaseContext(&r.ControllerBase, c))

	res, err := srv.Menu(&arg)
	if err != nil {
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, res)
}
func (r *ControllerPermit) Flag(c *gin.Context) {
	var arg pojos.ArgFlag
	err := c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000000, nil)
		return
	}
	arg.JwtUserMessage = r.GetUser(c)
	srv := services.NewPermitService(base.GetControllerBaseContext(&r.ControllerBase, c))
	res, err := srv.Flag(&arg)
	if err != nil {
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, res)
}
