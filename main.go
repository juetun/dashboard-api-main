package main

import (
	"github.com/juetun/base-wrapper/lib/app_start"
	_ "github.com/juetun/base-wrapper/lib/init" // 加载公共插件项
	. "github.com/juetun/base-wrapper/lib/plugins"
	_ "github.com/juetun/dashboard-api-main/web/router" // 加载路由信息
)

// https://github.com/izghua/go-blog
func main() {
	app_start.NewPluginsOperate().Use(
		PluginJwt, // 加载用户验证插件,必须放在Redis插件后
	).LoadPlugins() // 加载插件动作

	// 启动GIN服务
	app_start.NewWebApplication().
		LoadRouter(). // 记载gin 路由配置
		Run()
}
