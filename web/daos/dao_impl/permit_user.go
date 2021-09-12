/**
* @Author:changjiang
* @Description:
* @File:permit_user
* @Version: 1.0.0
* @Date 2021/9/12 12:05 下午
 */
package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/daos"
)

type DaoPermitUserImpl struct {
	base.ServiceDao
}

func NewDaoPermitUserImpl(ctx ...*base.Context) (res daos.DaoPermitUser) {
	p := &DaoPermitUserImpl{}
	p.SetContext(ctx...)
	return p
}