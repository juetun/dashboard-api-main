package admins

import "github.com/gin-gonic/gin"

type AdminConPermitImport interface {

	UserPageImport(c *gin.Context)

	GetImport(c *gin.Context)

	MenuImport(c *gin.Context)

	MenuImportSet(c *gin.Context)

	ImportList(c *gin.Context)

	UpdateImportValue(c *gin.Context)

	EditImport(c *gin.Context)

	DeleteImport(c *gin.Context)

	GetImportByMenuId(c *gin.Context)
}
