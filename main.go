package main

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	_ "github.com/juetun/base-wrapper/lib/app/init" // 加载公共插件项
	. "github.com/juetun/base-wrapper/lib/plugins"  // 组件目录
	"github.com/juetun/dashboard-api-main/basic/myplugins"
	_ "github.com/juetun/dashboard-api-main/web/router" // 加载路由信息
)

// https://github.com/izghua/go-blog
func main() {
	app_start.NewPlugins().Use(
		PluginJwt, // 加载用户验证插件,必须放在Redis插件后
		PluginOss,
		PluginAppMap,
		myplugins.PluginUser, // 用户登录,jwt等用户信息逻辑处理
	).LoadPlugins() // 加载插件动作
	// 启动GIN服务
	app_start.NewWebApplication(
		// middlewares.SignHttp(), // 添加签名中间件
	).LoadRouter(). // 记载gin 路由配置
		Run()
}
