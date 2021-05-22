package cons_outernet

import "github.com/gin-gonic/gin"

type Trash interface {
	TrashIndex(*gin.Context)
	UnTrash(*gin.Context)
}