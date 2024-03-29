package main

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	_ "github.com/juetun/base-wrapper/lib/app/init" // 加载公共插件项
	"github.com/juetun/base-wrapper/lib/app/middlewares"
	. "github.com/juetun/base-wrapper/lib/plugins" // 组件目录
	"github.com/juetun/dashboard-api-main/pkg/parameters"
	_ "github.com/juetun/dashboard-api-main/web/router" // 加载路由信息
	"github.com/juetun/library/common/plugins_lib"
)

// https://github.com/izghua/go-blog
func main() {

	app_start.NewPlugins().Use(
		PluginJwt, // 加载用户验证插件,必须放在Redis插件后
		PluginOss,
		PluginAppMap,
		plugins_lib.PluginWebMap,
		parameters.PluginAppendTrendType,
	).LoadPlugins() // 加载插件动作

	// 启动GIN服务
	_ = app_start.NewWebApplication(
		middlewares.HttpHeaderInfo(),
		middlewares.SignHttp(), // 添加签名中间件
	).LoadRouter(). // 记载gin 路由配置
		Run()
}
