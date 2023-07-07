package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type DaoCacheParam interface {
	GetListAsNext(arg *wrapper_admin.ArgCacheParamList, res *wrapper_admin.ResultCacheParamList) (list []*models.CacheKeyData, err error)

	GetCount(arg *wrapper_admin.ArgCacheParamList, res *wrapper_admin.ResultCacheParamList) (actResObject *base.ActErrorHandlerResult, err error)

	GetList(actResObject *base.ActErrorHandlerResult, arg *wrapper_admin.ArgCacheParamList) (res []*models.CacheKeyData, err error)

	BatchAddData(bases ...base.ModelBase) (err error)
}
