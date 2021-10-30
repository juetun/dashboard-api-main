// Package admins /**
package admins

import (
	"github.com/gin-gonic/gin"
)

type ConService interface {
	List(c *gin.Context)
	Detail(c *gin.Context)
	Edit(c *gin.Context)
}
