package admins

import "github.com/gin-gonic/gin"

type ConAdminCache interface {

	ClearCache(c *gin.Context)

	CacheParamList(c *gin.Context)

	ReloadAppCacheConfig(c *gin.Context)
}
