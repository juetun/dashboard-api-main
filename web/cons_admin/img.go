package cons_admin

import "github.com/gin-gonic/gin"

type Img interface {
	ImgUpload(*gin.Context)
}
