// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package router

import (
	_ "github.com/juetun/dashboard-api-main/web/router/intranet" // 加载内网访问路由
	_ "github.com/juetun/dashboard-api-main/web/router/outernet" // 加载外网访问路由
	_ "github.com/juetun/dashboard-api-main/web/router/page"     // 加载网页访问路由
)

func init() {

}
