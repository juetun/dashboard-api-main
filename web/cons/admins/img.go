package admins

import "github.com/gin-gonic/gin"

type Img interface {
	ImgUpload(*gin.Context)
}
