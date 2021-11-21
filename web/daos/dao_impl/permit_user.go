// Package dao_impl /**
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
