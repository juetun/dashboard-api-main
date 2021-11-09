package intranets

import (
	"github.com/gin-gonic/gin"
)

type ConPermitIntranet interface {

	GetImportPermit(c *gin.Context)

}
