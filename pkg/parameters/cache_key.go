package parameters

import (
	"github.com/juetun/base-wrapper/lib/common/redis_pkg"
	"time"
)

var (
	CacheKeyHelp = redis_pkg.CacheProperty{
		Key:    "p:hp:%s",
		Expire: 30 * 24 * time.Hour,
	}
	CacheKeyHelpWithId = redis_pkg.CacheProperty{
		Key:    "p:hp_id:%s",
		Expire: 12 * time.Hour,
	}
	//接口按照appname缓存的key
	CacheKeyImportWithAppKey = redis_pkg.CacheProperty{
		Key:    "p:app_c:%s",
		Expire: 7 * 24 * time.Hour,
	}

	//用户所属的管理组
	CacheKeyUserGroupWithAppKey = redis_pkg.CacheProperty{
		Key:    "p:u_grp:%d",
		Expire: 7 * 24 * time.Hour,
	}

	//用户组每个接口权限缓存
	CacheKeyUserGroupAppImportWithAppKey = redis_pkg.CacheProperty{
		Key:    "p:u_gai:%d_%s",
		Expire: 7 * 24 * time.Hour,
	}
)
