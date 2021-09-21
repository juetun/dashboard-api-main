/**
* @Author:changjiang
* @Description:
* @File:const
* @Version: 1.0.0
* @Date 2021/9/12 8:26 下午
 */
package parameters

const (
	MicroUser = "api-user" // 用户操作对象

	SystemPlatform = "platform"
	SystemBackend  = "backend"
	SystemSystem   = "system"
	SystemUser     = "user"

	DefaultSystem = SystemBackend // 默认系统
)

// 当前支持的系统列表
var SystemDescMap = map[string]SystemDescription{
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

type SystemDescription struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Desc  string `json:"desc"`
}
