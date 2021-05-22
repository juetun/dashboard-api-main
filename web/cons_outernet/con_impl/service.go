/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 5:16 下午
 */
package con_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons_outernet"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type ConServiceImpl struct {
	base.ControllerBase
}

func (r *ConServiceImpl) Edit(c *gin.Context) {
	var (
		arg wrappers.ArgServiceEdit
		err error
		res *wrappers.ResultServiceEdit
	)
	err = c.ShouldBind(&arg)
	if err != nil {
		r.Response(c, -1, nil, err.Error())
		return
	}
	srv := srv_impl.NewSrvServiceImpl(base.CreateContext(&r.ControllerBase, c))
	if res, err = srv.Edit(&arg); err != nil {
		r.Response(c, -1, res, err.Error())
		return
	}
	r.Response(c, 0, res, "success")
}

func (r *ConServiceImpl) List(c *gin.Context) {
	var (
		arg wrappers.ArgServiceList
		err error
		res *wrappers.ResultServiceList
	)
	err = c.ShouldBind(&arg)
	if err != nil {
		r.Response(c, -1, nil, err.Error())
		return
	}
	srv := srv_impl.NewSrvServiceImpl(base.CreateContext(&r.ControllerBase, c))
	if res, err = srv.List(&arg); err != nil {
		r.Response(c, -1, res, err.Error())
		return
	}
	r.Response(c, 0, res, "success")
}

func NewConServiceImpl() cons_outernet.ConService {
	p := &ConServiceImpl{}
	p.ControllerBase.Init()
	return p
}
