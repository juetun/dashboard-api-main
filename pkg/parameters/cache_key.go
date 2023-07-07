package parameters

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/common/redis_pkg"
	"time"
)

var (
	CacheParamConfig = map[string]*redis_pkg.CacheProperty{
		"CacheHelpDocRelateByKeyUpdating": { //店铺变跟资质缓存的缓存Key
			Key:    "m:help:relate_k:%d",
			Expire: 12 * time.Hour,
		},
		"CacheKeyAdminUserWithUserHId": { //管理员用户存储
			Key:    "p:adu:%s",
			Expire: 30 * 24 * time.Hour,
		},
		"CacheKeyHelp": {
			Key:    "p:hp:%s",
			Expire: 30 * 24 * time.Hour,
		},
		"CacheKeyHelpWithId": {
			Key:    "p:hp_id:%s",
			Expire: 12 * time.Hour,
		},
		"CacheKeyAllHelpRelate": {
			Key:    "p:hp_all:%s",
			Expire: 12 * time.Hour,
		},
		//接口按照appname缓存的key
		"CacheKeyImportWithAppKey": {
			Key:    "p:app_c:%s",
			Expire: 7 * 24 * time.Hour,
		},
		//用户所属的管理组
		"CacheKeyUserGroupWithAppKey": {
			Key:    "p:u_grp:%d",
			Expire: 7 * 24 * time.Hour,
		},
		//用户组每个接口权限缓存
		"CacheKeyUserGroupAppImportWithAppKey": {
			Key:    "p:u_gai:%d_%s",
			Expire: 7 * 24 * time.Hour,
		},
		"CacheMenuImportWithAppKey": { //菜单有的接口权限
			Key:    "p:menu_ipt:%d_%s",
			Expire: 7 * 24 * time.Hour,
		},
	}
)

func GetCacheParamConfig(key string) (res *redis_pkg.CacheProperty, err error) {
	var ok bool
	if res, ok = CacheParamConfig[key]; ok {
		return
	}
	err = fmt.Errorf("您当前未配置缓存信息(%v)", key)
	return
}
