/**
* @Author:changjiang
* @Description:
* @File:admin_user
* @Version: 1.0.0
* @Date 2020/10/18 11:59 上午
 */
package services

import (
	"github.com/juetun/base-wrapper/lib/base"
)

type AdminUserService struct {
	base.ServiceBase
}

func NewAdminUserService(context ...*base.Context) (p *AdminUserService) {
	p = &AdminUserService{}
	p.SetContext(context...)
	return
}

func (r *AdminUserService) GetUserByHid(userHid string) (res, err error) {
	return
}
