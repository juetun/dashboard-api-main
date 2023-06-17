// Package parameters /**
package parameters

import "github.com/juetun/base-wrapper/lib/base"

const (
	SystemPlatform = "platform"
	SystemBackend  = "backend"
	SystemSystem   = "system"
	SystemUser     = "user"

	DefaultSystem = SystemBackend // 默认系统
)

// 当前支持的系统列表
var (
	SliceTrendType = base.ModelItemOptions{}
	SystemDescMap = map[string]SystemDescription{
		SystemPlatform: {
			Key:   SystemPlatform,
			Label: "汽车",
		},
		SystemBackend: {
			Key:   SystemBackend,
			Label: "后台",
		},
		SystemSystem: {
			Key:   SystemSystem,
			Label: "系统管理",
		},
		SystemUser: {
			Key:   SystemUser,
			Label: "用户后台",
		},
	}
)

type SystemDescription struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Desc  string `json:"desc"`
}
