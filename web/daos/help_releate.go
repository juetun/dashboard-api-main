package daos

import "github.com/juetun/dashboard-api-main/web/models"

type DaoHelpRelate interface {
	AddOneHelpRelate(relate *models.HelpDocumentRelate) (err error)
}
