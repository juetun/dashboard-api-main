package parameters

import (
	"time"
)

const (
	// CacheKeyImportWithAppKey SKU缓存参数
	CacheKeyImportWithAppKey     = "p:app_c:%s"       // 接口按照appname缓存的key
	CacheKeyImportWithAppKeyTime = 7 * 24 * time.Hour // 接口按照appname缓存的生命周期

	CacheKeyUserGroupWithAppKey     = "p:u_grp:%d"       // 用户所属的管理组
	CacheKeyUserGroupWithAppKeyTime = 7 * 24 * time.Hour // 用户所属的管理组的生命周期

	CacheKeyUserGroupAppImportWithAppKey = "p:u_gai:%d_%s"       // 用户组每个接口权限缓存
	CacheKeyUserGroupAppImportTime       = 7 * 24 * time.Hour // 用户组每个接口权限缓存生命周期

)
