package daos

import (
	"github.com/juetun/dashboard-api-main/web/models"
)

type DaoGatewayPermit interface {

	// GetImportListByAppName 获取指定应用下的接口列表
	// refreshCache 是否强制刷新缓存，true 默认从数据库中读取数据更新缓存
	GetImportListByAppName(appName string, c bool) (res  map[int]models.AdminImport, err error)
}
