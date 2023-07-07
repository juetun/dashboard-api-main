package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/dashboard-api-main/web/cons/admins/admin_impl"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet, func(r *gin.Engine, urlPrefix string) {

		con := admin_impl.NewConAdminCache()
		routeGroup := getRouteGroup(r, urlPrefix, "cache", )

		routeGroup.POST("/cache_param_list", con.CacheParamList)              //缓存参数列表

		routeGroup.POST("/clear_cache", con.ClearCache)                       //清除缓存

		routeGroup.POST("/reload_app_cache_config", con.ReloadAppCacheConfig) //重新加载配置
	})
}
