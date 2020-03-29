package inter

import "github.com/gin-gonic/gin"

type Web interface {
	Index(*gin.Context)
	IndexTag(*gin.Context)
	IndexCate(*gin.Context)
	Detail(*gin.Context)
	Archives(*gin.Context)
	NoFound(*gin.Context)
}