/**

 */
package main

import (
	"github.com/juetun/study/app-dashboard/lib/app_start"
	_ "github.com/juetun/study/app-dashboard/lib/init"    // 加载公共插件项
	. "github.com/juetun/study/app-dashboard/lib/plugins" // 将基本插件包与当前包共用
	_ "github.com/juetun/study/app-dashboard/router"      // 加载路由信息
)

// https://github.com/izghua/go-blog
func main() {
	app_start.NewPluginsOperate().Use(
		PluginLog,   // 加载日志插件
		PluginMysql, // 加载数据库插件
		PluginRedis, // 加载Redis插件
	).LoadPlugins() // 加载插件动作

	// 启动GIN服务
	app_start.NewWebApplication().
		LoadRouter(). // 记载gin 路由配置
		Run()
}
