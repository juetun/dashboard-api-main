package intranets

import (
	"github.com/gin-gonic/gin"
)

type ConPermitIntranet interface {
	GetImportPermit(c *gin.Context)

	GetUerImportPermit(c *gin.Context) // 判断用户是否有接口权限

}
