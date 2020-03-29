package img_impl

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/services"
)

type ControllerImg struct {
	base.ControllerBase
}

func NewControllerImg() *ControllerImg {
	controller := &ControllerImg{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerImg) ImgUpload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "post.ImgUpload",
			"err":     err.Error(),
		})
		r.Response(c, 401000004, nil)
		return
	}

	filename := filepath.Base(file.Filename)
	dst := common.ConfigUpload.ImgUploadDst + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		r.Log.Error(map[string]string{
			"message": "post.ImgUpload",
			"err":     err.Error(),
		})
		r.Response(c, 401000005, nil)
		return
	}

	srv := services.NewQiuNiuService(&base.Context{Log: r.Log})
	// Default upload both
	data := make(map[string]interface{})
	if common.ConfigUpload.ImgUploadBoth {
		go srv.Qiniu(dst, filename)
		data["path"] = common.ConfigUpload.AppImgUrl + filename
		r.Response(c, 401000005, data)
		return
	}

	if common.ConfigUpload.QiNiuUploadImg {
		go srv.Qiniu(dst, filename)
		data["path"] = common.ConfigUpload.QiNiuHostName + filename
		r.Response(c, 401000005, data)
		return
	}

	data["path"] = common.ConfigUpload.AppImgUrl + filename
	r.Response(c, 401000005, data)
	return
}
