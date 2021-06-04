package cons_adminnet

import "github.com/gin-gonic/gin"

type Img interface {
	ImgUpload(*gin.Context)
}
