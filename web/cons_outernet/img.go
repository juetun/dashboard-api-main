package cons_outernet

import "github.com/gin-gonic/gin"

type Img interface {
	ImgUpload(*gin.Context)
}
