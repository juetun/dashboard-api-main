package srvs

import "github.com/juetun/library/common/app_param"

type SrvCacheOperate interface {
	ClearCacheByKeyPrefix(arg *app_param.ArgClearCacheByKeyPrefix) (res *app_param.ResultClearCacheByKeyPrefix, err error)

	GetCacheParamConfig(arg *app_param.ArgGetCacheParamConfig) (res app_param.ResultGetCacheParamConfig, err error)
}
