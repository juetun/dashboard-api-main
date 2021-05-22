package cons_outernet

import "github.com/gin-gonic/gin"

type System interface {
	Index(*gin.Context)
	Update(*gin.Context)
}

