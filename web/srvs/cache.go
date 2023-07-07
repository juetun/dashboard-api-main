package srvs

import "github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"

type (
	SrvCache interface {
		CacheParamList(arg *wrapper_admin.ArgCacheParamList) (res *wrapper_admin.ResultCacheParamList, err error)

		ClearCache(arg *wrapper_admin.ArgClearCache) (res *wrapper_admin.ResultClearCache, err error)

		ReloadAppCacheConfig(arg *wrapper_admin.ArgReloadAppCacheConfig) (res *wrapper_admin.ResultReloadAppCacheConfig, err error)
	}
)
