package parameters

import (
	"time"
)

const (
	// CacheKeyImportWithAppKey SKU缓存参数
	CacheKeyImportWithAppKey     = "p:app_client:%s"  // 接口按照appname缓存的key
	CacheKeyImportWithAppKeyTime = 7 * 24 * time.Hour // 接口按照appname缓存的生命周期

)
