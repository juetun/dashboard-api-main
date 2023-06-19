package parameters

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/library/common/app_param"
	"github.com/juetun/library/common/recommend"
)

//动态信息注册
func PluginAppendTrendType(arg *app_start.PluginsOperate) (err error) {
	if err = recommend.ReplaceTrendTypes(nil, &recommend.ArgReplaceTrendType{Types: SliceTrendType}); err != nil {
		return
	}
	app_param.AppendTrendsTypesAsModelItemOptions(SliceTrendType)
	return
}
