package init

import (
	"github.com/juetun/study/app-dashboard/lib/app_start"
	"github.com/juetun/study/app-dashboard/lib/common"
)

// 初始化加载内容
func init() {
	app_start.NewPluginsOperate().Use(
		common.PluginsApp, // 加载系统配置信息
		common.PluginsHashId,
	)
}
