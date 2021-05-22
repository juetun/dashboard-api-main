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
)

type ConServiceImpl struct {
	base.ControllerBase
}

func (c2 ConServiceImpl) List(c *gin.Context) {
	panic("implement me")
}

func NewConServiceImpl() cons_outernet.ConService {
	p := &ConServiceImpl{}
	p.ControllerBase.Init()
	return p
}

