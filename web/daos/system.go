package daos

import "github.com/juetun/dashboard-api-main/web/models"

type DaoSystem interface {

	GetSystemList()(res []*models.ZBaseSys,err error)

}
