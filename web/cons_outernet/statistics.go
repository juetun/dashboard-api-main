package cons_outernet

import "github.com/gin-gonic/gin"

type Statistics interface {
	Index(*gin.Context)
}