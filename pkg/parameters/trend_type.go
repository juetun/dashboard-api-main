package parameters

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/library/common/app_param"
)

//动态信息注册
func PluginAppendTrendType(arg *app_start.PluginsOperate) (err error) {
	app_param.AppendTrendsTypesAsModelItemOptions(SliceTrendType)
	return
}
